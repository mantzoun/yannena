output "instance_id" {
  value = aws_instance.app.id
}

output "instance_public_ip" {
  value = aws_instance.app.public_ip
}

output "ecr_engine_uri" {
  value = aws_ecr_repository.engine.repository_url
}

output "ecr_webserver_uri" {
  value = aws_ecr_repository.webserver.repository_url
}
