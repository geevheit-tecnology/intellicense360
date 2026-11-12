#!/usr/bin/env sh
set -eu

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"

echo "🚀 Iniciando Geevheit Intelligence 360..."
echo "📦 1. Backend Go API iniciando em http://localhost:8080"
echo "🖥️  2. Frontend Mission Control iniciando em http://localhost:3000"

cleanup() {
    echo "\n🛑 Encerrando serviços..."
    kill "$BACKEND_PID" "$FRONTEND_PID" 2>/dev/null || true
    wait "$BACKEND_PID" 2>/dev/null || true
    wait "$FRONTEND_PID" 2>/dev/null || true
    exit 0
}
trap cleanup INT TERM

(
    cd "$ROOT_DIR/backend/api"
    STORAGE_DRIVER="${STORAGE_DRIVER:-memory}" go run ./cmd/api
) &
BACKEND_PID=$!

(
    cd "$ROOT_DIR/apps/mission-control"
    npm run dev
) &
FRONTEND_PID=$!

echo "✅ Todos os serviços estão ativos! Acesse: http://localhost:3000"
echo "👤 Usuário ADM / Gestor: admin@geevheit.local | Senha: admin1234"

wait

