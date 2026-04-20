#!/usr/bin/env bash
# Configuración base para los scripts de prueba
# Edita estas variables antes de ejecutar los scripts

BASE_URL="${BASE_URL:-http://localhost:8080}"
EMAIL="${EMAIL:-usuario@ejemplo.cl}"

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

print_header() {
    echo -e "\n${CYAN}==== $1 ====${NC}"
}

print_ok() {
    echo -e "${GREEN}[OK]${NC} $1"
}

print_err() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

check_curl() {
    if ! command -v curl &>/dev/null; then
        echo "curl no está instalado" >&2
        exit 1
    fi
}

check_jq() {
    if ! command -v jq &>/dev/null; then
        print_info "jq no encontrado, el output no será formateado"
        JQ_AVAILABLE=false
    else
        JQ_AVAILABLE=true
    fi
}

pretty_print() {
    if [ "$JQ_AVAILABLE" = true ]; then
        echo "$1" | jq .
    else
        echo "$1"
    fi
}
