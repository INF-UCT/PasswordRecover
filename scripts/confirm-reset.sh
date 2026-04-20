#!/usr/bin/env bash
# Prueba el endpoint POST /api/v1/password-reset/confirm
# Uso: ./confirm-reset.sh <token> <nueva_contraseña>

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/config.sh"

check_curl
check_jq

TOKEN="${1}"
NEW_PASSWORD="${2}"

if [ -z "$TOKEN" ] || [ -z "$NEW_PASSWORD" ]; then
    echo "Uso: $0 <token> <nueva_contraseña>"
    echo "Ejemplo: $0 550e8400-e29b-41d4-a716-446655440000 'MiNuevaPass123!'"
    exit 1
fi

ENDPOINT="$BASE_URL/api/v1/password-reset/confirm"

print_header "Confirmar cambio de contraseña"
print_info "URL: $ENDPOINT"
print_info "Token: $TOKEN"

PAYLOAD=$(jq -n \
    --arg token "$TOKEN" \
    --arg pass "$NEW_PASSWORD" \
    '{token: $token, new_password: $pass, confirmed_password: $pass}')

RESPONSE=$(curl -s -w "\n%{http_code}" \
    -X POST \
    -H "Content-Type: application/json" \
    -d "$PAYLOAD" \
    "$ENDPOINT")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n -1)

echo -e "\nHTTP Status: $HTTP_CODE"
echo "Respuesta:"
pretty_print "$BODY"

case "$HTTP_CODE" in
    200) print_ok "Contraseña actualizada correctamente" ;;
    400) print_err "Payload inválido o contraseñas no coinciden" ;;
    404) print_err "Token no encontrado" ;;
    410) print_err "Token vencido o ya utilizado" ;;
    422) print_err "La contraseña no cumple la política de LDAP" ;;
    429) print_err "Demasiadas solicitudes (rate limit)" ;;
    500) print_err "Error interno del servidor" ;;
    *)   print_err "Respuesta inesperada (HTTP $HTTP_CODE)" ;;
esac

# Prueba: contraseñas que no coinciden
if [ "$1" = "--test-mismatch" ]; then
    print_header "Prueba: contraseñas no coinciden (debe retornar 400)"
    RESP=$(curl -s -w "\n%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -d '{"token": "fake-token", "new_password": "Pass1!", "confirmed_password": "Pass2!"}' \
        "$ENDPOINT")
    CODE=$(echo "$RESP" | tail -n1)
    BODY=$(echo "$RESP" | head -n -1)
    echo "HTTP Status: $CODE"
    pretty_print "$BODY"
    [ "$CODE" = "400" ] && print_ok "Validación correcta" || print_err "Se esperaba 400, se obtuvo $CODE"
fi
