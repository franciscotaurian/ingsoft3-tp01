# Especificación de Requerimientos — Realico Comidas

**Sistema:** Catálogo de comidas con gestión de pedidos y panel de administración.
**Stack:** Go (Gin + GORM) · React + Vite · PostgreSQL

---

## Índice

1. [Descripción General](#1-descripción-general)
2. [Actores del Sistema](#2-actores-del-sistema)
3. [Requerimientos Funcionales (RF)](#3-requerimientos-funcionales-rf)
4. [Requerimientos No Funcionales (RNF)](#4-requerimientos-no-funcionales-rnf)
5. [Modelo de Datos](#5-modelo-de-datos)
6. [API REST — Referencia de Endpoints](#6-api-rest--referencia-de-endpoints)

---

## 1. Descripción General

**Realico Comidas** es una aplicación web de un restaurante que permite a clientes explorar el catálogo de productos y realizar pedidos vía WhatsApp, y a un administrador gestionar el catálogo y el ciclo de vida de los pedidos a través de un panel protegido por autenticación.

---

## 2. Actores del Sistema

| Actor | Descripción |
|---|---|
| **Cliente** | Usuario anónimo que navega el catálogo, arma un carrito y genera pedidos |
| **Administrador** | Usuario autenticado (registrado en la tabla `admins`) que gestiona categorías, productos, pedidos y visualiza métricas |

---

## 3. Requerimientos Funcionales (RF)

### RF-01 — Catálogo Público de Productos

| ID | Descripción |
|---|---|
| RF-01.1 | El sistema muestra la lista de productos disponibles agrupados por categoría |
| RF-01.2 | Solo se muestran productos con `stock > 0`; productos con stock igual a cero o negativo **no son devueltos** por el backend en ningún caso para el catálogo público |
| RF-01.3 | Cada producto muestra: nombre, descripción, precio, stock disponible e imagen (si tiene URL cargada) |
| RF-01.4 | Si no hay productos disponibles, se muestra un mensaje indicando que el menú está vacío |
| RF-01.5 | Los productos sin categoría asignada se agrupan bajo la sección "Varios" |

---

### RF-02 — Carrito de Compras

| ID | Descripción |
|---|---|
| RF-02.1 | El cliente puede agregar productos al carrito desde el catálogo |
| RF-02.2 | El botón "Agregar al carrito" se deshabilita cuando la cantidad en el carrito iguala el stock disponible del producto |
| RF-02.3 | El cliente puede aumentar o disminuir la cantidad de un ítem dentro del carrito |
| RF-02.4 | Si la cantidad de un ítem se reduce a cero, el ítem se elimina automáticamente del carrito |
| RF-02.5 | El cliente puede eliminar cualquier ítem individualmente desde el carrito |
| RF-02.6 | El carrito muestra el subtotal por ítem y el total general del pedido |
| RF-02.7 | El carrito se persiste en el `localStorage` del navegador; si el usuario recarga la página, el contenido se mantiene |

---

### RF-03 — Generación de Pedido

| ID | Descripción |
|---|---|
| RF-03.1 | Para confirmar un pedido, el cliente debe ingresar: nombre completo, teléfono y dirección de entrega |
| RF-03.2 | El sistema valida en el **frontend** que todos los campos estén completos antes de enviar la solicitud al backend |
| RF-03.3 | El campo teléfono valida en el **frontend** que contenga solo dígitos numéricos; el input rechaza cualquier carácter no numérico en tiempo real |
| RF-03.4 | El sistema valida en el **backend** que todos los campos del cliente estén presentes y no vacíos |
| RF-03.5 | El sistema valida en el **backend** que el teléfono contenga exclusivamente dígitos numéricos (`^\d+$`); de lo contrario retorna un error descriptivo |
| RF-03.6 | El sistema valida que el carrito no esté vacío antes de procesar el pedido |
| RF-03.7 | El sistema valida que la cantidad solicitada de cada ítem no supere el stock disponible en la base de datos al momento de confirmar |
| RF-03.8 | Al confirmar el pedido, el stock de cada producto involucrado se descuenta de manera atómica en la base de datos |
| RF-03.9 | El total del pedido es calculado por el servidor (suma de `cantidad × precio unitario` por ítem) y guardado en la base de datos |
| RF-03.10 | El precio unitario de cada ítem se congela (snapshot) al momento de confirmar; cambios futuros en el precio del producto no afectan pedidos anteriores |
| RF-03.11 | Todo pedido se crea con estado inicial `pendiente` |
| RF-03.12 | Tras confirmar exitosamente, el frontend genera y abre un mensaje de WhatsApp preformateado con el número del restaurante, incluyendo: número de pedido, datos del cliente y detalle de ítems con precios |
| RF-03.13 | Tras confirmar el pedido, el catálogo se recarga automáticamente para reflejar el stock actualizado |

---

### RF-04 — Historial de Pedidos (Solo Administrador)

| ID | Descripción |
|---|---|
| RF-04.1 | El historial de pedidos es una **sección privada**; solo es accesible para administradores autenticados |
| RF-04.2 | El botón de navegación "Historial de Pedidos" en la barra de navegación **solo es visible** cuando el usuario ha iniciado sesión como administrador |
| RF-04.3 | Si el usuario intenta acceder al historial sin credenciales, el frontend muestra un mensaje de acceso restringido indicando que debe iniciar sesión desde el Panel Admin |
| RF-04.4 | El backend rechaza con HTTP 401 cualquier solicitud al historial que no incluya credenciales de administrador válidas |
| RF-04.5 | El sistema muestra una tabla con todos los pedidos registrados |
| RF-04.6 | Cada fila del historial muestra: número de pedido, fecha/hora, nombre del cliente, teléfono, dirección, detalle de ítems (nombre + precio unitario + cantidad), total y estado actual |
| RF-04.7 | Los pedidos se ordenan del más reciente al más antiguo |
| RF-04.8 | El estado del pedido se muestra con un badge de color diferenciado (`pendiente` = naranja, `confirmado` = azul, `entregado` = verde) |

---

### RF-05 — Autenticación del Administrador

| ID | Descripción |
|---|---|
| RF-05.1 | El panel de administración y el historial de pedidos requieren autenticación mediante usuario y contraseña |
| RF-05.2 | Las credenciales de administrador se configuran mediante variables de entorno (`ADMIN_USER`, `ADMIN_PASSWORD`) |
| RF-05.3 | Al iniciar el backend, el sistema verifica si el administrador configurado en las variables de entorno existe en la tabla `admins` de la base de datos; si no existe, lo crea automáticamente |
| RF-05.4 | Si el administrador ya existía en la tabla con contraseña en texto plano, el sistema la actualiza automáticamente al hash SHA-256 en cada arranque |
| RF-05.5 | El backend valida las credenciales recibidas consultando la tabla `admins` en la base de datos; un usuario que no exista en dicha tabla no puede autenticarse, aunque conozca la contraseña del archivo `.env` |
| RF-05.6 | El token de autenticación (HTTP Basic Auth en Base64) se persiste en el `localStorage` para mantener la sesión entre recargas; el token es compartido entre el Panel Admin y el Historial de Pedidos |
| RF-05.7 | El administrador puede cerrar sesión desde el panel; esto elimina el token del `localStorage` y oculta el botón de Historial de la barra de navegación de forma inmediata |
| RF-05.8 | Si las credenciales son inválidas, el sistema muestra un mensaje de error y no permite el acceso |
| RF-05.9 | Todos los endpoints del grupo `/api/admin/*` requieren autenticación HTTP Basic Auth; sin credenciales válidas retornan HTTP 401 |

---

### RF-06 — Gestión de Categorías (CRUD — Admin)

| ID | Descripción |
|---|---|
| RF-06.1 | El administrador puede crear una nueva categoría ingresando su nombre |
| RF-06.2 | El nombre de la categoría es obligatorio; no puede ser una cadena vacía o solo espacios |
| RF-06.3 | El administrador puede editar el nombre de una categoría existente |
| RF-06.4 | El administrador puede eliminar una categoría |
| RF-06.5 | **No se puede eliminar una categoría que tenga al menos un producto asociado** — el sistema retorna HTTP 409 Conflict con un mensaje de error descriptivo |
| RF-06.6 | La lista de categorías se muestra en una tabla con ID, nombre y botones de acción |
| RF-06.7 | Las categorías son visibles públicamente para el catálogo del cliente (sin autenticación) |

---

### RF-07 — Gestión de Productos (CRUD — Admin)

| ID | Descripción |
|---|---|
| RF-07.1 | El administrador puede crear un nuevo producto con los campos: nombre, descripción, precio, stock, categoría e imagen (URL opcional) |
| RF-07.2 | El nombre del producto es obligatorio y no puede ser vacío |
| RF-07.3 | El precio debe ser un número mayor a 0; de lo contrario el sistema rechaza la operación |
| RF-07.4 | El stock no puede ser negativo; de lo contrario el sistema rechaza la operación |
| RF-07.5 | La categoría es un campo **obligatorio** al crear o editar un producto; si `category_id` es 0 o no se proporciona, el backend retorna un error antes de consultar la base de datos |
| RF-07.6 | Además de la validación por ID, la categoría debe existir en la base de datos; de lo contrario el sistema rechaza la operación |
| RF-07.7 | El sistema impide abrir el formulario de creación de producto si no existe ninguna categoría previa (muestra un aviso al administrador) |
| RF-07.8 | El administrador puede editar cualquier campo de un producto existente, incluyendo restablecer el stock cuando se reciben nuevas unidades |
| RF-07.9 | El administrador puede eliminar un producto |
| RF-07.10 | La vista de administración devuelve **todos los productos**, incluyendo los con stock ≤ 0, para que puedan ser editados y reabastecidos |
| RF-07.11 | En la tabla de administración, los productos con stock = 0 se resaltan en rojo |

---

### RF-08 — Gestión de Pedidos y Transición de Estado (Admin)

| ID | Descripción |
|---|---|
| RF-08.1 | El administrador visualiza todos los pedidos con el detalle completo: cliente, teléfono, dirección, productos, precio unitario, total y estado |
| RF-08.2 | El administrador puede avanzar el estado de un pedido según el flujo permitido |
| RF-08.3 | Las únicas transiciones de estado válidas son: `pendiente` → `confirmado` y `confirmado` → `entregado` |
| RF-08.4 | No se puede retroceder un estado ni saltear pasos; cualquier otra transición es rechazada por el backend con un error descriptivo |
| RF-08.5 | Un pedido en estado `entregado` muestra el texto "✓ Finalizado" y no ofrece más acciones |
| RF-08.6 | Cada acción de cambio de estado actualiza la pantalla automáticamente |

---

### RF-09 — Dashboard de Métricas (Admin)

| ID | Descripción |
|---|---|
| RF-09.1 | El panel de administración muestra un dashboard con 5 métricas calculadas en el servidor |
| RF-09.2 | **Total Facturado**: suma del campo `total` de todos los pedidos registrados |
| RF-09.3 | **Pedidos por Estado**: cantidad de pedidos agrupados por estado (`pendiente`, `confirmado`, `entregado`) |
| RF-09.4 | **Producto Más Vendido**: producto con mayor suma de cantidades en todos los pedidos; muestra "Ninguno" si no hay pedidos |
| RF-09.5 | **Total de Productos en Catálogo**: cantidad total de productos registrados, incluyendo los sin stock |
| RF-09.6 | **Valor del Inventario**: suma de `precio × stock` de todos los productos |

---

### RF-10 — Estado de Salud de la API

| ID | Descripción |
|---|---|
| RF-10.1 | El endpoint `GET /api/health` retorna HTTP 200 con el estado del servidor, accesible sin autenticación |

---

## 4. Requerimientos No Funcionales (RNF)

### RNF-01 — Arquitectura y Diseño de Software

| ID | Descripción |
|---|---|
| RNF-01.1 | El backend sigue una **arquitectura de 3 capas estricta**: Handler (Controller) → Service → Repository |
| RNF-01.2 | Cada capa expone una **interfaz Go** que define su contrato; la capa superior depende únicamente de la interfaz, nunca de la implementación concreta |
| RNF-01.3 | La **inyección de dependencias** se realiza manualmente en `main.go`, construyendo de abajo hacia arriba: Repository → Service → Handler; no se usan variables globales |
| RNF-01.4 | Las **reglas de negocio** residen exclusivamente en la capa `service`; los handlers no contienen lógica de dominio |
| RNF-01.5 | Las **validaciones de datos de entrada** (campos vacíos, precio negativo, teléfono inválido, stock insuficiente) son responsabilidad del backend; el frontend replica las validaciones críticas para mejorar la experiencia del usuario, pero el backend es la fuente de verdad |

---

### RNF-02 — Seguridad

| ID | Descripción |
|---|---|
| RNF-02.1 | Las contraseñas de administrador se almacenan en la base de datos **únicamente como hash SHA-256** (cadena hexadecimal de 64 caracteres); nunca se guarda ni se compara texto plano en la base de datos |
| RNF-02.2 | El hash SHA-256 se calcula usando la librería estándar de Go (`crypto/sha256`), sin dependencias externas |
| RNF-02.3 | Al verificar credenciales, el sistema hashea la contraseña recibida y compara el hash resultante contra el hash almacenado; la contraseña en texto plano nunca llega a la capa de repositorio |
| RNF-02.4 | Los endpoints de administración están protegidos por **HTTP Basic Auth**; las credenciales se validan en el middleware antes de llegar al handler |
| RNF-02.5 | El backend retorna HTTP 401 con el header `WWW-Authenticate` ante credenciales inválidas o ausentes en rutas protegidas |
| RNF-02.6 | La autenticación se valida contra la tabla `admins` en la base de datos; cualquier usuario que no exista en esa tabla no puede acceder, aunque conozca la contraseña del `.env` |

---

### RNF-03 — Configuración y Portabilidad

| ID | Descripción |
|---|---|
| RNF-03.1 | Toda la configuración sensible (credenciales de DB, credenciales de admin, número de WhatsApp, puerto) se gestiona mediante **variables de entorno** cargadas desde archivos `.env` |
| RNF-03.2 | Cada variable de configuración tiene un **valor por defecto** en caso de que la variable de entorno no esté definida, garantizando que la app arranque sin configuración mínima |
| RNF-03.3 | El backend puede apuntar a distintas bases de datos (desarrollo, QA, producción, servidor remoto) **sin modificar el código**, solo cambiando el valor de `DB_HOST`, `DB_NAME` y demás variables de entorno |
| RNF-03.4 | Se provee un archivo `.env.example` tanto en `/backend` como en `/frontend` con todas las variables necesarias documentadas; estos archivos no contienen credenciales reales y pueden subirse al repositorio |
| RNF-03.5 | Los archivos `.env` con credenciales reales deben estar incluidos en `.gitignore` y **nunca** subirse al repositorio de código |

---

### RNF-04 — Base de Datos

| ID | Descripción |
|---|---|
| RNF-04.1 | La creación y actualización del esquema de base de datos se ejecuta automáticamente al iniciar el backend mediante **GORM AutoMigrate**; no se requieren scripts SQL manuales |
| RNF-04.2 | La tabla `admins` se crea automáticamente en el primer arranque; el usuario administrador se inserta desde las variables de entorno si no existe |
| RNF-04.3 | Las operaciones que afectan múltiples tablas (guardar pedido + descontar stocks) se ejecutan en una **transacción atómica**; si algún paso falla, la operación completa se revierte |
| RNF-04.4 | Las relaciones entre entidades (FK) están definidas a nivel de modelo GORM y de base de datos |
| RNF-04.5 | La base de datos puede estar en el mismo servidor que el backend o en un servidor remoto; la app solo requiere ajustar `DB_HOST` en el `.env` para adaptarse a cualquier topología |

---

### RNF-05 — Comunicación entre Frontend y Backend

| ID | Descripción |
|---|---|
| RNF-05.1 | El backend configura **CORS** para aceptar peticiones desde el origen del frontend (`http://localhost:5173`) |
| RNF-05.2 | Toda la comunicación entre frontend y backend se realiza mediante **JSON sobre HTTP/REST** |
| RNF-05.3 | El frontend centraliza todas las llamadas HTTP en la carpeta `src/api/`; los componentes no realizan `fetch` directamente |
| RNF-05.4 | Los errores del backend (campo `error` en el JSON de respuesta) se propagan y muestran al usuario en la interfaz mediante alertas visuales diferenciadas (rojo para error, verde para éxito) |

---

### RNF-06 — Disponibilidad y Comandos de Ejecución

| ID | Descripción |
|---|---|
| RNF-06.1 | El backend se levanta con un único comando: `go run main.go` desde el directorio `/backend` |
| RNF-06.2 | El frontend se levanta con un único comando: `npm run dev` desde el directorio `/frontend` |
| RNF-06.3 | El backend escucha en el puerto configurado mediante `PORT` (por defecto `8080`) |
| RNF-06.4 | El frontend escucha en el puerto `5173` (Vite por defecto) |
| RNF-06.5 | El sistema no requiere Docker ni herramientas adicionales más allá de Go, Node.js y PostgreSQL |

---

### RNF-07 — Mantenibilidad y Testabilidad

| ID | Descripción |
|---|---|
| RNF-07.1 | Las **reglas de negocio** (stock, precio, teléfono, estados de pedido, categorías) están aisladas en la capa `service`, permitiendo su verificación mediante tests unitarios sin levantar una base de datos real |
| RNF-07.2 | El uso de interfaces permite **mockear el Repository** en tests de Service sin depender de GORM ni PostgreSQL |
| RNF-07.3 | La función `HashSHA256` reside en el paquete `utils`, aislada y testeable de forma independiente |

---

### RNF-08 — Interfaz de Usuario

| ID | Descripción |
|---|---|
| RNF-08.1 | El frontend usa **CSS vanilla** sin frameworks de UI para facilitar su modificación |
| RNF-08.2 | Los componentes de la interfaz están claramente separados en archivos individuales (`ProductCard`, `Cart`, `ProductForm`, `Header`, etc.) para facilitar su localización y modificación |
| RNF-08.3 | La interfaz muestra mensajes de retroalimentación visual distintos para éxito (verde) y error (rojo) |
| RNF-08.4 | La barra de navegación adapta su contenido al estado de autenticación: el botón "Historial de Pedidos" es visible únicamente cuando el administrador ha iniciado sesión; para usuarios no autenticados solo se muestran "Catálogo" y "Panel Admin" |
| RNF-08.5 | El estado de autenticación es compartido entre los componentes `Header`, `Admin` y `Orders` mediante estado React gestionado en `App.jsx`, sin duplicar lógica |

---

## 5. Modelo de Datos

```
admins
├── id          UINT (PK, autoincremental)
├── username    VARCHAR (ÚNICO, NOT NULL)
└── password    VARCHAR (NOT NULL) ← almacenado como hash SHA-256 (64 hex chars)

categories
├── id          UINT (PK, autoincremental)
└── name        VARCHAR (ÚNICO, NOT NULL)

products
├── id          UINT (PK)
├── name        VARCHAR (NOT NULL)
├── description VARCHAR
├── price       FLOAT (NOT NULL, > 0)
├── stock       INT (NOT NULL, ≥ 0, DEFAULT 0)
├── category_id UINT (FK → categories.id, NOT NULL, > 0)
└── image_url   VARCHAR

orders
├── id               UINT (PK)
├── customer_name    VARCHAR (NOT NULL)
├── customer_phone   VARCHAR (NOT NULL, solo dígitos)
├── customer_address VARCHAR (NOT NULL)
├── total            FLOAT (calculado al confirmar)
├── status           VARCHAR ('pendiente' | 'confirmado' | 'entregado', DEFAULT 'pendiente')
├── created_at       TIMESTAMP
└── updated_at       TIMESTAMP

order_items
├── id          UINT (PK)
├── order_id    UINT (FK → orders.id)
├── product_id  UINT (FK → products.id)
├── quantity    INT (≥ 1)
└── unit_price  FLOAT (snapshot del precio al momento del pedido)
```

---

## 6. API REST — Referencia de Endpoints

### Rutas Públicas (sin autenticación)

| Método | Ruta | Descripción | RF relacionados |
|---|---|---|---|
| GET | `/api/health` | Estado de la API | RF-10.1 |
| GET | `/api/products` | Lista productos con `stock > 0` | RF-01.1, RF-01.2 |
| GET | `/api/categories` | Lista todas las categorías | RF-06.7 |
| POST | `/api/orders` | Crea un pedido; valida teléfono, stock y campos obligatorios | RF-03.4 a RF-03.13 |

### Rutas de Admin (requieren HTTP Basic Auth — credenciales validadas contra tabla `admins`)

| Método | Ruta | Descripción | RF relacionados |
|---|---|---|---|
| GET | `/api/admin/orders` | Lista todos los pedidos (historial privado) | RF-04.1 a RF-04.8 |
| PATCH | `/api/admin/orders/:id/status` | Avanza el estado del pedido | RF-08.2 a RF-08.4 |
| GET | `/api/admin/products` | Lista todos los productos (incl. stock ≤ 0) | RF-07.10, RF-07.11 |
| POST | `/api/admin/products` | Crea un producto; requiere categoría válida y existente | RF-07.1 a RF-07.6 |
| PUT | `/api/admin/products/:id` | Edita un producto; revalida categoría | RF-07.8 |
| DELETE | `/api/admin/products/:id` | Elimina un producto | RF-07.9 |
| POST | `/api/admin/categories` | Crea una categoría | RF-06.1, RF-06.2 |
| PUT | `/api/admin/categories/:id` | Edita una categoría | RF-06.3 |
| DELETE | `/api/admin/categories/:id` | Elimina (falla con 409 si tiene productos) | RF-06.4, RF-06.5 |
| GET | `/api/admin/metrics` | Devuelve las 5 métricas del dashboard | RF-09.2 a RF-09.6 |
