#!/bin/bash
set -e

# Load environment variables from .env if present
if [ -f ".env" ]; then
  export $(grep -v '^#' .env | xargs -d '\n')
fi

for dir in */; do
  if [ -f "$dir/main.go" ]; then
    echo "Running $dir"
    (cd "$dir" && go run .)
  fi
done
