#!/usr/bin/env zsh
set -euo pipefail

# Load env if present
if [ -f ".env" ]; then
  set -a
  source .env
  set +a
fi

# Run backend server
go run ./cmd/server

