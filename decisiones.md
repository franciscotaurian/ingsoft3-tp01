# Decisiones del TP1

## Por qué Git no pudo resolver el conflicto solo — y qué habría tenido que pasar para que nunca apareciera.

Git no resulve solo el conflicto porque no entiende de logica y le da el control total al usuario para que decida cual es el codigo que debe quedar en la main. No tiene la capacidad de comprender el codigo, por lo tanto no puede decidir sobre el codigo. 

Git prefiere detener el flujo y pedir intervencion antes que generar un bug en el codigo.

Para evitar los conflictos es importante no usar ramas viejas ya que es probable que el main sufra cambios, entonces al momento de realizar merge de esa rama vieja es probable que modifiquemos lineas que hayan sido modificadas en algun commit hecho sobre main, lo que genera el conflicto.
Para evitarlo debemos hacer merges periodicos y pulls para tener actualizado el codigo en mi rama y tener conocimiento de los ultimos cambios.

## Qué problemas encontraste y cómo los solucionaste.

No tuve inconvenientes al realizar la actividad. Seguí paso a paso la guía y pude resolver la consigna.

## Declaración de uso de IA

No utilice IA para resolver la actividad, si la utilice para profundizar sobre conceptos teoricos como conflictos y como git trabaja con ellos. Además le pedi los comando para hacer el pull request por consola.

# Decisiones TP2

## Elección de aplicación

La aplicación elegida es un catalogo de comidas y gestor de pedidos desarrollado con:
* Backend: Golang utilizando Gin para gestión de solicitudes y GORM como ORM
* Frontend: React + Vite (Js)
* Base de datos: PostgreSQL

La aplicación permite al usuario administrador:
* CRUD categorias de platos.
* CRUD productos.
* Seguimiento y cambio de estado de pedidos.
* Visualización de métricas.
* Visualización del historial de pedidos.

La aplicación permite al usuario cliente:
* Seleccionar productos y almacenarlos en el carrito.
* Realizar el pedido.
* Notificar sobre el pedido vía mensaje de Whatsapp utilizando el link de la aplicación

### Criterios de selección validados

1. **Ejecución inmediata:** La aplicación clona, levanta y funciona de forma local en mi máquina sin inconvenientes, cumpliendo con el requisito base para poder iterar sobre ella durante el semestre.
2. **Comandos de compilación y ejecución:** Conozco los comandos exactos para levantar el entorno de cada capa. Para el backend utilizo `go mod download` para las dependencias y `go run main.go` para la ejecución en desarrollo (o `go build` para compilar el binario de producción). Para el frontend utilizo `npm install` seguido de `npm run dev` (o `npm run build` para generar los assets estáticos de producción).
3. **Configuración de la base de datos:** La cadena de conexión a PostgreSQL no está hardcodeada en ninguna parte del código. Todos los parámetros de conexión (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`) se gestionan de forma centralizada mediante variables de entorno que el backend inyecta leyendo un archivo `.env`. Esto asegura que cambiar el host de la base de datos de local a un contenedor de Docker no requiera alterar ni una sola línea de Go.
4. **Lógica de negocio testeable:** El código tiene la complejidad necesaria para generar tests con sentido, esquivando las simples pruebas de "el endpoint responde 200".
    
    *   **Backend:** Cuenta con 8 reglas de negocio estrictas ya implementadas y separadas en la capa de Servicios, ideales para evaluar casos válidos y casos borde:
        1. **Catálogo con Stock:** Solo se le muestran a los clientes públicos aquellos productos que tengan un stock mayor a 0.
        2. **Descuento de Inventario:** Al generar un pedido, el sistema calcula matemáticamente el stock remanente (`product.Stock - itemDTO.Quantity`) para actualizar el inventario.
        3. **Disponibilidad Estricta:** Si un cliente intenta pedir una cantidad de un producto que supera el stock actual disponible, la operación se rechaza por completo.
        4. **Cálculo del Total del Pedido:** El total no viene del Frontend (para evitar manipulaciones y hackeos), sino que se calcula 100% en el servidor multiplicando la cantidad solicitada por el precio real del producto.
        5. **Transiciones de Estado Protegidas:** Implementa una máquina de estados estricta. Un pedido solo puede pasar de `pendiente` ➔ `confirmado`, o de `confirmado` ➔ `entregado`. Cualquier otro salto es bloqueado retornando un error.
        6. **Snapshot de Precios:** Cuando se crea un pedido, el precio actual del producto se copia al ítem del pedido. Si el administrador cambia el precio del producto en el futuro, el pedido histórico mantiene su valor original.
        7. **Integridad de Datos del Cliente:** El sistema exige obligatoriamente nombre y dirección, y valida mediante una expresión regular (`phoneRegex`) que el teléfono solo contenga números.
        8. **Integridad del Producto:** Un producto no puede ser creado ni actualizado con un precio `<= 0`, un stock negativo (`< 0`), o sin que se le asocie una categoría real existente en la base de datos.
    
    *   **Frontend:** Presenta validaciones de interfaz propias para testear, destacándose el recálculo dinámico de los totales del carrito y la inhabilitación reactiva del botón de compra cuando la cantidad solicitada alcanza el stock máximo disponible.

## Decisiones a la hora de Contenerizar

### Imágenes base elegidas
* Base de datos — postgres:16-alpine
Se eligió la imagen oficial de PostgreSQL en su variante Alpine por su tamaño reducido y por su version fija que nos garantiza la reproducibilidad.

* Backend Go — golang:1.26-alpine → alpine:latest

Stage 1 (build): golang:1.26-alpine — contiene todo lo necesario que Go necesita para compilar. Solo existe durante el build.
Stage 2 (final): alpine:latest — imagen mínima de ~5 MB que solo recibe el binario ya compilado. Sin Go, sin código fuente, sin dependencias de compilación.

* Frontend React/Vite — node:20-alpine → nginx:alpine

Stage 1 (build): node:20-alpine — ejecuta npm install y npm run build para generar los assets estáticos. Solo existe durante el build.
Stage 2 (final): nginx:alpine — imagen mínima que sirve los archivos estáticos compilados (/dist) y actúa como proxy inverso hacia el backend.

Utilizamos multi-stage builds ya que las herramientas de compilación son innecesarias en producción. Estas se utilizan en una primera etapa para compilar la aplicación y sus dependencias, pero la etapa final rescata solamente el ejecutable listo para funcionar, sin arrastrar las herramientas de desarrollo. Esto genera grandes ventajas al descargar las imágenes desde el registro: al ser mucho más livianas, los tiempos de despliegue son menores y se reduce significativamente la superficie de ataque del sistema.

### Diferencia entre depend_on y healthcheck

Un contenedor se muestra como iniciado cuando su proceso principal arranca. Pero, en el caso de Postgres, el motor de BD, aunque arranque tarda un par de minutos en iniciar sus archivos y estar listo para recibir peticiones.
Si solo usamos depends_on al momento de que arranque el contenedor el backend intentara conectarse lo que dara un error ya que la BD no esta lista para recibir peticiones.
En cambio, con Healthcheck podemos agregar un comando de validacion (en nuestro caso pg_isready) cada 5 segundos hasta que la BD nos diga que esta lista para recibir solicitudes.
Por lo tanto al poner en depends_on, service healthy, el back no intentara contectarse hasta que el comando del healthcheck devuelva que esta listo.

## Secretos
Las claves y credenciales no se guardan directamente en el archivo docker-compose ya que este queda en el repositorio publico. Las credenciales viven en un archivo .env el cual no se carga al repositorio ignorandolo en el archivo .gitignore y el compose lee estos valores de ese archivo. De esta forma las credenciales quedan locales en nuestros equipos pero todo aquel que quiera usar el sistema puede saber que credenciales necesita para levantarlo.

### Lo que persiste del contenedor
El volumen postgres_data se declara en la sección volumes: del docker-compose.yml y sobrevive a docker compose down.

## Dockerización

Las imagenes que cree para mi backend y mi frontend fueron creadas para procesadores AMD.

**Problemas encontrados**
Tuve un inconveniente al correr las imagenes descargadas del registry, el cual se debia a que habia subido una version antigua de las imagenes a las cuales les faltaban el dependecy healthcheck y luego tenia que actualizarlo. La solucion fue borrar las imagenes antiguas y subir unas nuevas con todas las correcciones.

## Uso de IA
Para la realizacion del TP me ayude de la inteligencia artificial para definir las imagenes base que utilizo en los dockerfile. 
Me ayude tambien con algunos comandos para el dockerfile del front en la parte de ngix. Me sirvio para aclarar algunas dudas.
Con el inconveniente al correr las imagenes descargadas del registry me ayudo a solucionarlo. Aunque basicamente me recomendo borrarlas y volverlas a subir.

# Decisiones TP3

## Duración del sprint
Estableci la duración del sprint en 2 semanas. Como menciona el video, la duración conviene fijarla en base a los plazos de entrega de los trabajos practicos. Por lo tanto, configure 2 semanas, para poder finalizar las tareas antes de la primera entrega solicitada por los profes.
Además 2 semanas es un ciclo ideal que brinda el tiempo suficiente para desarrollar e integrar una funcionalidad con valor (como el catálogo), pero es lo bastante corto como para recibir feedback temprano y corregir el rumbo si algo sale mal.

## Cantidad de tareas
En este trabajo configure la cantidad de tareas maximas en 2, esto ya que soy una unica persona trabajando en el proyecto y prefiero no sobrecargarme de tareas. A medida que avancemos con el proyecto se puede ir ajustando este numero, revisando si nos sobra tiempo, pero es recomendable comenzar con un margen de tiempo que sobre para no incumplir con los acuerdos. Además, es la cantidad de personas + 1 que esta trabajando en el proyecto. Esto garantiza que pueda continuar trabajando si alguna tarea necesita aprobación para moverse al estado de DONE.

## Historia mal escrita

"Como desarrollador quiero crear la tabla usuarios"

Esta historia no cuenta con ciertas características que estas deben tener:
* No presenta el formato correcto: COMO, QUIERO, PARA. Este formato nos permite definir claramente a quien favorece el cambio, que cambio hay que realizar, que beneficio brinda. En este caso sin el para no clarifica cual es el motivo de agregar esta funcionalidad.
* Brinda una solución por defecto. Las historias de usuario deben escribirse describiendo la funcionalidad o modificación que se debe realizar al sistema pero sin obligar al equipo de desarrolladores a elegir un camino para implementarla. Luego el equipo determinara el camino correcto para implementar la historia y lo definira en las tareas. Esta historia determina el camino de implementación haciendo que los desarrolladores, que conocen el código al 100%, no tengan otra opcion que adaptarse a esa solucion.
* En el body faltan agrega criterios de aceptación. Esto genera que la finalización de la implementacion de la historia sea ambigua.

### Historia bien escrita

"Como administrador del sistema, quiero poder registrar a los usuarios para que puedan iniciar sesión y utilizar la plataforma."

Criterios de Aceptación:

* Datos obligatorios: El sistema debe requerir como mínimo: Nombre completo, Correo electrónico y Contraseña.

* Unicidad: Si el administrador intenta registrar un correo electrónico que ya existe en el sistema, se debe mostrar un mensaje de error: "El correo ya se encuentra registrado".

* Seguridad: La contraseña no debe guardarse en texto plano (debe estar encriptada/hasheada en la base de datos).

* Confirmación: Al registrar correctamente al usuario, el sistema debe mostrar un mensaje de éxito y limpiar el formulario.

* Validación de formato: El campo de correo electrónico debe validar que el texto ingresado tenga un formato válido (ejemplo@dominio.com).

## Declaración de uso de IA y problemas
Para este trabajo no utilice herramientas de inteligencia artificial. Me guié con los videos del profesor.
Tuve inconvenientes al principio, ya que cree el project por terminal y no se me habia habilitado por defecto el workflow de auto-add. Cuando lo active no me di cuenta que tenia un filtro para solo agregar bugs, asi que tuve que cargar las issue a mano y modificar el filtro y luego se comenzaron a agregar automaticamente.

# Decisiones TP4

## Estructura elegida del pipeline 
El pipeline se diseñó dividiendo la carga de trabajo en dos jobs independientes (`build-backend` y `build-frontend`) que se ejecutan de manera simultánea.
* Optimización del tiempo: Al ejecutarse en paralelo, el tiempo total de validación del Pull Request está determinado por el job más lento, en lugar de la suma secuencial de la construcción de ambos componentes. Esto reduce drásticamente los tiempos de espera y agiliza el flujo de trabajo ante múltiples iteraciones diarias.
* Integridad y validación del sistema como un todo: Aunque los procesos de compilación ocurren en paralelo y arquitectónicamente separados, ambos actúan como compuertas de calidad indivisibles para proteger la rama principal (`main`). Si ocurre un fallo de compilación en cualquiera de los dos extremos (front o back), el pipeline general reporta un estado de error, bloqueando el merge. Esto garantiza que ninguna integración parcial o rota llegue al código de producción.

## Estrategia de caché
El pipeline implementa el sistema de caché nativo de GitHub Actions para Docker Buildx, asegurando el aislamiento de los contextos mediante el atributo `scope` (`scope=backend` y `scope=frontend`).
* **Capas reutilizadas:** Se preservan y reutilizan las capas base y las descargas de dependencias del proyecto. Si no existen modificaciones en los archivos que gestionan las dependencias, estas capas pesadas se recuperan del caché casi instantáneamente, evitando descargas repetitivas por la red.
* **Capas invalidadas (no reutilizadas):** Cualquier capa correspondiente a los archivos de código fuente modificados en el commit actual y todas las instrucciones que le siguen se invalidan automáticamente, forzando su reconstrucción para garantizar que se evalúe el código más reciente.
* **Tolerancia a fallos de caché:** En caso de que el caché desaparezca (por expiración de retención o limpieza en GitHub), la robustez del pipeline no se ve comprometida. El proceso simplemente ejecutará, descargando y construyendo todas las capas desde cero. El CI seguirá funcionando de manera exitosa y segura, experimentando únicamente una degradación temporal en su velocidad de ejecución hasta que se genere el nuevo caché.

## Construcción vía Dockerfile vs. Compilación nativa
Se optó por delegar la validación a la construcción de las imágenes mediante sus respectivos `Dockerfile` en lugar de instalar las herramientas y compilar el código directamente sobre el *runner* de Ubuntu.
* Construir mediante Docker garantiza que el entorno donde se compila la aplicación en el proceso de Integración Continua es exactamente el mismo que se utilizará en la fase de Despliegue.
* El pipeline se mantiene independiente a las tecnologías subyacentes. No es necesario instalar ni mantener los SDKs de los lenguajes utilizados dentro de las configuraciones de GitHub Actions. Si a futuro se requiere una actualización en la versión de un framework o lenguaje, esta modificación queda encapsulada únicamente en el `Dockerfile`, permitiendo que el pipeline CI siga operando sin necesidad de refactorización.

## Inconvenientes

No tuve inconvenientes a la hora de realizar el práctico. Siguiendo el video del profe pude completarlo y comprenderlo sin dificultades.

## Declaración del uso de IA

Utilice inteligencia artificial para comprender con mayor detalle por construir via Dockerfile y no utilizar compilación nativa. El resto del trabajo pude realizarlo sin incovenientes.



