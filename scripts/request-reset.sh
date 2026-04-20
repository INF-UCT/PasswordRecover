#!/usr/bin/env bash
# Prueba el endpoint POST /api/v1/password-reset/request
# Uso: ./request-reset.sh [email]
#      EMAIL=correo@dominio.cl ./request-reset.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/config.sh"

check_curl
check_jq

EMAIL="${1:-$EMAIL}"
ENDPOINT="$BASE_URL/api/v1/password-reset/request"

print_header "Solicitar restablecimiento de contraseña"
print_info "URL: $ENDPOINT"
print_info "Email: $EMAIL"

RESPONSE=$(curl -s -w "\n%{http_code}" \
    -X POST \
    -H "Content-Type: application/json" \
    -d "{\"email\": \"$EMAIL\"}" \
    "$ENDPOINT")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n -1)

echo -e "\nHTTP Status: $HTTP_CODE"
echo "Respuesta:"
pretty_print "$BODY"

if [ "$HTTP_CODE" = "200" ]; then
    print_ok "Solicitud enviada correctamente"
else
    print_err "Error en la solicitud (HTTP $HTTP_CODE)"
fi

# Casos de prueba adicionales
if [ "$1" = "--test-all" ]; then
    print_header "Prueba: email inválido (debe retornar 400)"
    RESP=$(curl -s -w "\n%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -d '{"email": "no-es-un-email"}' \
        "$ENDPOINT")
    CODE=$(echo "$RESP" | tail -n1)
    BODY=$(echo "$RESP" | head -n -1)
    echo "HTTP Status: $CODE"
    pretty_print "$BODY"
    [ "$CODE" = "400" ] && print_ok "Validación correcta" || print_err "Se esperaba 400, se obtuvo $CODE"

    print_header "Prueba: payload vacío (debe retornar 400)"
    RESP=$(curl -s -w "\n%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -d '{}' \
        "$ENDPOINT")
    CODE=$(echo "$RESP" | tail -n1)
    BODY=$(echo "$RESP" | head -n -1)
    echo "HTTP Status: $CODE"
    pretty_print "$BODY"
    [ "$CODE" = "400" ] && print_ok "Validación correcta" || print_err "Se esperaba 400, se obtuvo $CODE"
fi
