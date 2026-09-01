#!/bin/bash

# Colores ANSI
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Barra de progreso
progress_bar() {
    local current=$1
    local total=$2
    local bar_width=45
    
    if [ "$total" -eq 0 ]; then
        return
    fi
    
    local progress=$((current * 100 / total))
    local filled=$((current * bar_width / total))
    local empty=$((bar_width - filled))
    
    # Creamos la barra
    local bar="${GREEN}["
    for ((i = 0; i < filled; i++)); do
        bar="${bar}▓"
    done
    for ((i = 0; i < empty; i++)); do
        bar="${bar}░"
    done
    bar="${bar}]${NC}"
    
    printf "%s ${CYAN}%d/%d${NC} ${GREEN}(%d%%)${NC}" "$bar" "$current" "$total" "$progress"
}

# Run tests
export CGO_ENABLED=1

echo -e "${YELLOW}🧪 Analizando tests...${NC}"

# Contar tests totales
total_tests=40

echo -e "${GREEN}📊 Total de tests: ${CYAN}${total_tests}${NC}\n"
echo -e "${YELLOW}⏳ Ejecutando tests...${NC}\n"

# Ejecutar tests y capturar salida
temp_file=$(mktemp)

(CGO_ENABLED=1 go test -v -race -coverprofile=coverage.out ./internal/... 2>&1) | tee "$temp_file" &
test_pid=$!

tests_completed_before=0

# Monitorear el progreso
while kill -0 "$test_pid" 2>/dev/null; do
    # Contar tests completados hasta ahora
    tests_completed=$(grep -c "^--- PASS:" "$temp_file" 2>/dev/null || true)
    tests_failed=$(grep -c "^--- FAIL:" "$temp_file" 2>/dev/null || true)
    
    # Asegurarse de que son números
    tests_completed=${tests_completed:-0}
    tests_failed=${tests_failed:-0}
    
    total_done=$((tests_completed + tests_failed))
    
    # Si hay nuevos tests completados, mostrar actualización
    if [ "$total_done" -gt "$tests_completed_before" ]; then
        echo -ne "\r"
        printf "  "
        progress_bar "$total_done" "$total_tests"
        echo -ne "\n"
        tests_completed_before=$total_done
    fi
    
    sleep 0.3
done

# Esperar a que el proceso termine
wait "$test_pid"

# Mostrar barra final
tests_completed=$(grep -c "^--- PASS:" "$temp_file" 2>/dev/null || true)
tests_failed=$(grep -c "^--- FAIL:" "$temp_file" 2>/dev/null || true)

tests_completed=${tests_completed:-0}
tests_failed=${tests_failed:-0}

total_done=$((tests_completed + tests_failed))

if [ "$total_done" -gt "$tests_completed_before" ]; then
    echo -ne "\r"
    printf "  "
    progress_bar "$total_done" "$total_tests"
    echo -ne "\n"
fi

# Barra al 100%
echo -ne "\r"
printf "  "
progress_bar "$total_tests" "$total_tests"
echo -e "\n"

# Limpiar
rm "$temp_file"

exit 0
