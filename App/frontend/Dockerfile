# ---- Etapa 1: build (imagen con Node.js) ----
FROM node:20-alpine AS build
WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# ---- Etapa 2: tests ----
FROM build AS test
ENTRYPOINT ["npm", "run", "test:ci"]

# ---- Etapa 3: runtime (imagen con Nginx) ----
FROM nginx:alpine AS final
WORKDIR /usr/share/nginx/html

# Borramos los archivos por defecto de Nginx
RUN rm -rf ./*

# Copiamos la carpeta 'dist' compilada en la etapa 1
COPY --from=build /app/dist .

# Copiamos como plantilla: nginx:alpine aplica envsubst al iniciar, reemplazando ${BACKEND_URL}
COPY nginx.conf /etc/nginx/templates/default.conf.template

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]