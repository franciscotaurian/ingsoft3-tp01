# ---- Etapa 1: build (imagen con Node.js) ----
FROM node:20-alpine AS build
WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# ---- Etapa 2: runtime (imagen con Nginx) ----
FROM nginx:alpine AS final
WORKDIR /usr/share/nginx/html

# Borramos los archivos por defecto de Nginx
RUN rm -rf ./*

# Copiamos la carpeta 'dist' compilada en la etapa 1
COPY --from=build /app/dist .

# Copiamos nuestra configuración personalizada de Nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
ENTRYPOINT ["nginx", "-g", "daemon off;"]