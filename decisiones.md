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

Actualmente la aplicación puede ser corrida en mi maquina de forma local sin inconvenientes y cumple con los requisitos de la materia para, a futuro, por realizar los test y demás etapas del trabajo.

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

Utilizamos MultiStage ya que las herramientas de compilacion no deben estar en produccion.

### Lo que persiste del contenedor
El volumen postgres_data se declara en la sección volumes: del docker-compose.yml y sobrevive a docker compose down.

## Dockerización

Las imagenes que cree para mi backend y mi frontend fueron creadas para procesadores AMD.

## Uso de IA
Para la realizacion del TP me ayude de la inteligencia artificial para definir las imagenes base que utilizo en los dockerfile. 
Me ayude tambien con algunos comando para el dockerfile del front en la parte de ngix. Me sirvio para aclarar algunas dudas.
Tambien tuve un inconveniente al correr las imagenes descargadas del registry y me ayudo a solucionarlo. Aunque basicamente me recomendo borrarlas y volverlas a subir.

# Decisiones TP3

## Duración del sprint
Estableci la duración del sprint en 2 semanas. Como menciona el video, la duración conviene fijarla en base a los plazos de entrega de los trabajos practicos. Por lo tanto, configure 2 semanas, para poder finalizar las tareas antes de la primera entrega solicitada por los profes.

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

# Desiciones TP3

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



