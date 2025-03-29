#!/bin/bash

source .env

set -e

# === REMOTE SSH CONFIG ===
REMOTE_HOST="${REMOTE_HOST:-}"
REMOTE_USER="${REMOTE_USER:-}"
REMOTE_PASS="${REMOTE_PASS:-}"
CONTAINER_NAME="automations" # or full container name if different
LOCAL_LOG_DIR="./logs"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
LOCAL_LOG_FILE="$LOCAL_LOG_DIR/${CONTAINER_NAME}_$TIMESTAMP.log"

# === VALIDATION ===
if [[ -z "$REMOTE_HOST" || -z "$REMOTE_USER" || -z "$REMOTE_PASS" ]]; then
  echo "❌ ERROR: REMOTE_HOST, REMOTE_USER, and REMOTE_PASS must be set."
  exit 1
fi

# === CREATE LOCAL DIR IF NEEDED ===
mkdir -p "$LOCAL_LOG_DIR"

# === FETCH LOGS ===
echo "📡 Fetching logs from $CONTAINER_NAME on $REMOTE_HOST..."
sshpass -p "$REMOTE_PASS" ssh -o StrictHostKeyChecking=no "$REMOTE_USER@$REMOTE_HOST" \
  "docker logs $CONTAINER_NAME 2>&1 | sed -r 's/\x1B\[[0-9;]*[mK]//g'" > "$LOCAL_LOG_FILE"

echo "✅ Logs saved to $LOCAL_LOG_FILE"
