# Variables de entorno

Este documento lista todas las variables de entorno necesarias para ejecutar el servicio.

## Aplicación

| Variable | Requerida | Descripción | Ejemplo |
|----------|-----------|-------------|---------|
| `APP_PORT` | No | Puerto en el que escucha el servidor HTTP. Por defecto `8080`. | `8080` |
| `APP_BASE_URL` | Sí | URL base del servicio backend. | `https://chpass.inf.uct.cl` |
| `APP_FRONTEND_URL` | Sí | URL base del frontend, usada para construir el enlace de restablecimiento que se envía al correo. | `https://chpass.inf.uct.cl` |
| `APP_TOKEN_TTL_MINUTES` | No | Tiempo de vida del token en minutos. Por defecto `10`. | `10` |

## LDAP

| Variable | Requerida | Descripción | Ejemplo |
|----------|-----------|-------------|---------|
| `LDAP_HOST` | Sí | Host del servidor LDAP. | `ldap.inf.uct.cl` |
| `LDAP_PORT` | No | Puerto del servidor LDAP. Por defecto `389`. | `389` |
| `LDAP_USE_TLS` | No | Habilita TLS (LDAPS). Por defecto `false`. | `true` |
| `LDAP_BIND_DN` | Sí | DN de la cuenta de servicio usada para autenticarse en LDAP. | `cn=admin,dc=inf,dc=uct,dc=cl` |
| `LDAP_BIND_PASSWORD` | Sí | Contraseña de la cuenta de servicio. | `secret` |
| `LDAP_BASE_DN` | Sí | DN base desde donde se realizan las búsquedas de usuarios. | `dc=inf,dc=uct,dc=cl` |

## Redis

| Variable | Requerida | Descripción | Ejemplo |
|----------|-----------|-------------|---------|
| `REDIS_URL` | Sí | URL de conexión a Redis. | `redis://redis:6379` |
| `REDIS_PASSWORD` | No | Contraseña de Redis, si está configurada. | `secret` |

## SMTP

| Variable | Requerida | Descripción | Ejemplo |
|----------|-----------|-------------|---------|
| `SMTP_HOST` | Sí | Host del servidor SMTP. | `smtp.inf.uct.cl` |
| `SMTP_PORT` | No | Puerto del servidor SMTP. Por defecto `587`. | `587` |
| `SMTP_USER` | Sí | Usuario para autenticarse en el servidor SMTP. | `noreply@inf.uct.cl` |
| `SMTP_PASSWORD` | Sí | Contraseña del usuario SMTP. | `secret` |
| `SMTP_FROM` | Sí | Dirección que aparecerá como remitente en los correos enviados. | `noreply@inf.uct.cl` |
