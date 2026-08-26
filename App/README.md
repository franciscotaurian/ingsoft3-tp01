# Realico Comidas - Catálogo y Pedidos

[![CI](https://github.com/franciscotaurian/ingsoft3-tp01/actions/workflows/ci.yml/badge.svg)](https://github.com/franciscotaurian/ingsoft3-tp01/actions/workflows/ci.yml)

Aplicación web para un restaurante desarrollada para la materia Ingeniería de Software III (Pipeline CI/CD).

## Configuración de Variables de Entorno

Antes de levantar la aplicación, se debe crear el archivo `.env` en la raíz del proyecto (`App/`):

```bash
cp .env.example .env
```

Luego editar `.env` con los valores deseados:

```env
# Base de datos PostgreSQL
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=catalogo

# Administrador del Panel Admin
ADMIN_USER=admin
ADMIN_PASSWORD=admin123

# Número de WhatsApp del restaurante (código país + número sin espacios)
WHATSAPP_NUMBER=5491112345678
```

> **Nota:** `DB_HOST=db` hace referencia al nombre del servicio de base de datos dentro de la red Docker. No se debe cambiar al levantar con Docker Compose.

---

## Levantar la Aplicación con Docker

### Requisito previo

Tener instalado [Docker](https://docs.docker.com/get-docker/) con el plugin Docker Compose (versión v2+).

---

### Opción A — Compilando las imágenes localmente (desarrollo)

Usa `docker-compose.yml`. Compila el backend en Go y el frontend en Node.js directamente desde el código fuente local.

```bash
# 1. Crear el archivo de entorno
cp .env.example .env

# 2. Construir las imágenes y levantar los contenedores
sudo docker compose up -d --build
```

La aplicación estará disponible en: **http://localhost**

Para detener y eliminar los contenedores:

```bash
sudo docker compose down
```

---

### Opción B — Descargando las imágenes desde el Registry (producción)

Usa `docker-compose.registry.yml`. Descarga las imágenes ya compiladas directamente desde **GitHub Container Registry (GHCR)** sin necesidad de tener Go ni Node.js instalados.

```bash
# 1. Crear el archivo de entorno
cp .env.example .env

# 2. Levantar descargando las imágenes publicadas
sudo docker compose -f docker-compose.registry.yml up -d --pull always
```

La aplicación estará disponible en: **http://localhost**

Para detener:

```bash
sudo docker compose -f docker-compose.registry.yml down
```

#### Imágenes publicadas en GHCR

| Servicio | Imagen |
|----------|--------|
| Backend  | `ghcr.io/franciscotaurian/realico-comidas-backend:v0.1.0` |
| Frontend | `ghcr.io/franciscotaurian/realico-comidas-frontend:v0.1.0` |

---

## Comandos Docker Útiles

```bash
# Ver contenedores en ejecución
sudo docker ps

# Ver logs de un contenedor
sudo docker logs realico_comidas_backend
sudo docker logs realico_comidas_frontend
sudo docker logs realico_comidas_db

# Ver todas las imágenes locales
sudo docker images

# Eliminar todos los contenedores detenidos
sudo docker container prune -f

# Eliminar imágenes sin uso
sudo docker image prune -a -f

# Limpiar todo el sistema Docker (contenedores, imágenes, redes, volúmenes)
sudo docker system prune -a --volumes -f
```

---

## Ejecución sin Docker (desarrollo local)

### Requisitos Previos
- PostgreSQL corriendo localmente con una base de datos creada (ej: `catalogo`)
- Go ≥ 1.22
- Node.js ≥ 18

### Backend (Go)
```bash
cd backend
cp .env.example .env
# Editar .env con DB_HOST=localhost
go mod download
go run main.go
```
*El backend escuchará en `http://localhost:8080` y ejecutará la migración de tablas automáticamente en PostgreSQL.*

### Frontend (React + Vite)
```bash
cd frontend
npm install
npm run dev
```
*El frontend escuchará en `http://localhost:5173`.*

---
## Stack Tecnológico

- **Backend**: Go (framework Gin + ORM GORM)
- **Base de Datos**: PostgreSQL
- **Frontend**: React + Vite (JavaScript) servido por Nginx
- **Configuración**: Variables de entorno mediante `.env`
- **Contenedores**: Docker + Docker Compose

---

## Arquitectura del Backend

El backend está estructurado en una **arquitectura de 3 capas con Interfaces e Inyección de Dependencias**:

```
Handler (Controller)
    │  recepta peticiones HTTP, valida input, llama al Service
    ▼
Service (interfaz + impl)
    │  contiene la lógica y las 8 reglas de negocio (para tests del TP5)
    ▼
Repository (interfaz + impl)
    │  ejecuta consultas contra PostgreSQL usando GORM
    ▼
PostgreSQL
```

### Inyección de Dependencias en `main.go`
Las dependencias se conectan de abajo hacia arriba sin variables globales:
- `Repository` recibe la conexión `gorm.DB`
- `Service` recibe la interfaz `Repository`
- `Handler` recibe la interfaz `Service`

---

## Reglas de Negocio Implementadas

1. **R1**: Los productos solo se muestran en el catálogo si `stock > 0`.
2. **R2**: El botón de agregar al carrito se inhabilita cuando la cantidad solicitada alcanza el stock disponible.
3. **R3**: Al confirmar un pedido, se descuenta el stock de cada producto en una transacción atómica.
4. **R4**: No se puede confirmar un pedido si la cantidad solicitada supera el stock actual disponible.
5. **R5**: El total del pedido es calculado por el servidor y guardado al confirmar.
6. **R6**: Transiciones de estado estrictas de un pedido: `pendiente` ➔ `confirmado` ➔ `entregado`.
7. **R7**: No se puede eliminar una categoría si tiene productos asociados (retorna HTTP 409 Conflict).
8. **R8**: El precio unitario de cada producto se congela (snapshot) al momento de registrar el pedido.

---

## Arquitectura Docker

La aplicación se compone de **3 contenedores** orquestados con Docker Compose:

```
┌─────────────────────────────────────────────────────────┐
│                      Docker Network                      │
│                                                         │
│  ┌─────────────────┐        ┌──────────────────────┐   │
│  │    frontend      │──/api/─▶      backend          │   │
│  │  (Nginx:alpine)  │        │  (Go binary:alpine)  │   │
│  │   Puerto 80      │        │    Puerto 8080        │   │
│  └────────┬─────────┘        └──────────┬───────────┘   │
│           │                             │               │
│    Expuesto al host                     │               │
│    localhost:80                         ▼               │
│                              ┌──────────────────────┐   │
│                              │         db            │   │
│                              │  (postgres:16-alpine) │   │
│                              │    Puerto 5432        │   │
│                              └──────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

- El **frontend** (Nginx) actúa como proxy inverso: las peticiones a `/api/` se redirigen internamente al contenedor `backend:8080`.
- El **backend** espera a que la base de datos esté lista (`healthcheck`) antes de iniciar.
- Los datos de PostgreSQL persisten en un **volumen Docker** (`postgres_data`).

---

## Pantallas de la Aplicación

1. **Catálogo (`/`)**: Muestra el menú disponible, carrito de compras y formulario para realizar el pedido por WhatsApp.
2. **Historial de Pedidos (`/pedidos`)**: Muestra la lista pública de pedidos y su estado en tiempo real.
3. **Panel Admin (`/admin`)**: Requiere autenticación HTTP Basic (`admin` / `admin123`). Permite gestionar categorías, productos, avanzar el estado de los pedidos y visualizar el dashboard de métricas.
