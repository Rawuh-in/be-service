#!/usr/bin/env bash
# generate.sh
# Regenerate Swagger docs using swaggo/swag. Run this from the repository root.
# It runs swag in module-mode and generates docs into ./docs

set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Allow overriding the swag command or flags via environment vars
SWAG_CMD="${SWAG_CMD:-go run github.com/swaggo/swag/cmd/swag}"
SWAG_FLAGS="${SWAG_FLAGS:-init -g ./cmd/server/main.go -o ./docs}"

echo "Generating swagger docs..."
GOFLAGS=-mod=mod $SWAG_CMD $SWAG_FLAGS

echo "Swagger docs generated at ./docs/swagger.json and ./docs/docs.go"
