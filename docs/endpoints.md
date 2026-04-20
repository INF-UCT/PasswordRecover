# Diseño de Endpoints

Este documento describe los endpoints disponibles en la API de cambio de contraseña, incluyendo los datos que reciben y las respuestas que retornan.

---

## POST /api/v1/password-reset/request

Recibe el correo del usuario y, si existe en LDAP, envía un enlace de restablecimiento al correo indicado.

**Request**

```json
{
  "email": "juan@dominio.cl"
}
```

**Response `200 OK`**

Siempre retorna `200` independientemente de si el correo existe o no, para evitar enumeración de usuarios.

```json
{
  "success": true,
  "meta": {
    "request_id": "1d5e-c50...",
    "timestamp": "2025-05-13T16:04:55Z"
  },
  "data": {
    "message": "Si el correo existe en el sistema, recibirás un enlace para cambiar tu contraseña."
  },
  "errors": null
}
```

**Errores posibles**

| Status | `type`                              | Descripción                        |
|--------|-------------------------------------|------------------------------------|
| 400    | `.../errors/invalid-input`          | El campo `email` falta o es inválido |
| 429    | `.../errors/rate-limit`             | Demasiadas solicitudes desde la misma IP o para el mismo correo |
| 500    | `.../errors/internal`               | Error interno o LDAP no disponible |

---

## GET /api/v1/password-reset/{token}

Valida si un token de restablecimiento sigue siendo válido. El frontend puede llamar a este endpoint antes de mostrar el formulario de nueva contraseña, evitando que el usuario complete el formulario con un token ya vencido.

**Path parameter**

| Parámetro | Tipo   | Descripción       |
|-----------|--------|-------------------|
| `token`   | string | Token de reseteo  |

**Response `200 OK`**

```json
{
  "success": true,
  "meta": {
    "request_id": "1d5e-c50...",
    "timestamp": "2025-05-13T16:04:55Z"
  },
  "data": {
    "valid": true,
    "expires_at": "2025-05-13T16:14:55Z"
  },
  "errors": null
}
```

**Errores posibles**

| Status | `type`                              | Descripción                  |
|--------|-------------------------------------|------------------------------|
| 404    | `.../errors/not-found`              | Token inexistente             |
| 410    | `.../errors/token-expired`          | Token vencido o ya utilizado  |
| 500    | `.../errors/internal`               | Error interno                 |

---

## POST /api/v1/password-reset/confirm

Recibe el token y la nueva contraseña para efectuar el cambio en LDAP. El token se invalida inmediatamente tras un cambio exitoso.

**Request**

```json
{
  "token": "fjdksalio389vn98wehfv9dhdfs9s",
  "new_password": "NuevaContraseña123!",
  "confirmed_password": "NuevaContraseña123!"
}
```

**Response `200 OK`**

```json
{
  "success": true,
  "meta": {
    "request_id": "1d5e-c50...",
    "timestamp": "2025-05-13T16:05:30Z"
  },
  "data": {
    "message": "Contraseña actualizada correctamente."
  },
  "errors": null
}
```

**Errores posibles**

| Status | `type`                              | Descripción                                          |
|--------|-------------------------------------|------------------------------------------------------|
| 400    | `.../errors/invalid-input`          | Payload malformado o `new_password` != `confirmed_password` |
| 404    | `.../errors/not-found`              | Token inexistente                                    |
| 410    | `.../errors/token-expired`          | Token vencido o ya utilizado                         |
| 422    | `.../errors/password-policy`        | La nueva contraseña no cumple la política de LDAP    |
| 429    | `.../errors/rate-limit`             | Demasiadas solicitudes                               |
| 500    | `.../errors/internal`               | Error interno o LDAP no disponible                   |
