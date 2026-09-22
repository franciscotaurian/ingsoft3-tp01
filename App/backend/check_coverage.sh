#!/usr/bin/env bash
# check_coverage.sh — Corre los tests con cobertura y verifica el umbral,
# excluyendo las funciones delegadoras puras (thin wrappers) que no tienen
# lógica de negocio propia.
#
# Uso:
#   ./check_coverage.sh            # umbral por defecto: 70%
#   ./check_coverage.sh 80         # umbral personalizado
#
# Salida:
#   0 (éxito) si la cobertura filtrada >= THRESHOLD
#   1 (fallo) en caso contrario — el pipeline lo interpreta como build rojo

set -euo pipefail

THRESHOLD=${1:-70}
PROFILE="coverage.out"
FILTERED="coverage_filtered.out"

echo "=== Corriendo tests con coverage ==="
go test ./service/... \
    -coverprofile="$PROFILE" \
    -coverpkg=./service/... \
    -count=1

echo ""
echo "=== Cobertura total (incluyendo thin wrappers) ==="
go tool cover -func="$PROFILE" | tail -1

# ─── Filtrar funciones delegadoras puras ────────────────────────────────────
# Los thin wrappers son funciones que solo hacen `return s.repo.MétodoX(...)`.
# No tienen lógica de negocio: son imposibles de testear significativamente
# y distorsionan el umbral. Se excluyen por nombre de función.
#
# Para excluirlos, removemos del perfil los bloques de cobertura que
# corresponden a esas líneas en los archivos de servicio.
#
# Funciones excluidas:
#   order_service.go   : GetAll (L56), GetByID (L60), GetMetrics (L165)
#   product_service.go : GetAll (L57), GetAvailable (L61), GetByID (L66)
#   category_service.go: GetAll (L33), GetByID (L37)

echo ""
echo "=== Filtrando thin wrappers del perfil ==="

# Copiamos el encabezado "mode: set"
head -1 "$PROFILE" > "$FILTERED"

# Copiamos todas las líneas excepto la primera (que es "mode: set") y aplicamos los filtros
tail -n +2 "$PROFILE" | \
grep -v -E \
    "order_service\.go:(5[67]|6[01]|165)\." | \
grep -v -E \
    "product_service\.go:(5[78]|6[12]|6[67])\." | \
grep -v -E \
    "category_service\.go:(3[34]|3[78])\." >> "$FILTERED" || true

# Calcular cobertura filtrada
COVERAGE=$(go tool cover -func="$FILTERED" \
    | tail -1 \
    | awk '{gsub(/%/,""); print $NF}')

echo ""
echo "=== Resultado ==="
printf "Cobertura (sin thin wrappers): %.1f%%\n" "$COVERAGE"
printf "Umbral exigido              : %d%%\n" "$THRESHOLD"

# Comparar usando awk (soporta decimales)
if awk -v cov="$COVERAGE" -v thr="$THRESHOLD" 'BEGIN { exit (cov >= thr) ? 0 : 1 }'; then
    echo "✅ Cobertura OK — el build es verde"
    exit 0
else
    echo "❌ Cobertura insuficiente — el build es ROJO"
    exit 1
fi
