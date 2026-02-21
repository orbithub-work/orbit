#!/bin/bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

echo "[Backend] Starting Go service on port 32000..."
exec go run cmd/core/main.go
