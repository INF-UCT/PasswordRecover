#!/usr/bin/env bash
# Prueba el endpoint GET /api/v1/password-reset/{token}
# Uso: ./validate-token.sh <token>

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/config.sh"

check_curl
check_jq

TOKEN="${1}"

if [ -z "$TOKEN" ]; then
    echo "Uso: $0 <token>"
    echo "Ejemplo: $0 550e8400-e29b-41d4-a716-446655440000"
    exit 1
fi

ENDPOINT="$BASE_URL/api/v1/password-reset/$TOKEN"

print_header "Validar token de restablecimiento"
print_info "URL: $ENDPOINT"
print_info "Token: $TOKEN"

RESPONSE=$(curl -s -w "\n%{http_code}" \
    -X GET \
    "$ENDPOINT")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n -1)

echo -e "\nHTTP Status: $HTTP_CODE"
echo "Respuesta:"
pretty_print "$BODY"

case "$HTTP_CODE" in
    200) print_ok "Token válido" ;;
    404) print_err "Token no encontrado" ;;
    410) print_err "Token vencido o ya utilizado" ;;
    500) print_err "Error interno del servidor" ;;
    *)   print_err "Respuesta inesperada (HTTP $HTTP_CODE)" ;;
esac

# Prueba con token inexistente
if [ "$1" = "--test-invalid" ]; then
    print_header "Prueba: token inexistente (debe retornar 404)"
    FAKE_TOKEN="00000000-0000-0000-0000-000000000000"
    RESP=$(curl -s -w "\n%{http_code}" \
        -X GET \
        "$BASE_URL/api/v1/password-reset/$FAKE_TOKEN")
    CODE=$(echo "$RESP" | tail -n1)
    BODY=$(echo "$RESP" | head -n -1)
    echo "HTTP Status: $CODE"
    pretty_print "$BODY"
    [ "$CODE" = "404" ] && print_ok "Comportamiento correcto" || print_err "Se esperaba 404, se obtuvo $CODE"
fi
