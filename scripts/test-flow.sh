#!/usr/bin/env bash
# Flujo completo de prueba interactivo
# Prueba los 3 endpoints en secuencia guiada
# Uso: ./test-flow.sh [email] [nueva_contraseña]

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/config.sh"

check_curl
check_jq

EMAIL="${1:-$EMAIL}"
NEW_PASSWORD="${2:-TestPassword123!}"

echo -e "${CYAN}"
echo "╔══════════════════════════════════════════════╗"
echo "║   Test de flujo completo - chpassword API    ║"
echo "╚══════════════════════════════════════════════╝"
echo -e "${NC}"
print_info "Base URL: $BASE_URL"
print_info "Email:    $EMAIL"

# ─── Paso 1: solicitar reset ───────────────────────────────────────────────
print_header "Paso 1/3 — Solicitar restablecimiento"

RESP=$(curl -s -w "\n%{http_code}" \
    -X POST \
    -H "Content-Type: application/json" \
    -d "{\"email\": \"$EMAIL\"}" \
    "$BASE_URL/api/v1/password-reset/request")

CODE=$(echo "$RESP" | tail -n1)
BODY=$(echo "$RESP" | head -n -1)

pretty_print "$BODY"

if [ "$CODE" != "200" ]; then
    print_err "Fallo en paso 1 (HTTP $CODE). Abortando."
    exit 1
fi
print_ok "Solicitud enviada (HTTP $CODE)"

# ─── Paso 2: validar token ────────────────────────────────────────────────
print_header "Paso 2/3 — Validar token"
echo ""
echo -e "${YELLOW}Ingresa el token recibido en el correo:${NC}"
read -r TOKEN

if [ -z "$TOKEN" ]; then
    print_err "Token vacío. Abortando."
    exit 1
fi

RESP=$(curl -s -w "\n%{http_code}" \
    -X GET \
    "$BASE_URL/api/v1/password-reset/$TOKEN")

CODE=$(echo "$RESP" | tail -n1)
BODY=$(echo "$RESP" | head -n -1)

pretty_print "$BODY"

case "$CODE" in
    200) print_ok "Token válido" ;;
    404) print_err "Token no encontrado. Verifica el valor ingresado."; exit 1 ;;
    410) print_err "Token vencido o ya utilizado."; exit 1 ;;
    *)   print_err "Error inesperado (HTTP $CODE)"; exit 1 ;;
esac

# ─── Paso 3: confirmar nueva contraseña ──────────────────────────────────
print_header "Paso 3/3 — Confirmar nueva contraseña"
print_info "Nueva contraseña: $NEW_PASSWORD"

PAYLOAD=$(jq -n \
    --arg token "$TOKEN" \
    --arg pass "$NEW_PASSWORD" \
    '{token: $token, new_password: $pass, confirmed_password: $pass}')

RESP=$(curl -s -w "\n%{http_code}" \
    -X POST \
    -H "Content-Type: application/json" \
    -d "$PAYLOAD" \
    "$BASE_URL/api/v1/password-reset/confirm")

CODE=$(echo "$RESP" | tail -n1)
BODY=$(echo "$RESP" | head -n -1)

pretty_print "$BODY"

case "$CODE" in
    200) print_ok "¡Contraseña actualizada correctamente!" ;;
    400) print_err "Payload inválido (HTTP $CODE)" ;;
    410) print_err "Token vencido entre pasos (HTTP $CODE)" ;;
    422) print_err "Contraseña no cumple política LDAP (HTTP $CODE)" ;;
    *)   print_err "Error inesperado (HTTP $CODE)" ;;
esac

echo ""
