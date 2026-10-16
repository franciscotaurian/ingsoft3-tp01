terraform {
  required_version = ">= 1.5"

  required_providers {
    neon = {
      source  = "kislerdm/neon"
      version = "0.18.0" # Versión fijada: evita incompatibilidades de actualización
    }
    render = {
      source  = "render-oss/render"
      version = "1.9.1"
    }
  }
}

provider "neon" {}   # Lee automáticamente NEON_API_KEY del entorno
provider "render" {} # Lee RENDER_API_KEY y RENDER_OWNER_ID del entorno

# --- Base de datos de PREPROD en Neon (no existe hasta que hagamos apply) ---

resource "neon_project" "preprod" {
  name                      = "miapp-preprod"
  history_retention_seconds = 21600
}

resource "neon_role" "preprod" {
  project_id = neon_project.preprod.id
  branch_id  = neon_project.preprod.default_branch_id # 🔴 default_branch_id, NO branch_id
  name       = "app_preprod_owner"
}

resource "neon_database" "preprod" {
  project_id = neon_project.preprod.id
  branch_id  = neon_project.preprod.default_branch_id
  name       = "app_preprod"
  owner_name = neon_role.preprod.name
}

# --- Servicios Web en Render para PREPROD ---

locals {
  # La cadena de conexión para PostgreSQL en Go (GORM).
  # Se DERIVA automáticamente de lo que devolvió Neon. Nadie la copia ni la pega.
  conn = "postgres://${neon_role.preprod.name}:${neon_role.preprod.password}@${neon_project.preprod.database_host}/${neon_database.preprod.name}?sslmode=require"
}

resource "render_web_service" "api" {
  name   = "miapp-api-preprod"
  plan   = "free"
  region = var.region

  runtime_source = {
    image = {
      image_url = var.image_repo_api # 🔴 SIN la etiqueta pegada
      tag       = var.image_tag      #    la etiqueta va en su propio campo
    }
  }

  env_vars = {
    DATABASE_URL = { value = local.conn }
  }
}

resource "render_web_service" "front" {
  name   = "miapp-front-preprod"
  plan   = "free"
  region = var.region

  runtime_source = {
    image = {
      image_url = var.image_repo_front
      tag       = var.image_tag
    }
  }

  env_vars = {
    BACKEND_URL  = { value = render_web_service.api.url } # ← La URL sale de Terraform, no de un copy-paste
    DNS_RESOLVER = { value = "8.8.8.8" }
  }
}