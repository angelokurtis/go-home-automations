#!/bin/bash

source .env

set -e

# === REMOTE SSH CONFIG ===
REMOTE_HOST="${REMOTE_HOST:-}"
REMOTE_USER="${REMOTE_USER:-}"
REMOTE_PASS="${REMOTE_PASS:-}"
CONTAINER_NAME="automations"  # Adjust if needed

# === VALIDATION ===
if [[ -z "$REMOTE_HOST" || -z "$REMOTE_USER" || -z "$REMOTE_PASS" ]]; then
  echo "❌ ERROR: REMOTE_HOST, REMOTE_USER, and REMOTE_PASS must be set."
  exit 1
fi

# === FOLLOW LOGS ===
echo "🔐 Connecting to $REMOTE_HOST and following logs for $CONTAINER_NAME..."
sshpass -p "$REMOTE_PASS" ssh -o StrictHostKeyChecking=no "$REMOTE_USER@$REMOTE_HOST" \
  "docker logs -f $CONTAINER_NAME" # Follow logs in real-time

echo "✅ Now following logs for $CONTAINER_NAME on $REMOTE_HOST."
