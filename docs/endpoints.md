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
  "new_password": "MiClaveSegura2024!A",
  "confirmed_password": "MiClaveSegura2024!A"
}
```

**Política de contraseñas**

La nueva contraseña debe cumplir todas las reglas siguientes:

| # | Regla                          | Código en `errors[]`     |
|---|--------------------------------|--------------------------|
| 1 | Mínimo 8 caracteres            | `too_short`              |
| 2 | Al menos una letra minúscula   | `missing_lowercase`      |
| 3 | Al menos una letra mayúscula   | `missing_uppercase`      |
| 4 | Al menos un número             | `missing_digit`          |
| 5 | Al menos un símbolo especial   | `missing_symbol`         |

**Letras y dígitos** siguen el estándar Unicode (`\p{Ll}`, `\p{Lu}`, `\p{Nd}`); puedes usar con normalidad `ñ`, `á`, `é`, `ü`. **Símbolos** son una whitelist ASCII cerrada: `! @ # $ % ^ & * ( ) _ + - = [ ] { } ; ' : " \ | , . < > / ? ~ \``. Esto evita homoglyphs y caracteres ambiguos.

Además, `new_password` debe ser igual a `confirmed_password` (código `password_mismatch`).

Las mismas reglas se aplican en frontend (Zod) y en backend (Go). El servidor LDAP puede imponer reglas adicionales — si las rechaza, el endpoint retorna `422 password-policy`.

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

| Status | `type`                              | Descripción                                                       |
|--------|-------------------------------------|-------------------------------------------------------------------|
| 400    | `.../errors/validation`             | Uno o más campos inválidos. Ver `errors[]` para detalle por campo |
| 404    | `.../errors/not-found`              | Token inexistente                                                 |
| 410    | `.../errors/token-expired`          | Token vencido o ya utilizado                                      |
| 422    | `.../errors/password-policy`        | La nueva contraseña pasó la validación de Go pero LDAP la rechazó |
| 429    | `.../errors/rate-limit`             | Demasiadas solicitudes                                            |
| 500    | `.../errors/internal`               | Error interno o LDAP no disponible                                |

**Formato del error 400 `validation`**

```json
{
  "type": "https://chpass.inf.uct.cl/api/errors/validation",
  "title": "Datos inválidos",
  "status": 400,
  "detail": "Los datos enviados no cumplen los requisitos.",
  "instance": "/api/v1/password-reset/confirm",
  "errors": [
    {
      "field": "new_password",
      "code": "too_short",
      "message": "Debe tener al menos 8 caracteres."
    },
    {
      "field": "new_password",
      "code": "missing_symbol",
      "message": "Debe incluir al menos un símbolo especial."
    },
    {
      "field": "confirmed_password",
      "code": "password_mismatch",
      "message": "Las contraseñas no coinciden."
    }
  ]
}
```

El array `errors` se omite del JSON cuando está vacío (`omitempty`), por lo que las respuestas de error existentes no se ven afectadas.
