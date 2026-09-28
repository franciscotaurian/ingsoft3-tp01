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

# Decisiones TP5

## Qué lógica elegí testear y por qué

Enfoqué los tests en la **capa de servicios del backend** (`order_service.go`, `product_service.go`, `category_service.go`) y en los **módulos de lógica pura del frontend** (`pedido.js`, `client.js`).

La razon por la que se eligieron estas capas es por su importancia de la logica de negocio. Un error o bug en alguno de estos componentes pone en riesgo la integridad de los datos del sistema y el posible incumplimiento con el cliente. Por lo tanto es fundamental validar esta capa.

En el frontend, testear los componentes visuales de React no aportaba mucho valor real porque son capas de presentación. Lo que sí aporta valor es testear la lógica de armado del pedido (`pedido.js`) y el módulo que se comunica con la API (`client.js`).

## Umbral de coverage

**Backend:** Umbral de **70% sobre sentencias** (`statements`), que es la métrica nativa de `go tool cover`.

Elegí ese número porque es el umbral estándar que balancea exigencia real con practicidad. Tenía funciones puras de delegación (como `GetAll` y `GetByID`) que no tienen lógica de bifurcación, simplemente llaman al repositorio y devuelven el resultado. Incluirlas en la cuenta habría sido injusto penalizarían el porcentaje sin aportar seguridad real.

Filtrando esas funciones, la cobertura real sobre el código que importa fue de **71.2%**.

Sobre **ramas** (`branches`): no usé esa métrica como umbral en el backend porque `go tool cover` no la expone directamente de forma nativa. Sin embargo, revisando el reporte HTML manualmente, estimé que la cobertura de ramas ronda el **65%**, lo cual es coherente con que hay varios caminos del tipo `if err != nil` que no ejercita ningún test.

**Frontend:** Umbral de **80% sobre líneas, funciones, ramas y sentencias** configurado vía Vitest. La cobertura real cerró en **97.22%** porque los módulos testeados son pequeños y su lógica está muy acotada.

## Qué dejé afuera de la cuenta de cobertura

**Backend:**

- **`main.go`:** Es el punto de arranque de la aplicación. Inicializa la base de datos, levanta el servidor y conecta las dependencias. No tiene lógica de negocio testeable y no puede testearse de forma unitaria (requeriría una base de datos real corriendo).
- **Capa de repositorios (`repository/`):** Son adaptadores a la base de datos. Testearlos requeriría una base de datos real corriendo (test de integración), lo cual está fuera del alcance de este TP.
- **Capa de handlers (`handler/`):** Son los controladores HTTP. Testearlos implicaría levantar el servidor completo, lo cual también es un test de integración.
- **Modelos (`models/`):** Son structs de datos puros sin lógica. No hay nada que testear.
- **"Thin wrappers" de servicios:** Funciones como `GetAll()`, `GetByID()` o `GetMetrics()` que simplemente hacen `return s.repo.MétodoX(...)`. No tienen bifurcación ni regla de negocio: son pasamanos entre el handler y el repositorio. Se excluyeron filtrando sus líneas del archivo `.out` de cobertura antes de calcular el porcentaje.

**Frontend:**

- **Componentes React (`.jsx`):** Son capas de presentación. Testearlos requeriría herramientas extras como `@testing-library/react` y en la práctica testean más el renderizado visual que la lógica de negocio.
- La herramienta se configuró para apuntar únicamente a `src/pedido.js` y `src/client.js` mediante el atributo `include` de la configuración de Vitest.

## Por qué coverage alto no garantiza calidad

Un test puede ejecutar todas las líneas de una función y aun así no verificar nada útil.

**Ejemplo concreto en la app:** La función `CancelOrder` tiene un camino de éxito donde actualiza el estado del pedido a "cancelado" y lo retorna. Un test que llame a `svc.CancelOrder(1)` y solo verifique `assert.NoError(err)` le daría al reporte un 100% de cobertura en esa función. Sin embargo, si la función por algún bug guardara el estado "confirmado" en vez de "cancelado", ese test seguiría pasando. La cobertura no detectaría el error porque nunca verificó *qué* guardó, solo que no tiró una excepción.

Por eso, en nuestros tests siempre verificamos el comportamiento esperado y no solo la ausencia de error. Por ejemplo, `mockOrderRepo.AssertExpectations(t)` verifica que el mock fue llamado exactamente con los argumentos correctos.

## Pull Request bloqueado

Para demostrar el freno del pipeline, implementé la funcionalidad `CancelOrder` (con 3 caminos de ejecución distintos) y la subí **sin tests** en un Pull Request.

El check `build-backend` se puso en **rojo** con este mensaje en el log:

```text
Cobertura (sin thin wrappers): 65.2%
Umbral exigido              : 70%
❌ Cobertura insuficiente — el build es ROJO
```

Para solucionarlo, escribí los tests TestCancelOrder_Success, TestCancelOrder_NotFound y TestCancelOrder_InvalidStatus, que ejercitan los tres caminos de la función. La cobertura subió a 71.2% y el pipeline volvió a verde. Ese PR fue mergeado.

Luego escribi una nueva funcionalidad `ApplyDiscount`, con el mismo objetico -que de cobertura insufisiente-. Actualmente este PR con ApplyDiscount sin tests queda abierto y en rojo hasta la defensa, con el botón de merge desactivado.

Este freno es cualitativamente distinto al del TP4. En el TP4, el pipeline verificaba únicamente que el código compilara: si el build era verde, el PR podía mergearse. Un desarrollador podía eliminar todos los tests o escribir código completamente sin testear y el pipeline no lo detectaba. En el TP5, el pipeline tiene una segunda compuerta que verifica que el código nuevo esté cubierto por tests. Escribir código sin tests, aunque compile perfectamente, bloquea el merge.

Lo que este freno no detecta: tests con aserciones débiles o incorrectas. Si un test llama a la función pero no verifica el resultado con precisión, el pipeline se pone verde aunque el test no proteja nada real.

## Url's

https://github.com/franciscotaurian/ingsoft3-tp01/pull/23 (Corrida roja por cobertura)
https://github.com/franciscotaurian/ingsoft3-tp01/pull/22 (Corrida verde con reporte)
https://github.com/franciscotaurian/ingsoft3-tp01/pull/25 (PR abierto en rojo por cobertura)

## Refactorizacion para mockear

La aplicación ya estaba organizada con inyección de dependencias desde el TP2, por lo que no fue necesario refactorizar para habilitar los mocks. Cada servicio recibe sus repositorios por el constructor (NewOrderService(orderRepo, productRepo)), lo cual permite sustituirlos por dobles de prueba en los tests sin tocar el código de producción.

## Reglas de negocio agregadas
La aplicación ya tenía lógica de negocio real (validación de stock, formato de teléfono, transiciones de estado de pedidos). Sin embargo, para el ejercicio del freno de cobertura, agregué dos funcionalidades nuevas:

- CancelOrder: Permite al administrador cancelar un pedido, pero solo si está en estado pendiente. Si ya está confirmado o entregado, retorna un error de transición inválida.
- ApplyDiscount: Permite al administrador aplicar un descuento porcentual al total de un pedido pendiente. Valida que el porcentaje esté entre 1 y 100, que el pedido exista, y que no haya sido ya procesado.

## Stack utilizado

| Qué hace falta | En .NET (cátedra) | En nuestra app (Go + Vitest) |
|---|---|---|
| Test parametrizado | `[Theory]` + `[InlineData]` de xUnit | *Table-driven tests*: un `[]struct{...}` iterado con `t.Run()` |
| Doble (Mock) | Moq | `github.com/stretchr/testify/mock` |
| Medidor de cobertura | `coverlet` integrado en `dotnet test` | `go test -coverprofile` + `go tool cover` |
| Umbral que frena el build | `/p:Threshold=70` en el `ENTRYPOINT` del Dockerfile | Script bash con `awk` en el `ci.yml` que hace `exit 1` si el porcentaje es < 70 |
| Filtro de qué entra en la cuenta | `/p:Exclude=[DemoApi]Program*` en el `ENTRYPOINT` | `grep -v -E` sobre el archivo `.out` crudo en el `ci.yml` |

## Camino no cubierto

Inspeccionando el reporte HTML generado por go tool cover -html=coverage.out, identifiqué el siguiente camino no cubierto en product_service.go:

Línea: 75: if dto.Price <= 0 { return nil, ErrProductInvalidPrice }
Entrada concreta que lo recorrería: Llamar a svc.Create(CreateProductDTO{Name: "Empanada", Price: -5.0, Stock: 10, CategoryID: 1}). Con un precio negativo, la validación debería fallar y retornar ErrProductInvalidPrice.
Qué decidí hacer: No agregué el test. Los tests existentes sobre Create de productos cubren el camino de nombre vacío y de categoría inexistente, que son las reglas de negocio más críticas. El camino del precio negativo es importante pero, dado que ya superamos el umbral del 70% con holgura, lo dejé documentado como deuda técnica. Si el umbral fuera del 80%, sería el primer test a agregar.

## Problemas Encontrados

- Go no tiene flags nativas para umbral ni para excluir archivos. A diferencia de .NET con coverlet, go test no permite decirle "falla si la cobertura es menor al 70%" ni "excluí este paquete". Lo resolví completamente desde el ci.yml usando herramientas nativas de Linux (grep, awk) sobre el archivo de cobertura crudo que genera Docker.

- Conflicto en el PR de ApplyDiscount. Había creado la rama sobre un main desactualizado, por lo que al hacer el PR aparecieron conflictos con la rama que ya tenía CancelOrder mergeada. Lo resolví cerrando el PR sin mergear, borrando la rama remota desde GitHub, haciendo git pull en main local para actualizarlo, y recreando la rama desde el main actualizado.

## Declaración de uso de IA

Utilicé Antigravity (IA) como asistente de desarrollo durante todo el TP5. El rol de la IA fue de par de programación: propuse decisiones, la IA las implementó, y yo verifiqué y aprobé cada paso.

- Implementó el código de los tests en order_service_test.go, product_service_test.go y category_service_test.go.
- Reescribió el ci.yml completo, traduciendo la lógica de .NET a Go/Bash.
- Agregó las etapas test a los Dockerfiles del backend y del frontend.
- Implementó el código de las funcionalidades CancelOrder y ApplyDiscount.
Cómo lo verifiqué:

Corrí localmente ./check_coverage.sh después de cada cambio para confirmar que el porcentaje subía o bajaba según lo esperado.
Leí cada test generado y puedo explicar qué verifica cada assert. Por ejemplo, TestCreate_LlamaAlRepositorioConStockDescontado no verifica el valor de retorno: verifica que CreateWithTx fue llamado con el mapa {1: 7} (stock original 10 menos la cantidad pedida 3), usando AssertExpectations. Si la fórmula de descuento cambiara, ese test se rompería.
Identifiqué qué casos no están cubiertos: el camino de precio negativo en ProductService.Create (documentado en la sección anterior) y los errores de persistencia en base de datos (cuando el repositorio falla al guardar), que requieren mocks más complejos.
