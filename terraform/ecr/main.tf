provider "aws" {
  region = "eu-north-1"
  profile = "terraform"
}

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
