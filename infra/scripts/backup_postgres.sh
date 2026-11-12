#!/usr/bin/env bash
# ==============================================================================
# Geevheit Intelligence 360 — Script de Backup Automatizado do PostgreSQL
# Retenção: 30 dias (configurável via BACKUP_RETENTION_DAYS)
# Formato: dump custom binário compactado (.dump.gz)
# ==============================================================================

set -euo pipefail

BACKUP_DIR="${BACKUP_DIR:-/backups}"
BACKUP_RETENTION_DAYS="${BACKUP_RETENTION_DAYS:-30}"
DB_HOST="${DB_HOST:-postgres}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${POSTGRES_DB:-geevheit_intelligence}"
DB_USER="${POSTGRES_USER:-geevheit}"
TIMESTAMP="$(date +'%Y%m%d_%H%M%S')"
BACKUP_FILE="${BACKUP_DIR}/${DB_NAME}_${TIMESTAMP}.sql.gz"

mkdir -p "${BACKUP_DIR}"

echo "[$(date -u +'%Y-%m-%dT%H:%M:%SZ')] [INFO] Iniciando backup da base de dados '${DB_NAME}'..."

# Executa pg_dump e compacta com gzip em streaming
PGPASSWORD="${POSTGRES_PASSWORD}" pg_dump -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" --clean --if-exists --no-owner --no-privileges | gzip -9 > "${BACKUP_FILE}"

FILESIZE=$(du -h "${BACKUP_FILE}" | cut -f1)
echo "[$(date -u +'%Y-%m-%dT%H:%M:%SZ')] [SUCCESS] Backup concluído com sucesso: ${BACKUP_FILE} (${FILESIZE})"

# Remove backups mais antigos que a janela de retenção (padrão 30 dias)
echo "[$(date -u +'%Y-%m-%dT%H:%M:%SZ')] [INFO] Aplicando política de retenção (${BACKUP_RETENTION_DAYS} dias)..."
find "${BACKUP_DIR}" -name "${DB_NAME}_*.sql.gz" -type f -mtime +"${BACKUP_RETENTION_DAYS}" -exec rm -f {} +
echo "[$(date -u +'%Y-%m-%dT%H:%M:%SZ')] [INFO] Limpeza de arquivos legados concluída."

