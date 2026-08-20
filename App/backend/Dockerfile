# ---- Etapa 1: build ----
FROM golang:1.26-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/binario-backend main.go

# ---- Etapa 2: runtime ----
FROM alpine:latest AS final
WORKDIR /app

COPY --from=build /app/binario-backend .

EXPOSE 8080
ENTRYPOINT ["./binario-backend"]