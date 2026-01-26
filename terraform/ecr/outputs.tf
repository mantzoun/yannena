output "ecr_engine_uri" {
  value = aws_ecr_repository.engine.repository_url
}

output "ecr_webserver_uri" {
  value = aws_ecr_repository.webserver.repository_url
}
