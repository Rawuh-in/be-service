#!/usr/bin/env bash
# generate.sh
# Regenerate Swagger docs using swaggo/swag. Run this from the repository root.
# It runs swag in module-mode and generates docs into ./docs

set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Allow overriding the swag command or flags via environment vars
SWAG_CMD="${SWAG_CMD:-go run github.com/swaggo/swag/cmd/swag}"

# By default run swag from the cmd/server directory so it doesn't try to
# `go list ./` on the repository root (which has no .go files) and emit the
# harmless "no Go files in ./" warning. You can still override SWAG_FLAGS
# via environment if you want a custom invocation.
# Scan only cmd/server and the handler packages under internal to avoid
# scanning the whole internal tree twice (which can produce duplicate-route
# warnings). Override SWAG_FLAGS if you need different behavior.
SWAG_FLAGS="${SWAG_FLAGS:-init -g main.go -o ../../docs \
	--parseInternal --parseDependency --parseDependencyLevel 3 --parseFuncBody \
	--dir .,../../internal/event/handler,../../internal/guest/handler,../../internal/project/handler,../../internal/user/handler,../../internal/auth/handler}"

echo "Generating swagger docs..."

# If the user provided SWAG_FLAGS explicitly, respect it but still run inside
# cmd/server so relative -g/main paths work and root scanning is avoided.
(
	cd cmd/server
	GOFLAGS=-mod=mod $SWAG_CMD $SWAG_FLAGS
)

echo "Swagger docs generated at ./docs/swagger.json and ./docs/docs.go"
