provider "aws" {
  region = "eu-north-1"
  profile = "terraform"
}

# -------------------
# ECR Repositories
# -------------------
resource "aws_ecr_repository" "engine" {
  name = "${var.project_name}-engine"
}

resource "aws_ecr_lifecycle_policy" "engine" {
  repository = aws_ecr_repository.engine.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Delete untagged images older than 7 days"
        selection = {
          tagStatus    = "untagged"
          countType    = "sinceImagePushed"
          countUnit    = "days"
          countNumber  = 7
        }
        action = { type = "expire" }
      }
    ]
  })
}

resource "aws_ecr_repository" "webserver" {
  name = "${var.project_name}-webserver"
}


resource "aws_ecr_lifecycle_policy" "webserver" {
  repository = aws_ecr_repository.webserver.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Delete untagged images older than 7 days"
        selection = {
          tagStatus    = "untagged"
          countType    = "sinceImagePushed"
          countUnit    = "days"
          countNumber  = 7
        }
        action = { type = "expire" }
      }
    ]
  })
}

# -------------------
# EC2 Instance
# -------------------
resource "aws_instance" "app" {
  ami           = data.aws_ami.amazon_linux.id
  instance_type = "t3.micro"
  associate_public_ip_address = true

  iam_instance_profile = aws_iam_instance_profile.ec2_profile.name
  vpc_security_group_ids = [aws_security_group.app_sg.id]

  user_data = <<-EOF
    #!/bin/bash
    set -e

    dnf update -y
    dnf install -y docker amazon-ssm-agent
    systemctl enable docker
    systemctl start docker
    systemctl enable amazon-ssm-agent
    systemctl start amazon-ssm-agent
    usermod -aG docker ec2-user
    mkdir -p /tmp/deployment_home/www

    curl -SL https://github.com/docker/compose/releases/latest/download/docker-compose-linux-$(uname -m) \
       -o /usr/libexec/docker/cli-plugins/docker-compose
    sudo chmod +x /usr/libexec/docker/cli-plugins/docker-compose

    cat > /tmp/deployment_home/config.json <<'CONFIGEOF'
    ${file("${path.module}/../config-docker.json")}
    CONFIGEOF

    cat > /tmp/deployment_home/engine.json <<'ENGINEEOF'
    ${file("${path.module}/../engine.json")}
    ENGINEEOF

    cat > /tmp/deployment_home/app-compose.yaml <<'COMPOSEEOF'
    ${file("${path.module}/../deploy/app-compose.yaml")}
    COMPOSEEOF

    cat > /tmp/deployment_home/www/index.html <<'INDEXEOF'
    ${file("${path.module}/../www/index.html")}
    INDEXEOF
  EOF

  tags = {
    Name = "${var.project_name}-app"
  }
}

# -------------------
# Roles
# -------------------
# IAM Role for EC2 to pull images + SSM
resource "aws_iam_role" "ec2_role" {
  name = "${var.project_name}-ec2-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "ec2_ssm" {
  role       = aws_iam_role.ec2_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_role_policy_attachment" "ec2_ecr" {
  role       = aws_iam_role.ec2_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

resource "aws_iam_instance_profile" "ec2_profile" {
  name = "${var.project_name}-ec2-profile"
  role = aws_iam_role.ec2_role.name
}

resource "aws_security_group" "app_sg" {
  name        = "${var.project_name}-sg"
  description = "Frontend + backend on one instance"

  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "App frontend dev port"
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_iam_openid_connect_provider" "github" {
  url = "https://token.actions.githubusercontent.com"

  client_id_list = [
    "sts.amazonaws.com"
  ]

  thumbprint_list = [
    "6938fd4d98bab03faadb97b34396831e3780aea1"
  ]
}

resource "aws_iam_role" "github_actions" {
  name = "${var.project_name}-github-actions"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = {
        Federated = aws_iam_openid_connect_provider.github.arn
      }
      Action = "sts:AssumeRoleWithWebIdentity"
      Condition = {
        StringEquals = {
          "token.actions.githubusercontent.com:aud" = "sts.amazonaws.com"
        }
        StringLike = {
          "token.actions.githubusercontent.com:sub" = "repo:mantzoun/yannena:*"
        }
      }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "github_ecr" {
  role       = aws_iam_role.github_actions.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser"
}

resource "aws_iam_role_policy_attachment" "github_ssm" {
  role       = aws_iam_role.github_actions.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMFullAccess"
}

data "aws_ami" "amazon_linux" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-*-x86_64"]
  }
}
