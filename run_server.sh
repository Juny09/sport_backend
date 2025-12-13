#!/bin/bash
# Load env vars from .env file, ignoring comments
set -a
[ -f .env ] && . .env
set +a

# Run the server
./server
