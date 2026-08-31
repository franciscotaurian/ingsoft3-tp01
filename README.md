# Ingenieria de Software III - Realico Comidas - Catálogo y Pedidos

[![CI](https://github.com/franciscotaurian/ingsoft3-tp01/actions/workflows/ci.yml/badge.svg)](https://github.com/franciscotaurian/ingsoft3-tp01/actions/workflows/ci.yml)

## Instalación

```bash
git clone https://github.com/franciscotaurian/ingsoft3-tp01
```

## Propósito del Proyecto

Este repositorio aloja el desarrollo práctico integral para la materia **Ingeniería de Software III**. Su objetivo principal no es solo presentar un producto de software funcional, sino servir como un entorno de laboratorio para implementar, iterar y madurar prácticas de **Cultura DevOps e Integración/Despliegue Continuo (CI/CD)**.

A lo largo de este repositorio se refleja la evolución desde el código fuente hasta una arquitectura lista para producción, aplicando metodologías ágiles, contenerización, pipelines como código, testing automatizado e Infraestructura como Código (IaC).

## Sobre la Aplicación

**Realico Comidas** es una aplicación web full-stack (Go + React + PostgreSQL) diseñada para modernizar la gestión operativa de un restaurante. Está construida bajo una arquitectura de tres capas con reglas de negocio estrictas (como validación transaccional de stock y congelamiento de precios). 

La plataforma resuelve dos necesidades principales:
*   **Experiencia del Cliente:** Permite a los usuarios explorar un catálogo dinámico, armar un carrito de compras respetando el stock disponible y concretar sus pedidos directamente a través de WhatsApp.
*   **Gestión Administrativa:** Provee un panel privado para que el restaurante administre sus categorías, controle el inventario de productos, visualice métricas y gestione el ciclo de vida de los pedidos (Pendiente ➔ Confirmado ➔ Entregado).

## Ejes Temáticos Abordados

El proyecto implementa de manera progresiva las siguientes prácticas del ciclo de vida del software:

*   **Contenedores:** Dockerización de servicios y orquestación local con Docker Compose.
*   **Continuous Integration (CI):** Pipelines as code para automatizar compilaciones.
*   **Testing y Calidad:** Pruebas unitarias, análisis estático de código, coverage y pruebas End-to-End (E2E) dentro del pipeline.
*   **Continuous Deployment (CD):** Gestión de entornos, despliegue de imágenes mediante GitHub Container Registry (GHCR) y deployment patterns.
*   **Operaciones:** Observabilidad, DevSecOps y Continuous Feedback.



