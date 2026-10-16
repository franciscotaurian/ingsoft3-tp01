variable "image_repo_api" {
  description = "Repositorio de la imagen de Docker para la API (sin etiqueta)"
  type        = string
  default     = "ghcr.io/franciscotaurian/ingsoft3-tp01-backend"
}

variable "image_repo_front" {
  description = "Repositorio de la imagen de Docker para el Frontend (sin etiqueta)"
  type        = string
  default     = "ghcr.io/franciscotaurian/ingsoft3-tp01-frontend"
}

variable "image_tag" {
  description = "Tag de la imagen generada en el TP7 (sha-del-commit)"
  type        = string
  default     = "sha-5cddcf1b94bb59822ea03aac1b83982f9c81b98f"
}

variable "region" {
  description = "Región de Render para desplegar los servicios"
  type        = string
  default     = "oregon" # frankfurt | ohio | oregon | singapore | virginia
}