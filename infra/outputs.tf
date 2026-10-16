output "db_host" {
  description = "Host asignado por Neon para la base de preprod"
  value       = neon_project.preprod.database_host
}

output "database" {
  description = "Nombre de la base de datos de preprod"
  value       = neon_database.preprod.name
}

output "api_url" {
  description = "URL pública de la API de preprod en Render"
  value       = render_web_service.api.url
}

output "front_url" {
  description = "URL pública del Frontend de preprod en Render"
  value       = render_web_service.front.url
}