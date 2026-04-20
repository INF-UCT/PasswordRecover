---
name: Preferir nginx para reverse proxy
description: El usuario prefiere nginx como reverse proxy en lugar de soluciones a nivel de aplicación
type: feedback
---

Usar nginx para reverse proxy en vez de hooks/middleware de la aplicación (como SvelteKit hooks.server.js, Express middleware, etc.).

**Why:** Preferencia explícita del usuario por nginx como capa de proxy.

**How to apply:** Al necesitar proxy inverso entre servicios, crear un servicio nginx con su configuración en docker-compose en vez de implementar el proxy dentro del código de la aplicación.
