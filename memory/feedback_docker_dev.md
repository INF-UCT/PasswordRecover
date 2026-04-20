---
name: Usar Docker para desarrollo Node/JS
description: El usuario no tiene Node.js instalado localmente, usar contenedores Docker para todo lo relacionado con frontend/JS
type: feedback
---

Usar Docker (o docker compose) para cualquier tarea de desarrollo que requiera Node.js, npm, o herramientas de frontend (Svelte, Vite, etc.).

**Why:** El usuario no tiene un entorno de desarrollo Node instalado en su máquina.

**How to apply:** Al crear o ejecutar proyectos frontend, incluir siempre un Dockerfile o servicio en docker-compose para el dev server. No asumir que `node`, `npm`, `npx` están disponibles en el host.
