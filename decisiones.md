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

