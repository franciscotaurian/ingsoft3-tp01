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

**Frontend:** Umbral de **70% sobre líneas y ramas (`lines: 70, branches: 70`)** configurado en `vite.config.js` vía `@vitest/coverage-v8`. La cobertura real cerró muy por encima (**97.22%**) porque los módulos testeados (`src/lib/pedido.js` y `src/api/client.js`) tienen lógica pura y acotada.

## Qué dejé afuera de la cuenta de cobertura

**Backend:**

- **`main.go`:** Es el punto de arranque de la aplicación. Inicializa la base de datos, levanta el servidor y conecta las dependencias. No tiene lógica de negocio testeable y no puede testearse de forma unitaria (requeriría una base de datos real corriendo).
- **Capa de repositorios (`repository/`):** Son adaptadores a la base de datos. Testearlos requeriría una base de datos real corriendo (test de integración), lo cual está fuera del alcance de este TP.
- **Capa de handlers (`handler/`):** Son los controladores HTTP. Testearlos implicaría levantar el servidor completo, lo cual también es un test de integración.
- **Modelos (`models/`):** Son structs de datos puros sin lógica. No hay nada que testear.
- **"Thin wrappers" de servicios:** Funciones como `GetAll()`, `GetByID()` o `GetMetrics()` que simplemente hacen `return s.repo.MétodoX(...)`. No tienen bifurcación ni regla de negocio: son pasamanos entre el handler y el repositorio. Se excluyeron filtrando sus líneas del archivo `.out` de cobertura antes de calcular el porcentaje.

**Frontend:**

- **Componentes React (`.jsx`):** Son capas de presentación. Testearlos requeriría herramientas extras como `@testing-library/react` y en la práctica testean más el renderizado visual que la lógica de negocio.
- La herramienta se configuró para apuntar únicamente a `src/lib/pedido.js` y `src/api/client.js` mediante el atributo `include` de la configuración de Vitest.

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

# Decisiones TP6

## Enlaces de este TP

### Paquetes públicos del registry
- **Backend:** https://github.com/franciscotaurian/ingsoft3-tp01/pkgs/container/ingsoft3-tp01-backend
- **Frontend:** https://github.com/franciscotaurian/ingsoft3-tp01/pkgs/container/ingsoft3-tp01-frontend


### Cadena de evidencia del registry
- **Corrida de PR — "Entrar al registry" salteado:** https://github.com/franciscotaurian/ingsoft3-tp01/actions/runs/36577456147
- **Corrida de main — "Construir y publicar la imagen" como último step:** https://github.com/franciscotaurian/ingsoft3-tp01/actions/runs/36578043434

### URLs de los entornos en vivo
| Entorno | Frontend | Backend |
|---|---|---|
| **QA** | https://ingsoft3-tp01.onrender.com | https://miapp-api-qa-w60e.onrender.com |
| **PROD** | https://miapp-front-prod-kdsx.onrender.com | https://miapp-api-prod-5djk.onrender.com |

---

## Por qué el artefacto se publica solo con la verificación en verde

El pipeline tiene una cadena de dependencias explícitas en el `ci.yml`: el job `deploy-qa` tiene `needs: [build-backend, build-frontend]`, y el step de publicación dentro de cada job de build solo llega a ejecutarse si todos los steps anteriores (compilación y tests) terminaron con éxito. Si los tests fallan, el job se marca en rojo y los steps siguientes nunca se ejecutan.

Si se publicara aunque los tests fallen, el registry dejaría de ser una fuente de verdad. Cualquier imagen ahí podría tener bugs conocidos y habría perdido el sentido de usarlo como artefacto verificado. El registry vale como garantía justamente porque solo llegan imágenes que pasaron por el pipeline completo.

La misma lógica aplica a la rama: la condición `push: ${{ github.event_name == 'push' && github.ref == 'refs/heads/main' }}` garantiza que solo los commits que llegaron a `main` (es decir, que pasaron por una PR con CI verde) generan una imagen publicada. Un commit de una rama de feature nunca produce un paquete en el registry.

**Los tres eslabones de la cadena y qué NO garantizan:**
La confiabilidad del registry se sostiene por tres controles encadenados en la configuración: (1) nada entra a `main` sin los checks del PR y el umbral de cobertura en verde, (2) solo los eventos `push` sobre `refs/heads/main` habilitan la publicación, y (3) el step de construir y publicar la imagen final está ubicado **al final del mismo job que corrió los tests** (sin `always()` ni `!cancelled()`), por lo que si un test falla el job se corta antes de llegar al push.
Sin embargo, esta cadena tiene dos límites que hay que conocer:
- **No impide un `docker push` manual:** Garantiza lo que publica *el pipeline*, pero alguien con permisos de escritura en el paquete podría subir una imagen a mano desde su terminal salteando el CI.
- **Las etiquetas (tags) son mutables:** Un tag como `sha-...` o `latest` puede ser sobreescrito en el registry; la única referencia criptográfica verdaderamente inmutable del contenido exacto de una imagen es su **digest** (`@sha256:...`).

## Continuous Delivery vs. Continuous Deployment

Implementé **Continuous Delivery**, no Continuous Deployment.

La diferencia está en el último paso: en Continuous Deployment, cada commit verde se despliega automáticamente en producción sin intervención humana. En Continuous Delivery, el pipeline automatiza todo hasta tener el artefacto listo en un entorno de QA verificado, pero el pase final a producción requiere una aprobación manual explícita.

Esta decisión es apropiada para el contexto de este proyecto porque:
- Producción tiene usuarios reales del restaurante. Un deploy automático fuera de horario o con un cambio no validado por el negocio podría interrumpir pedidos activos.
- El equipo es pequeño (un solo desarrollador), lo que hace que la aprobación manual sea liviana y no un cuello de botella burocrático.
- El restaurante tiene ventanas operativas claras (mediodía y noche), lo que hace que el deploy tenga que coincidir con un horario específico, algo que una aprobación manual garantiza y un deploy automático no.

---

## Diseño de la cadena: needs / if / environments y alcance de los secrets

El pipeline tiene cuatro jobs encadenados:

```
build-backend ──┐
                ├──► deploy-qa ──► deploy-prod
build-frontend ─┘
```

- **`build-backend` y `build-frontend`**: corren en paralelo en cualquier evento (PR o push). Validan compilación y tests. El step de publicación al registry tiene `if: github.event_name == 'push'` para saltear el login en PRs, y la condición `push: ${{ github.event_name == 'push' && github.ref == 'refs/heads/main' }}` para que solo main genere imagen.

- **`deploy-qa`**: tiene `needs: [build-backend, build-frontend]` (no corre si algún build falla) y `if: github.ref == 'refs/heads/main'` (no corre en PRs). Usa los secrets `RENDER_HOOK_API_QA` y `RENDER_HOOK_FRONT_QA` que viven en el environment `qa` de GitHub. Dispara el deploy del commit verificado usando `&ref=$GITHUB_SHA` y luego ejecuta el smoke test con reintentos.

- **`deploy-prod`**: tiene `needs: deploy-qa` (solo corre si QA está vivo) y `environment: production`. Este environment tiene un required reviewer configurado en GitHub: el workflow se pausa aquí hasta que el aprobador da el OK o lo rechaza. Los secrets `RENDER_HOOK_API_PROD` y `RENDER_HOOK_FRONT_PROD` están encapsulados en el scope del environment `production` y son inaccesibles para los jobs anteriores.

**Por qué el alcance de los secrets importa**: los secrets de producción solo pueden ser leídos por jobs que corren bajo el environment `production`. Esto garantiza que ningún job de build ni de QA pueda accidentalmente (o maliciosamente) disparar un deploy a producción. El aprobador, al aprobar, también está autorizando que esos secrets se expongan a ese job específico.

---

## Qué mira el aprobador antes de aprobar (criterios del gate)

Antes de aprobar el pase a producción, el aprobador verifica:

1. **QA está vivo y en verde**: el smoke test del job `deploy-qa` terminó con `✅ QA responde (API + BD + front)`. Si el smoke test falló, el gate de producción ni siquiera está disponible.
2. **Las funcionalidades validadas por el equipo de QA**: se realizó una prueba manual en el entorno de QA (navegando la app, cargando un pedido, verificando el flujo de admin) y el Product Manager dio el OK explícito.
3. **Ventana operativa**: el deploy se realiza en el horario permitido (tarde, entre el servicio del mediodía y el servicio nocturno del restaurante). No se aprueban pases en fin de semana ni fuera del horario de soporte.
4. **Sin cambios destructivos pendientes**: si el commit incluye migraciones de base de datos irreversibles, se verifica previamente que existe un plan de rollback de datos.

5. **Práctica del rechazo y límites de visibilidad del aprobador:** Se practicó el rechazo explícito de un despliegue en el environment `production` dejando asentado el motivo en el comentario de la corrida, comprobando que el job `deploy-prod` se cancela/falla sin exponer los secrets ni tocar producción. Asimismo, el aprobador sabe qué **no** puede ver desde la pantalla de GitHub: un `deploy-qa` en verde solo garantiza que los 3 `curl` del smoke test devolvieron HTTP 200, pero no detecta regresiones visuales ni errores funcionales que el smoke no consulta.

## La letra chica del free tier y cómo la maneja el pipeline

**Render (free tier):**
- **Cold start**: el backend tarda entre 30 y 60 segundos en despertar cuando no recibió tráfico en los últimos 15 minutos. El smoke test maneja esto con un loop de hasta 30 reintentos con `sleep 20` entre cada uno (hasta 10 minutos de espera total) y `--max-time 10` por cada `curl` para no colgarse si el servicio acepta la conexión pero no responde.
- **Sleep automático**: los servicios gratuitos se duermen tras 15 minutos sin tráfico. Para la defensa, conviene abrir las URLs manualmente antes para que estén despiertas.
- **750 horas mensuales compartidas en el Workspace:** Las 750 horas gratuitas no son por servicio, sino compartidas entre todos los Web Services gratuitos de la cuenta. Con 4 servicios (backend y frontend en QA y PROD), si estuvieran activos 24/7 consumirían 2.880 horas al mes y Render suspendería todo el workspace; el presupuesto cierra únicamente gracias al sleep automático a los 15 minutos sin tráfico.
- **500 minutos de build mensuales:** El plan gratuito de Render tiene un tope de 500 minutos de compilación al mes. En el TP6 cada deploy reconstruye los Dockerfiles consumiendo minutos de esa cuota (si se agota, los deploys quedan bloqueados hasta fin de mes).

**Neon (free tier):**
- La base de datos nunca se apaga (el proyecto free de Neon tiene una instancia activa permanente). No hay cold start en la BD.
- El límite es de almacenamiento (0.5 GB) y compute hours, que con el volumen de este proyecto no se alcanzan.
- **Por qué Neon y no el PostgreSQL de Render:** La base de datos PostgreSQL gratuita de Render **expira a los 30 días de su creación** (y se borra tras 14 días de gracia), por lo que no sirve para sostener los entornos durante el semestre; el plan gratuito de Neon no tiene vencimiento temporal.


## Garantía perdida: Render reconstruye desde el repo en vez de usar la imagen verificada

Actualmente el pipeline dispara un deploy hook de Render que le pide a Render que construya el commit `$GITHUB_SHA` desde el código fuente. Esto significa que Render hace su propio `docker build` en sus servidores.

La garantía que se pierde es la siguiente: la imagen que corre en QA y en PROD **no es exactamente la misma imagen que el pipeline verificó y publicó en GHCR**. Son imágenes construidas del mismo commit, pero compiladas en entornos diferentes (el runner de GitHub Actions vs los builders de Render). En teoría, el resultado debería ser idéntico, pero pueden diferir si alguna dependencia externa no está pineada con un hash exacto, si hay diferencias en el sistema operativo del builder, o si Render usa una capa de caché desactualizada.

## Qué prueba el smoke test y qué no

**Lo que prueba:**
- `GET /api/health` responde HTTP 200: el proceso del backend arrancó y el servidor HTTP está escuchando.
- `GET /api/products` responde HTTP 200: el backend está conectado a la base de datos (esta ruta hace una query real) y devuelve una respuesta válida.
- `GET /` del frontend responde HTTP 200: el contenedor de Nginx está vivo y sirve el HTML de la SPA.

**Lo que NO prueba:**
- Que la versión desplegada sea la del commit verificado (hasta antes de agregar `version` al health check, era imposible saberlo solo con el smoke test). A partir del commit donde se agregó `APP_VERSION` al endpoint `/api/health`, la respuesta incluye el campo `"version": "v6.0.1"` que permite verificarlo manualmente, aunque el smoke test no valida este campo explícitamente. Aunque luego fue eliminado con un rollback de produccion.
- Que la lógica de negocio funcione end-to-end (crear un pedido, descontar stock, etc.).
- Que el frontend esté conectado correctamente al backend de ese entorno (podría estar apuntando al backend equivocado y el smoke test pasaría igual).
- Que no haya regresiones visuales ni de UX.
- El Deploy Hook responde inmediatamente que recibió la orden y Render compila en segundo plano; mientras tanto (o si ese build nuevo falla), Render sigue sirviendo el contenedor de la versión anterior. Por lo tanto, el smoke test puede dar `200 OK` en los primeros intentos contra la versión vieja sin enterarse de que la nueva todavía no terminó de subir.

## Estrategia de despliegue elegida para producción real y plan de rollback

### Deployment pattern: Recreate en ventana de mantenimiento de la tarde

Para esta aplicación en producción real, la estrategia elegida es **Recreate** (apagar la versión vieja, levantar la nueva), ejecutada durante la **ventana de la tarde** (entre el servicio del mediodía y el servicio nocturno del restaurante).

**Justificación:**
- El restaurante opera en dos turnos claros: mediodía (~12:00-15:00) y noche (~20:00-23:00). La tarde (15:00-19:00) es una ventana natural sin usuarios activos, lo que hace que el tiempo de inactividad del Recreate (aproximadamente 2 minutos en Render) sea completamente invisible para los clientes.
- **Costo**: Recreate es la estrategia más simple de operar. No requiere infraestructura duplicada (como Blue-Green) ni un balanceador de carga sofisticado (como Canary). En un free tier de Render, no hay recursos para mantener dos ambientes de producción corriendo simultáneamente.
- **Riesgo**: bajo, dado que el deploy ocurre fuera del horario operativo. Si algo falla, se ejecuta el rollback (ver abajo) antes de que abra el turno de la noche.
- **Rollback**: al ser una sola instancia, el rollback es inmediato: se dispara el hook con el SHA del commit anterior.

**Por qué no Blue-Green ni Canary**: Blue-Green requiere tener dos ambientes de producción idénticos corriendo en paralelo (el doble del costo). Canary requiere un balanceador de carga que pueda dividir el tráfico por porcentaje. Ninguno de los dos es viable en el free tier actual, y ninguno agrega valor real si el deploy ya ocurre en horario sin usuarios.

### Plan de rollback paso a paso

Si tras un deploy a producción el smoke test falla o se detecta un problema en el turno de la noche:

1. Identificar el SHA del commit estable anterior (el tag de la release previa o el commit inmediatamente anterior en `main`):
   ```bash
   export SHA_ANTERIOR=$(git rev-list -n1 v6.0.0)
   ```
2. Leer los deploy hooks de PROD desde las variables de entorno o desde los secrets de GitHub (en una máquina con acceso):
   ```bash
   read -rs HOOK_API_PROD && export HOOK_API_PROD
   read -rs HOOK_FRONT_PROD && export HOOK_FRONT_PROD
   ```
3. Disparar el rollback del backend y del frontend al commit estable:
   ```bash
   curl -fsS "$HOOK_API_PROD&ref=$SHA_ANTERIOR"; echo
   curl -fsS "$HOOK_FRONT_PROD&ref=$SHA_ANTERIOR"; echo
   ```
4. Verificar que Render muestra el deploy del commit anterior como `Live`.
5. Confirmar manualmente que `/api/health` responde y que la app funciona.

**Tiempo medido de rollback**: el rollback realizado durante el TP tomó **~2 minutos** entre el inicio del deploy y el estado `Live` en Render (Deploy Started: September 30, 2026 a las 13:02 hs → Deploy Live: September 30, 2026 a las 13:04 hs).

**Qué NO deshace este rollback (Base de datos):**
Redesplegar el commit anterior vuelve atrás los contenedores de código en ~2 minutos, pero **no deshace los cambios en la base de datos de Neon**. Si el deploy que estamos revirtiendo ejecutó una migración destructiva (por ejemplo, borrar o renombrar una columna), el binario viejo fallará al querer leerla. Por eso los cambios de esquema requieren migraciones compatibles hacia atrás (*expand and contract*) o restaurar la base mediante Point-in-Time Recovery / branches de Neon.

---

## Problemas encontrados y cómo los resolviste

1. **Smoke test con 404 en el primer deploy**: el smoke test fallaba porque la URL configurada era `$URL_API/health` pero el backend de Go registra la ruta bajo `/api/health` (con el prefijo del grupo). Se corrigió la URL a `$URL_API/api/health` en el `ci.yml`.

2. **Problema al renombrar nginx.conf a default.conf.template**: al renombrar el archivo de configuración de Nginx para que coincida con el mecanismo de plantillas nativo de `nginx:alpine`, el `Dockerfile` del frontend quedó apuntando al nombre viejo (`COPY nginx.conf ...`). Se actualizó el `COPY` del Dockerfile para usar el nuevo nombre.

3. **Variable `${DNS_RESOLVER}` no definida en el compose**: al agregar la directiva `resolver ${DNS_RESOLVER}` al template de Nginx, el contenedor fallaba al arrancar en local porque esa variable no estaba en el `.env`. Se resolvió agregando el valor por defecto en el `Dockerfile` del frontend (`ENV DNS_RESOLVER=127.0.0.11`) y documentando en `.env.example` que en entornos cloud debe cambiarse a un DNS público.

4. **Segundo deploy a PROD con endpoints mal configurados**: en el primer deploy a producción el smoke test del job `deploy-prod` tenía configuradas las URLs de QA en vez de las de PROD. Se detectó al revisar el `ci.yml` y se corrigieron las variables de entorno del job de PROD para apuntar a `https://miapp-api-prod-5djk.onrender.com` y `https://miapp-front-prod-kdsx.onrender.com`.

---

## Declaración de uso de IA

Utilicé **Antigravity** (IA) como asistente principal durante todo el TP6. El rol fue de par de programación permanente: yo tomé todas las decisiones de diseño y arquitectura, y la IA implementó, analizó y diagnosticó problemas concretos.

Usos específicos:
- Diagnosticó por qué la imagen no aparecía en GitHub Packages a partir del log de la corrida de CI.
- Analizó el `ci.yml` e identificó que el smoke test llamaba a `/health` en vez de `/api/health`.
- Explicó la diferencia entre el enfoque con `if:` en el step de build vs. la condición embebida en `push:`.
- Detectó los tres problemas del nuevo `default.conf.template` (Dockerfile desactualizado, variable DNS_RESOLVER sin valor, y omisión de headers de proxy).
- Implementó el campo `version` en el endpoint `/api/health` modificando `config.go` y `main.go`.
- Redactó esta sección de `decisiones.md` en base a la información que yo fui aportando durante la sesión de trabajo.

**Cómo lo verifiqué:**
- El fix del smoke test (`/health` → `/api/health`) lo verifiqué corriendo `curl -fsS https://miapp-api-qa-w60e.onrender.com/api/health` desde la terminal antes de hacer el commit, confirmando el `200 OK` con el JSON de respuesta.
- Los cambios en `default.conf.template` y el `Dockerfile` los verifiqué levantando el stack local con `docker compose up --build` y navegando la app en `localhost:80`, confirmando que el frontend cargaba y el proxy al backend funcionaba correctamente.
- El campo `version` en el health check lo verifiqué leyendo el código modificado en `config.go` y `main.go` línea por línea y corriendo `go test ./...` para confirmar que los tests existentes seguían en verde después del cambio.

# Decisiones TP7 – Contenedores en el Pipeline, Integración y Pruebas e2e

## Enlaces del TP7

### 1. Paquetes públicos del Container Registry (ghcr.io)
- **Backend:** https://github.com/franciscotaurian/ingsoft3-tp01/pkgs/container/ingsoft3-tp01-backend
- **Frontend:** https://github.com/franciscotaurian/ingsoft3-tp01/pkgs/container/ingsoft3-tp01-frontend

### 2. Evidencia Central del Gate de Promoción (Integración 🟢 / e2e 🔴)
- **Commit que rompió el Front:** `d8054faf5de3bf71d88eb32a83dfc049af2f9041`
- **Corrida donde falló e2e y frenó la promoción:** https://github.com/franciscotaurian/ingsoft3-tp01/actions/runs/38204489860
- **Artefacto Reporte Integración (Verde):** `playwright-report-integracion` (sha256: `5b383e5dda2145703953cc81f92a9092aee2afbaaec8885cbab41c9c12f64d3d`)
- **Artefacto Reporte e2e (Rojo):** `playwright-report-e2e` (sha256: `349942968161a6aefc38f6f86ef83a4cf8a67d32778ac7fd798f4ac274331685`)

### 3. Corrida completa en verde posterior (hasta PROD)
- **Corrida exitosa verificada y desplegada a PROD:** https://github.com/franciscotaurian/ingsoft3-tp01/actions/runs/38205040873

### 4. URLs de los entornos en vivo
| Entorno | Frontend | Backend |
|---|---|---|
| **QA** | https://ingsoft3-tp01.onrender.com | https://miapp-api-qa-w60e.onrender.com |
| **PROD** | https://miapp-front-prod-kdsx.onrender.com | https://miapp-api-prod-5djk.onrender.com |

---

## 1. Build once, deploy many: La imagen como unidad de release

### Problema que resuelve respecto al TP6
En el TP6 la cadena de despliegue promovía *commits*: en cada deploy a QA y a PROD el pipeline disparaba un hook con `&ref=$GITHUB_SHA`, lo que provocaba que Render clonara el repositorio y **reconstruyera la aplicación desde cero** ejecutando `docker build`.
Esto rompía el principio fundamental de reproducibilidad: **lo que llegaba a producción no era idéntico bit a bit a lo verificado en QA**. 
Existía el riesgo latente de que dependencias externas resueltas en momentos distintos produjeran binarios dispares, o que variaciones en el tiempo de compilación generaran comportamientos imprevistos.

En este TP7 cerramos ese hueco implementando el principio **Build Once, Deploy Many**:
- El pipeline compila y empaqueta la aplicación **una sola vez** en la etapa de CI.
- Las imágenes inmutables resultantes se publican en GitHub Container Registry (`ghcr.io`) bajo la etiqueta unívoca del commit (`sha-${{ github.sha }}`).
- Tanto QA como PROD pasan a ser servicios **Image-backed**: Render no construye código; simplemente descarga y ejecuta **exactamente el mismo binario inmutable**.
- Promover a producción ya no significa "volver a compilar en PROD", sino invocar el Deploy Hook de producción pasándole el parámetro `imgURL` apuntando a la misma imagen que QA verificó con sus suites de pruebas.

---

## 2. Estrategia de etiquetas e inmutabilidad

### Qué garantiza cada etiqueta
1. **`sha-<commit>` en ghcr.io:** Garantiza trazabilidad estricta y bidireccional entre el código fuente y el contenedor. Conociendo el commit, sabemos exactamente qué imagen le corresponde; conociendo la imagen, sabemos con precisión matemática qué commit de git la produjo. Es la etiqueta que el pipeline usa para promover entre entornos.
2. **`v7.0.0` en git:** Es el nombre semántico de la release que corrió exitosamente en producción. Se coloca exclusivamente sobre el commit que está desplegado en producción (obtenido de *Deployments*), no a ciegas sobre el último commit de `main`.

### Por qué el pipeline NO publica `latest`
`latest` es un puntero mutable: hoy apunta a una imagen y mañana a otra. No ofrece ninguna garantía de contenido. Si el pipeline publicara `latest`, cualquier despliegue que lo referencie ejecutaría "lo último que se subió", sin importar si pasó las pruebas o si correspondió a una corrida intermedia. Con una única etiqueta por commit (`sha-...`), es imposible desplegar algo sin declarar explícitamente su identidad.

### Configuración en Render Web vs lo que realmente corre
En Render Web (*Settings → Existing Image*) configuramos los 4 servicios inicialmente con la imagen del primer merge verificado (`sha-e4f917f510...`).
Sin embargo, **esa no es la imagen que corre actualmente**. Render utiliza esa configuración como plantilla base para validar el registry y el nombre del paquete. Cada vez que nuestro pipeline ejecuta el Deploy Hook pasándole `--data-urlencode "imgURL=..."`, Render sobreescribe la etiqueta y despliega la imagen nueva enviada por el hook.
**Peligro del botón *Manual Deploy* en Render:** Como la pantalla *Settings → Image* del dashboard conserva escrita la etiqueta inicial (`sha-e4f917f510...`) y los Deploy Hooks con `imgURL` no actualizan ese texto, si alguien presiona **"Manual Deploy"** desde el panel web de Render **provocará una regresión en producción o QA**, ya que Render volverá a desplegar aquella primera imagen vieja en lugar de la última promovida por el pipeline.

### Conservación de servicios existentes
No creamos servicios nuevos en Render; en los 4 servicios existentes (APIs y Frontends de QA y PROD) editamos la fuente pasando de *Git Provider* a *Existing Image*. Esto permitió conservar las URLs públicas, las variables de entorno ya configuradas, los Deploy Hooks existentes y no duplicar el consumo de horas del plan gratuito (750 horas compartidas por workspace).

---

## 3. Trazabilidad: De la Release a la Imagen en un solo paso

Para pasar del tag de la release `v7.0.0` a las imágenes exactas desplegadas en producción, basta un único comando de Git:
```bash
git rev-list -n1 v7.0.0
```
Este comando devuelve el SHA del commit (los 40 caracteres). Sabiendo que el pipeline publica las imágenes con la convención fija `sha-<commit>`, las imágenes exactas en el registry son de forma inmediata:
- **Backend:** `ghcr.io/franciscotaurian/ingsoft3-tp01-backend:sha-<SHA>`
- **Frontend:** `ghcr.io/franciscotaurian/ingsoft3-tp01-frontend:sha-<SHA>`

---

## 4. Verificación externa del entorno y limitaciones del Smoke Test

### Cómo se comprueba desde afuera que el entorno ejecuta la imagen
Se comprueba ingresando a la pestaña **Events** de cada servicio en Render:
- El evento de despliegue indica textualmente: **«Triggered via Deploy Hook»** y nombra la imagen del registro con su etiqueta `sha-...`.
- En un servicio que compila desde el repo, el evento indicaría «Deploying commit...» y registraría la fase de construcción de Docker.

### Qué NO alcanza a probar el Smoke Test
El smoke test realiza peticiones `curl` a `/api/health`, `/api/products` y `/`. Solo comprueba que el servicio **responde HTTP 200**.
**Limitación:** Si un deploy en Render fallara silenciosamente y Render continuara sirviendo el contenedor anterior, el smoke test seguiría dando verde porque el entorno sigue vivo, sin poder advertir que no se levantó la imagen nueva. Tampoco prueba que los botones de la pantalla funcionen o que un formulario web se comunique adecuadamente con el backend.



---

## 5. Diseño de las Suites de Pruebas: Integración y e2e

### Suite de Integración (`App/frontend/e2e/api.spec.js`)
Ejecutada con el fixture `request` de Playwright, **sin navegador**, comunicándose directamente contra la API de QA (`API_BASE_URL`):
1. **Alta, verificación y borrado:** Crea un producto administrativo con nombre único (`POST /api/admin/products`), verifica que la base de datos real lo expone en el catálogo (`GET /api/products`), lo elimina (`DELETE /api/admin/products/:id`) y comprueba que desapareció de la base.
2. **Rechazo de payload inválido:** Envía un producto con datos inválidos (nombre vacío / precio negativo), comprueba que la API rechaza con `400 Bad Request` y verifica que la cantidad de registros en la base no se modificó.
3. **Validación de regla de negocio en pedidos (elegida por el equipo):** Envía un pedido (`POST /api/orders`) con un teléfono con caracteres no numéricos. Comprueba que el backend aplica la regla estricta de regex y responde `400`.
4. **Ciclo completo de stock (decremento y restauración):** Crea un pedido real que descuenta stock en la BD, comprueba que el stock disponible disminuyó exactamente en la cantidad pedida, y luego invoca el endpoint administrativo de cancelación (`POST /api/admin/orders/:id/cancel`) verificando que el stock de la base de datos real vuelve a su valor original.

### Suite End-to-End (`App/frontend/e2e/products.spec.js`)
Ejecutada con Chromium real contra el Frontend de QA (`E2E_BASE_URL`):
1. **Flujo crítico de compra diario con cleanup y stock inmutable:** El usuario ingresa a la aplicación, registra el stock inicial visible en la tarjeta del producto, lo agrega al carrito, completa los datos de entrega y confirma el pedido. Posteriormente, para cumplir con el principio de **idempotencia y cleanup**, el test navega al Panel Admin, se autentica, ubica la orden creada por el cliente de la prueba y presiona **"Cancelar Pedido ✕"**. Finalmente, regresa al catálogo y comprueba con aserción estricta de Playwright que el stock en pantalla retornó exactamente a su valor inicial.
2. **Validación de formulario incompleto:** Intenta enviar el pedido sin completar los datos requeridos y verifica que la interfaz retiene el foco y previene el envío.
3. **Flujo de administración y login (elegido por el equipo):** El usuario ingresa a la pestaña "Panel Admin", introduce las credenciales (`admin` / `admin123`), presiona "Ingresar", verifica el acceso al "Panel de Administración" con sus métricas y realiza el logout exitoso.

### Qué NO pusimos en cada suite (Pirámide de Testing)
- **No duplicamos pruebas unitarias:** Las validaciones de lógica pura de pedidos (`pedido.js`) y cálculo de descuentos ya están cubiertas en milisegundos en el TP5 con Vitest.
- **En e2e no testeamos combinaciones exhaustivas:** Levantar un navegador real es costoso en tiempo y recursos. En e2e solo van los 3 flujos críticos indispensables del negocio.

---

## 6. Diagnóstico del par Verde/Rojo: Integración vs e2e

En la rotura provocada intencionalmente en la aplicación, modificamos el texto del botón en `ProductCard.jsx` a `"Comprar Plato"`.

El resultado del pipeline fue:
- **Smoke test:** 🟢 VERDE (el servidor y la API responden 200).
- **Integración:** 🟢 VERDE (la API Go y la base de datos PostgreSQL funcionan y responden perfecto a las llamadas directas).
- **e2e:** 🔴 ROJO (el navegador no encontró el botón que el usuario necesita para operar).
- **`deploy-prod`:** ⚪ BLOQUEADO (el gate impidió la promoción a producción).

### Lectura diagnóstica (Tabla §2.5)
| Integración | e2e | Diagnóstico |
|---|---|---|
| 🟢 | 🔴 | **El frontend está roto:** La API y la BD están sanas, pero la pantalla no utiliza o no presenta los controles correctamente. |
| 🔴 | ⚪ (Salteado) | **La API o la base de datos están rotas:** El backend falló; la e2e ni siquiera llega a gastar un navegador. |
| 🟢 | 🟢 | **Sistema íntegro:** Ambos niveles verificados. |

Este par demostró el valor del diagnóstico: supimos inmediatamente que la falla residía en la interfaz del frontend sin necesidad de revisar el código del backend.

---

## 7. Integración Amplia vs Estrecha

Elegimos la **Integración Amplia** (hablarle a la API de QA ya desplegada en Render con su base PostgreSQL real):
- **Ventajas:** Máximo realismo; prueba la aplicación en condiciones de red reales, con el driver de PostgreSQL real y la infraestructura del proveedor cloud, sin necesidad de orquestar bases de datos temporales en el runner.
- **Desventajas:** Depende de que el entorno de QA esté disponible y despierto (cold start); las corridas deben ser secuenciales para no colisionar datos.
- **Frente a la Integración Estrecha (en memoria en el runner):** La estrecha sería más rápida y aislada, pero requeriría levantar contenedores de base de datos descartables en GitHub Actions y no probaría el despliegue real en la nube.

- **Límite conocido de tener un único QA compartido entre corridas:** Si se realizan dos merges seguidos a `main` (corridas A y B), el `deploy-qa` de B puede redesplegar QA mientras las suites de `integracion` o `e2e` de A están corriendo. Esto puede hacer fallar a A o provocar un falso verde donde A pruebe parcialmente la imagen de B. Se reconoce comparando la imagen en *Events* de Render con la hora de ejecución en *Actions*; si ocurre, se rechaza manualmente la corrida A en el gate de producción (dado que `concurrency` en `deploy-prod` no cancela corridas en estado *Waiting for review*) y se aprueba únicamente la B, manteniendo la regla operativa de hacer un merge por vez.
---

## 8. Cold Start, Timeouts y Tests Flaky

En el plan gratuito de Render, las instancias se suspenden tras periodos de inactividad. Para evitar falsos negativos:
- Configuramos `timeout: 60_000` (60 segundos) por test y `expect: { timeout: 15_000 }` en `playwright.config.js`.
- Configuramos `retries: 1` para absorber demoras transitorias de red.
- **Test Flaky:** Es un test que a veces pasa y a veces falla sin cambios en el código. Es sumamente perjudicial porque acostumbra al equipo a desestimar las alertas rojas. Playwright marca explícitamente en el reporte cuando un test requirió un retry para pasar, permitiendo auditar la estabilidad.

---

## 9. Una sola imagen de Frontend para QA y PROD

La imagen de frontend empaqueta los archivos estáticos compilados de React servidos por Nginx (`nginx:alpine`):
- **Adentro de la imagen:** Los assets estáticos de producción (`dist/`) y la plantilla de Nginx (`default.conf.template`).
- **Configurado por variable de entorno:** `BACKEND_URL` y `DNS_RESOLVER`. 
En tiempo de arranque del contenedor, Nginx sustituye `${BACKEND_URL}` con la URL de la API correspondiente (QA o PROD). De esta forma, **la misma imagen inmutable funciona en ambos entornos** sin necesidad de recompilar el frontend.

---

## 10. Problemas encontrados y cómo los resolvimos

1. **Conexión rechazada `ECONNREFUSED 127.0.0.1:8080` en local:** El contenedor backend no exponía el puerto al host. Se mapeó `ports: - "8080:8080"` en `App/docker-compose.yml` para permitir pruebas locales directas sin intermediación de Nginx.
2. **Fallo de suite e2e por falta de datos en QA:** En la primera corrida remota en QA, la suite e2e falló por timeout al encontrar el catálogo vacío. Se resolvió ingresando al Panel Admin de QA, creando una categoría y un producto con stock, y utilizando "Re-run failed jobs" en GitHub Actions.
3. **Ajuste de aserciones en e2e:** La redirección hacia WhatsApp requería asertar inmediatamente el mensaje `¡Pedido confirmado correctamente!` y el vaciado del carrito antes del cambio de URL del navegador.
4. **Consumo progresivo de stock en QA y necesidad de Cleanup:** Los tests E2E deben ser idempotentes y dejar el sistema exactamente como lo encontraron. Cada pedido ejecutado por la prueba restaba stock de forma permanente, por lo que en sucesivas corridas el inventario se agotaba inevitablemente. Se diseñó e implementó la funcionalidad completa de cancelación administrativa con restauración de stock transaccional en backend y frontend, incorporando el paso de cleanup en el test.

---

## 11. Declaración de uso de IA

Durante el desarrollo del Trabajo Práctico 07 se utilizó **Antigravity** como asistente de pair programming.
- **Rol de la IA:** Asistencia en la adaptación de la guía oficial (.NET) a nuestro stack (Go + React), diagnóstico de errores de conexión y colisión de runners, diseño de la transacción atómica de restauración de stock, diagnóstico del *strict mode violation* en Playwright y soporte en la redacción técnica.
- **Verificación humana:** Todas las pruebas (unitarias, integración y e2e) fueron ejecutadas, comprobadas y validadas directamente en mi terminal local y en las corridas de GitHub Actions y Render Web.

---

## 12. Idempotencia y Cleanup en Pruebas e2e: Restauración Atómica de Stock y Mejora de Cobertura

### 12.1 Justificación Teórica y Regla de Idempotencia
Una suite de pruebas automatizadas contra un entorno compartido (como QA en Render) no puede ser destructiva con los datos de prueba. Si un test consume recursos no renovables (como stock de inventario) sin restaurarlos:
- La suite se vuelve dependiente del orden de ejecución.
- Genera fallos intermitentes (*flaky tests*) cuando el stock llega a cero.
- Obliga a una intervención manual periódica para "recargar datos", lo cual rompe la autonomía del pipeline de Integración Continua.

### 12.3 Incremento de Cobertura en Backend (del 68.1% al 84.1%)
Al incorporar los tests unitarios con formato **AAA** (Arrange-Act-Assert) y Mocks:
- `CancelOrder` quedó cubierto al **100%** (`TestCancelOrder_Success_LlamaCancelWithStockRestore`, `TestCancelOrder_NotFound`, `TestCancelOrder_InvalidStatus` parametrizado).
- Diagnosticamos que `ProductService` tenía baja cobertura (41.2% en `Create` y 44.0% en `Update`) porque los tests anteriores solo verificaban rechazos tempranos (errores de validación inicial), cortando en las primeras 4 líneas sin recorrer jamás el **camino feliz** (*happy path*).
- Agregamos tests AAA con mocks para la creación y edición exitosa de productos con categorías (`TestCreate_Exitoso_GuardaProducto`, `TestUpdate_Exitoso`, y validaciones de categorías obligatorias y existentes), excluyendo lecturas triviales o tests de infraestructura ajenos a la lógica de negocio.
- **Resultado:** La cobertura en `service` subió al **84.1%** de sentencias (`Create` al 94.1% y `Update` al 96.0%), superando holgadamente el umbral de calidad del 70% exigido.

# Decisiones TP8 – Infraestructura como Código (IaC)

## 1. Riel Elegido y Versiones de Providers Fijadas

- **Riel:** Canónico (Terraform + Neon para PostgreSQL + Render para Web Services).
- **Herramienta:** Terraform (>= 1.5).
- **Entorno aprovisionado:** `preprod` (entorno nuevo, efímero e independiente; QA y PROD no se modificaron ni importaron, permaneciendo intactos).
- **Providers fijados:**
  - `kislerdm/neon`: versión `0.18.0` (fijada explícitamente debido al estado archivado del repositorio comunitario y la espera de la publicación del provider oficial `neondatabase/neon`).
  - `render-oss/render`: versión `1.9.1`.
- **Higiene de estado y credenciales:**
  - El estado (`terraform.tfstate`) vive exclusivamente en la máquina local y está excluido de Git mediante `.gitignore` para evitar la filtración de credenciales sensibles (passwords de Neon en texto plano).
  - Las credenciales (`NEON_API_KEY`, `RENDER_API_KEY`, `RENDER_OWNER_ID`) se inyectan mediante variables de entorno en la terminal y mediante *Repository Secrets* en GitHub Actions, sin estar escritas en ningún archivo `.tf`.
  - El archivo `.terraform.lock.hcl` sí se versiona en Git para asegurar paridad de hashes y versiones entre operadores y CI.

## 1.b Conceptos Clave Verificados (Declarativo, Idempotencia, Grafo y Drift)

- **Modelo Declarativo e Idempotencia:** En `infra/main.tf` declaramos el estado final deseado (los 5 recursos de `preprod`) y no una secuencia imperativa de pasos. Verificamos la **idempotencia** ejecutando un segundo `terraform plan` y `apply` inmediatamente después de crear el entorno, obteniendo `No changes. Your infrastructure matches the configuration.` sin duplicar recursos ni arrojar errores.
- **Grafo de dependencias y orden de destrucción:** Sin usar `depends_on` explícitos, Terraform armó el grafo de dependencias leyendo las referencias entre recursos (`neon_project.preprod.id`, `default_branch_id`, `local.conn` y `render_web_service.api.url`). Aunque el texto del `plan` lista los recursos en **orden alfabético**, al ejecutar `terraform destroy` comprobamos que los elimina recorriendo el **grafo en orden inverso**: primero el Frontend, luego la API, luego la base de datos y el rol, y por último el proyecto de Neon.
- **Drift sobre lo declarado vs. lo no declarado:** Al borrar manualmente la base `app_preprod` en la consola de Neon, `terraform plan` refrescó el estado contra la API, detectó el **drift** y propuso recrearla (`Plan: 1 to add`). En cambio, al crear manualmente una `base_intrusa` no declarada en el HCL (al igual que ocurre con nuestros servicios preexistentes de QA y PROD), `terraform plan` siguió respondiendo `No changes`, demostrando que Terraform solo gobierna los recursos registrados en su `terraform.tfstate`.

## 2. Problemas Encontrados y Cómo los Resolvimos

1. **Límite de retención de historial en Neon (Free Tier):**
   - *Problema:* El provider `kislerdm/neon` solicita por defecto 86400 segundos (24 horas) de retención para Point-in-Time Recovery. La capa gratuita de Neon admite un máximo de 21600 segundos (6 horas). El primer `plan`/`apply` fallaba con `requested history retention seconds exceeds allowed maximum`.
   - *Solución:* Se configuró explícitamente `history_retention_seconds = 21600` en el recurso `neon_project.preprod`.
2. **Adaptación de la cadena de conexión para Go / GORM:**
   - *Problema:* La guía presenta el ejemplo en .NET con `ConnectionStrings__Default` (formato clave-valor con punto y coma). Nuestra aplicación backend está desarrollada en Go y utiliza la variable de entorno `DATABASE_URL` con formato URI de PostgreSQL.
   - *Solución:* Se definió `DATABASE_URL = { value = local.conn }` en `main.tf`, donde `local.conn` se deriva dinámicamente como `postgres://${neon_role.preprod.name}:${neon_role.preprod.password}@${neon_project.preprod.database_host}/${neon_database.preprod.name}?sslmode=require`. No se copiaron credenciales a mano.
3. **Owner ID inválido en Render:**
   - *Problema:* Al ejecutar el primer `apply` contra Render, el provider rechazó el `RENDER_OWNER_ID` configurado inicialmente (`usr-...`) indicando que el ID no era válido para el aprovisionamiento de servicios de ese workspace.
   - *Solución:* Se consultó directamente el endpoint `/v1/owners` de Render mediante `curl -fsS -H "Authorization: Bearer $RENDER_API_KEY" https://api.render.com/v1/owners`, extrayendo el ID exacto del workspace asociado a la clave de API y actualizando la variable de entorno.
4. **Plan del Pipeline vs. Estado Local (Sección 3.6 Opcional):**
   - *Problema:* Al implementar el workflow de GitHub Actions (`infra.yml`), el job `infra-plan` reporta `5 to add` a pesar de que los recursos ya existen en Neon y Render.
   - *Solución:* Se comprendió que el runner de CI arranca sin estado porque `terraform.tfstate` vive en la máquina local y está en `.gitignore`. Para que un pipeline ejecute diffs incrementales o haga `apply` automático, requeriría un *Remote Backend* compartido (ej. S3, Azure Blob o HCP Terraform) y gates de aprobación.
5. **Cuota compartida del Free Tier de Render (750 horas/mes) y estrategia de entornos efímeros:**
   - *Problema:* El plan gratuito de Render comparte 750 horas mensuales entre todos los servicios del workspace. Al sumar preprod pasamos de 4 a 6 servicios, lo que aumentaría el riesgo de suspensión de todo el workspace (incluyendo producción). Además, el provider de Render rechaza modificaciones en servicios gratuitos por el bug de `maintenance_mode`.
   - *Solución:* Se adoptó la estrategia de entornos efímeros: preprod no se modifica "in place" sino que se destruye (`terraform destroy`) cuando no está en uso y se recrea (`terraform apply`) a demanda, garantizando que QA y PROD sigan siempre disponibles para la cursada y la defensa.

## 3. Declaración de Uso de IA

Durante el desarrollo del Trabajo Práctico 08 se utilizó **Antigravity** como asistente de pair programming.
- **Rol de la IA:** 
  - Explicación conceptual y pedagógica de los fundamentos de IaC: modelo declarativo vs. imperativo, idempotencia, el estado como fuente de correlación y no de verdad, drift y su alcance real.
  - Adaptación técnica del HCL al stack Go + GORM (`DATABASE_URL`).
  - Diagnóstico del error de autorización de Render Owner ID mediante la API REST de Render.
  - Diseño e implementación del workflow opcional de GitHub Actions (`infra.yml`) con jobs desacoplados (`infra-check` e `infra-plan`).
- **Verificación humana:** 
  - Todos los comandos (`init`, `fmt`, `validate`, `plan`, `apply`, `destroy`), las pruebas de drift (borrado manual en Neon, prueba de la `base_intrusa`), las comprobaciones de red de seguridad (`prevent_destroy`), las peticiones `curl` a endpoints de vida (`/api/health`) y base de datos (`/api/products`), y la verificación de aislamiento respecto a QA y PROD fueron ejecutadas, comprobadas y validadas directamente por el alumno.

# Decisiones TP9 – DevSecOps + Observabilidad + Continuous Feedback


## Enlaces del TP9
- **PR de Dependabot analizado y mergeado (PR #67):** https://github.com/franciscotaurian/ingsoft3-tp01/pull/67
- **Historial de ejecuciones de CI/CD y CodeQL:** https://github.com/franciscotaurian/ingsoft3-tp01/actions
- **Release de cierre (`v9.0.1`):** https://github.com/franciscotaurian/ingsoft3-tp01/releases/tag/v9.0.1

## 1. Actualización de Dependabot (SCA) y Criterio de Confianza

- **Pull Request mergeado:** PR #67 (`chore(deps-dev): bump @playwright/test from 1.64.0 to 1.64.1 in /App/frontend`).
- **Qué actualizaba:** Actualizó la suite de pruebas e2e `@playwright/test` en el frontend (`App/frontend/package.json`) de la versión `1.64.0` a la `1.64.1`.
- **Análisis SemVer:** Se trató de una actualización de tipo **PATCH** (incremento del tercer dígito). Según Semantic Versioning, un parche contiene exclusivamente correcciones de errores (*bug fixes*), garantizando compatibilidad hacia atrás sin breaking changes ni modificaciones de API.
- **Detalle del changelog:** Corrigió dos problemas puntuales del motor de pruebas:
  1. Un fallo intermitente al capturar la primera ventana en aplicaciones empaquetadas con Electron (`#43247`).
  2. Una regresión en Chromium 155+ donde se abortaban redirecciones *cross-origin* cuando había rutas interceptadas (`#43242`).
- **Por qué me animé a mergearlo:** Al ser una `devDependency` que corre los tests e2e en el pipeline y no código empaquetado en el bundle final que viaja al cliente en producción, no existía riesgo de regresión funcional para el usuario final. Además, todo el pipeline de Integración Continua (build de contenedores, tests unitarios, umbral de cobertura del 70% y suite de integración) corrió sobre la rama de Dependabot y finalizó en **verde**, confirmando empíricamente que la nueva versión del test runner mantiene la estabilidad y compatibilidad de las pruebas existentes.

---

## 2. Observabilidad: Monitor Externo y Telemetría con Sentry

### Monitor Externo (UptimeRobot) y Caída
- **Configuración:** Monitor HTTP(s) externo apuntando a `https://miapp-api-prod-5djk.onrender.com/api/health`, con sondeo cada 5 minutos y notificación vía email.
- **Tiempo de detección y duración de la caída:** El monitor detectó el incidente en su primer ciclo de sondeo (~5 minutos), y la caída total registrada en el historial y en el correo de alerta fue de **35 minutos y 35 segundos**.
- **Causa del incidente real:** El plan gratuito de UptimeRobot realiza peticiones mediante el método HTTP `HEAD` (para optimizar ancho de banda solicitando solo cabeceras). Por defecto, el router Gin de Go no mapea `HEAD` en rutas registradas con `r.GET` si `HandleMethodNotAllowed` está desactivado, respondiendo `404 Not Found`.
- **Mitigación y recuperación:** Se refactorizó la ruta de salud usando `r.Match([]string{"GET", "HEAD"}, "/api/health", ...)` en conjunto con la instrumentación de Sentry. El tiempo de 35 minutos reflejó el ciclo completo de la ingeniería del sistema: detección por el monitor ➔ diagnóstico del método HTTP ➔ desarrollo de la solución ➔ corrida completa del pipeline de CI/CD (QA, smoke tests, integración, e2e, aprobación y deploy en PROD) ➔ siguiente ciclo de sondeo del monitor donde confirmó el `200 OK` y emitió automáticamente la alerta de **UP**. El monitor fue **pausado** inmediatamente tras la prueba para preservar la cuota mensual de 750 horas de Render.

### Telemetría de Producción (Sentry)
- **Operación más lenta:** `GET /api/admin/metrics`, registrando una latencia **P95 de 1.023 segundos** (1023 ms). El percentil P95 indica que el 5% de las solicitudes más lentas tardaron más de un segundo, atribuible a consultas de agregación sobre PostgreSQL sin cacheo. Con este dato se abrió una tarjeta de optimización en el tablero Kanban del TP3 para cerrar el ciclo de Continuous Feedback.
- **Operación más usada:** `GET /api/products`, con un volumen inferior a `0.01 req/min` debido al tráfico acotado de pruebas de la cursada.
- **Error provocado:** Se verificó la captura del `panic` en el endpoint de prueba `/api/falla` en la sección *Issues* de Sentry, reportando archivo, línea exacta y stack trace completo sin tumbar el proceso del servidor gracias al middleware `sentrygin`.

---

## 3. Seguridad Estática (CodeQL) y Protección de Secretos (Push Protection)

- **Push Protection (Capa 1):** Se comprobó el bloqueo en tiempo real intentando pushear a una rama remota un Personal Access Token (PAT) real sin permisos en un archivo de texto. GitHub bloqueó el push de forma preventiva (`remote: error: GH013: Repository rule violations found - GHE004: Push cannot contain secrets`), demostrando que las credenciales con firmas conocidas no ingresan al historial del repositorio. Posteriormente la rama fue eliminada localmente y el token revocado.
- **CodeQL (Capa 3 - SAST):** El escaneo por defecto de CodeQL identificó dos hallazgos principales en el repositorio:
  1. *Workflow permissions missing:* Ausencia de declaración explícita de permisos de menor privilegio en ciertos jobs de GitHub Actions.
  2. *Uso de hash criptográfico inseguro (OWASP A02: Cryptographic Failures):* Identificó el uso de SHA-256 para el almacenamiento de contraseñas de administradores. Se analizó que SHA-256 es una función diseñada para ser rápida en verificación de integridad, pero vulnerable a ataques de fuerza bruta asistidos por hardware (GPU/ASIC) cuando se aplica a contraseñas, recomendándose su migración a algoritmos con factor de costo ajustable como `bcrypt` o `argon2`.

---

## 4. Problemas Encontrados y Soluciones

1. **Rechazo de HEAD por Gin en UptimeRobot:**
   - *Problema:* UptimeRobot marcaba la API de producción como DOWN con código 404 a pesar de que el servicio estaba operativo.
   - *Causa:* Gin aísla las rutas por método HTTP. Un `r.GET` no atiende peticiones `HEAD`.
   - *Solución:* Se utilizó `r.Match([]string{"GET", "HEAD"}, "/api/health", ...)` permitiendo que el endpoint atienda ambos métodos transparentemente.
---

## 5. Declaración de Uso de IA

Durante el desarrollo del Trabajo Práctico 09 se utilizó **Antigravity** como asistente de pair programming.
- **Rol de la IA:**
  - Explicación pedagógica de los conceptos centrales: Shift-Left Security, las 3 capas automáticas (secretos, dependencias y código propio), blackbox monitoring vs. observabilidad interna (logs, métricas y trazas), percentiles P50/P95, y el ciclo de continuous feedback.
  - Asesoramiento en la adaptación del SDK de Sentry para Go (`sentry-go` + `sentry-go/gin`) y React (`@sentry/react`).
  - Diagnóstico del comportamiento de Gin frente al método HTTP HEAD utilizado por UptimeRobot.
  - Redacción del cuerpo de la issue para el tablero Kanban y consolidación de la sección de decisiones.
- **Verificación humana:**
  - Todas las configuraciones en las interfaces web (GitHub Advanced Security, CodeQL, UptimeRobot, Sentry y Render) fueron realizadas y validadas por el alumno.
  - La ejecución de comandos locales, la creación de ramas, la verificación de los bloqueos de push, el análisis del PR de Dependabot #67 y las pruebas en producción fueron ejecutadas y comprobadas directamente por el alumno.
