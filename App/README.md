# Realico Comidas - Catálogo y Pedidos

Aplicación web para un restaurante desarrollada para la materia Ingeniería de Software III (Pipeline CI/CD).

## Stack Tecnológico

- **Backend**: Go (framework Gin + ORM GORM)
- **Base de Datos**: PostgreSQL
- **Frontend**: React + Vite (JavaScript)
- **Configuración**: Variables de entorno mediante `.env`

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

## Configuración y Variables de Entorno

### Backend (`backend/.env`)
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=catalogo

ADMIN_USER=admin
ADMIN_PASSWORD=admin123

WHATSAPP_NUMBER=5491112345678

PORT=8080
```

### Frontend (`frontend/.env`)
```env
VITE_API_URL=http://localhost:8080
VITE_WHATSAPP_NUMBER=5491112345678
```

---

## Comandos para Compilar y Ejecutar

### 1. Requisitos Previos
- PostgreSQL corriendo localmente con una base de datos creada (ej: `catalogo`).
- Go ≥ 1.22
- Node.js ≥ 18

### 2. Levantar el Backend (Go)
```bash
cd backend
cp .env.example .env
go mod download
go run main.go
```
*El backend escuchará en `http://localhost:8080` y ejecutará la migración de tablas automáticamente en PostgreSQL.*

### 3. Levantar el Frontend (React + Vite)
```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```
*El frontend escuchará en `http://localhost:5173`.*

---

## Pantallas de la Aplicación

1. **Catálogo (`/`)**: Muestra el menú disponible, carrito de compras y formulario para realizar el pedido por WhatsApp.
2. **Historial de Pedidos (`/pedidos`)**: Muestra la lista pública de pedidos y su estado en tiempo real.
3. **Panel Admin (`/admin`)**: Requiere autenticación HTTP Basic (`admin` / `admin123`). Permite gestionar categorías, productos, avanzar el estado de los pedidos y visualizar el dashboard de métricas.
