# Política de contraseñas

Este documento describe las reglas que debe cumplir la nueva contraseña al usar el endpoint `POST /api/v1/password-reset/confirm`, y cómo se aplican en cada capa del sistema.

---

## Reglas

Toda contraseña nueva debe cumplir **todas** las reglas siguientes:

| # | Regla                          | Código de error (backend)  |
|---|--------------------------------|----------------------------|
| 1 | Mínimo 8 caracteres            | `too_short`                |
| 2 | Al menos una letra minúscula   | `missing_lowercase`        |
| 3 | Al menos una letra mayúscula   | `missing_uppercase`        |
| 4 | Al menos un número             | `missing_digit`            |
| 5 | Al menos un símbolo especial   | `missing_symbol`           |

Además, `new_password` debe ser idéntico a `confirmed_password`. Si no coinciden, el error es:

| Campo               | Código              |
|---------------------|---------------------|
| `confirmed_password`| `password_mismatch` |

### Símbolos especiales aceptados

```
!  @  #  $  %  ^  &  *  (  )  _  +  -  =  [  ]  {  }  ;  '  :  "  \  |  ,  .  <  >  /  ?  ~  `
```

> **Nota:** la lista de símbolos es deliberadamente cerrada (whitelist ASCII) para evitar edge cases: espacios, comillas tipográficas y caracteres unicode visualmente similares (homoglyphs) **no** cuentan como símbolo válido.

### Letras y dígitos

Las letras y dígitos siguen el estándar Unicode (`\p{Ll}`, `\p{Lu}`, `\p{Nd}`). Puedes usar con total normalidad caracteres como `ñ`, `Ñ`, `á`, `é`, `í`, `ó`, `ú`, `ü` en tu contraseña. Los símbolos, en cambio, son ASCII fijos (ver sección anterior) para evitar homoglyphs.

---

## Capas de validación (defensa en profundidad)

La contraseña se valida en tres capas, en este orden:

### 1. Frontend — Zod (`frontend/src/lib/validators.js`)
- Validación inmediata mientras el usuario tipea, sin esperar al servidor.
- Muestra un panel de "requisitos" con cinco checks (✓/○) que se actualizan reactivamente.
- Bloquea el envío si la contraseña no cumple las reglas o si las contraseñas no coinciden.
- Mensajes claros y granulares en español, uno por regla incumplida.
- `safeParse` se usa para no lanzar excepciones; el resultado se mapea a `{ new_password?, confirmed_password?, token? }`.

### 2. Backend — Go (`internal/handler/validate.go`)
- Re-validación en el servidor: la API no confía en el cliente.
- Si la validación falla, responde `400 Bad Request` con `type: .../errors/validation` y un array `errors[]` con un objeto por cada regla incumplida (`{ field, code, message }`).
- Si la contraseña cumple las reglas de Go pero LDAP la rechaza, el servidor responde `422 Unprocessable Entity` con `type: .../errors/password-policy` y el mensaje de LDAP como `detail`.

### 3. LDAP — servidor externo
- Política de contraseñas definitiva. Si LDAP tiene reglas adicionales (historial, longitud máxima, etc.) y la contraseña las viola, retorna código 19 (`LDAPResultConstraintViolation`), que el backend mapea a `ErrPasswordPolicy` → `422 password-policy`.

---

## Formato del error 400 `validation`

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

El array `errors` se omite del JSON cuando está vacío (`omitempty`), por lo que las respuestas de error existentes en otros endpoints no se ven afectadas.

---

## Ejemplos

**Válidas** (cumplen las 5 reglas):

- `MiClave2024!`
- `Cambio#Seguro9`
- `A1b2C3d4$`

**Inválidas** (con el motivo):

| Contraseña          | Reglas que falla                                  |
|---------------------|---------------------------------------------------|
| `miclave2024`       | muy corta, sin mayúscula, sin símbolo             |
| `Miclave2024`       | sin símbolo                                       |
| `Miclave!`          | muy corta, sin número                             |
| `MICLAVE2024!`      | sin minúscula                                     |
| `MiClaveSegura`     | sin número, sin símbolo                           |
| `MiClave2024!` vs `MiClave2024` | `password_mismatch` (no coinciden)         |
