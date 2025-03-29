#!/bin/bash

source .env

set -e

# === CONFIGURATION ===
IMAGE_NAME="kurtis/automations"
TAG=$(date +%Y%m%d%H%M%S)
ARCH="linux/arm64"

# 🔐 Remote SSH credentials
REMOTE_HOST="${REMOTE_HOST:-}"
REMOTE_USER="${REMOTE_USER:-}"
REMOTE_PASS="${REMOTE_PASS:-}"
REMOTE_DIR="$HOME/wrkspc"
NEW_IMAGE="$IMAGE_NAME:$TAG"

# === VALIDATION ===
if [[ -z "$REMOTE_HOST" || -z "$REMOTE_USER" || -z "$REMOTE_PASS" ]]; then
  echo "❌ ERROR: REMOTE_HOST, REMOTE_USER, and REMOTE_PASS must be set."
  exit 1
fi

# === BUILD & PUSH IMAGE ===
echo "🛠️ Building Docker image for $ARCH..."
docker buildx build --platform "$ARCH" -t "$NEW_IMAGE" --push .

# === REMOTE COMMAND (YQ via Docker) ===
REMOTE_COMMAND=$(cat <<EOF
set -e
cd "$REMOTE_DIR"
pwd
ls -lha
echo "📦 Updating docker-compose.yaml using Dockerized yq..."
docker run --rm -v "$REMOTE_DIR":/workdir \
  mikefarah/yq \
  -i '.services.automations.image = "$NEW_IMAGE"' /workdir/docker-compose.yaml

echo "🚀 Restarting Docker Compose stack..."
docker compose up -d
EOF
)

# Inject NEW_IMAGE value into the remote command safely
REMOTE_COMMAND="${REMOTE_COMMAND//\$NEW_IMAGE/$NEW_IMAGE}"

# === EXECUTE REMOTE COMMAND ===
echo "🔐 Connecting to $REMOTE_HOST via sshpass..."
sshpass -p "$REMOTE_PASS" ssh -o StrictHostKeyChecking=no "$REMOTE_USER@$REMOTE_HOST" "$REMOTE_COMMAND"

echo "✅ Done. Deployed image: $NEW_IMAGE to $REMOTE_HOST"
