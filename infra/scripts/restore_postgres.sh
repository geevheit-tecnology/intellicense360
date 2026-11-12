#!/usr/bin/env bash
# ==============================================================================
# Geevheit Intelligence 360 — Script de Restauração Segura do PostgreSQL
# Uso: ./restore_postgres.sh /caminho/do/arquivo.sql.gz
# ==============================================================================

set -euo pipefail

if [ "$#" -ne 1 ]; then
    echo "Uso: $0 <caminho_do_backup.sql.gz>"
    exit 1
fi

BACKUP_FILE="$1"
DB_HOST="${DB_HOST:-postgres}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${POSTGRES_DB:-geevheit_intelligence}"
DB_USER="${POSTGRES_USER:-geevheit}"

if [ ! -f "${BACKUP_FILE}" ]; then
    echo "[ERRO] Arquivo de backup não encontrado: ${BACKUP_FILE}"
    exit 1
fi

echo "=============================================================================="
echo "ATENÇÃO: A restauração irá sobrescrever os dados da base '${DB_NAME}'!"
echo "Arquivo de origem: ${BACKUP_FILE}"
echo "=============================================================================="

echo "[$(date -u +'%Y-%m-%dT%H:%M:%SZ')] [INFO] Descompactando e restaurando base de dados..."
gunzip -c "${BACKUP_FILE}" | PGPASSWORD="${POSTGRES_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}"

echo "[$(date -u +'%Y-%m-%dT%H:%M:%SZ')] [SUCCESS] Restauração finalizada com sucesso."

