# Tecnologías

Este documento lista las tecnologías y librerías seleccionadas para el desarrollo del servicio de cambio de contraseña.

## Lenguaje

| Tecnología | Versión | Descripción |
|------------|---------|-------------|
| Go | latest stable | Lenguaje principal del servicio. Compilado, de alto rendimiento y con soporte nativo para concurrencia. |

## Base de datos

| Tecnología | Descripción |
|------------|-------------|
| Redis | Almacenamiento en memoria para los UUID de restablecimiento. Los tokens se almacenan con TTL de 10 minutos y se eliminan automáticamente al vencer o al ser utilizados. |

## Librerías

| Librería | Propósito |
|----------|-----------|
| `github.com/wneessen/go-mail` | Envío de correos via SMTP. Soporta TLS/STARTTLS y contenido HTML. |
| `github.com/go-ldap/ldap/v3` | Conexión y operaciones contra el servidor LDAP (búsqueda de usuarios y actualización de contraseña). |
| `github.com/redis/go-redis/v9` | Cliente de Redis. Soporta Redis Cluster, compatible con despliegues en Kubernetes. |
