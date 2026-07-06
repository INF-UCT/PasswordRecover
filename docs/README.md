# Documento de diseño API para cambiar contraseña

Este documento tiene como objetivo planificar y modelar el funcionamiento de un servicio API REST de cambio de contraseña para usuarios almacenados en un servicio de LDAP

## Funcionamiento
Un usuario hace una petición HTTP POST con su correo como payload, la API recibe esa información verifica que el correo existe en LDAP, para proceder a generar un UUID v4 el cual se utilizara como ruta dentro del enlace que el usuario utilizara para cambiar su contraseña, el enlace generado se le enviara al usuario al correo que envió al servidor en primera instancia.

Los enlaces duran un máximo de 10 minutos

## Estructuras de datos
Respuestas de la API:
```json
{
  "success": true,
  "meta": {
    "request_id": "1d5e-c50...",
    "timestamp": "2025-05-13T16:04:55Z"
  },
  "data": { /* payload */ } || null,
  "errors": null
} 
```

Errors de la API:
```json
{
  "type": "https://chpass.inf.uct.cl/api/errors/unauthorize",
  "title": "Acceso denegado",
  "status": 401,
  "detail": "Token invalido o caducado",
  "instance": "/api/chpass/klfdjskl3893289429fdjfi93j9f2j94jf893f3j9
}
```

Payload que envia el usuario en primera instancia:
```json
{
  "email": "juan@dominio.cl"
}
```

Payload que envia el usuario desde el link para cambiar su contraseña
```json
{
  "token": "550e8400-e29b-41d4-a716-446655440000",
  "new_password": "MiClaveSegura2024!A",
  "confirmed_password": "MiClaveSegura2024!A"
}
```

La nueva contraseña debe cumplir la **política de contraseñas** del sistema:
- Mínimo 8 caracteres.
- Al menos una letra minúscula, una mayúscula y un número.
- Al menos un símbolo especial (`! @ # $ % ^ & * ( ) _ + - = [ ] { } ; ' : " \ | , . < > / ? ~ \``).
- `new_password` y `confirmed_password` deben coincidir.

Estas reglas se aplican en el frontend (Zod), en el backend (Go) y finalmente en LDAP. Ver `docs/politica-contrasenas.md` para el detalle completo.
