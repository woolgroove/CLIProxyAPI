#!/bin/sh
set -e

cd /CLIProxyAPI

# Image ships config.example.yaml only; create a runtime config if missing.
if [ ! -f config.yaml ]; then
  cp config.example.yaml config.yaml
fi

# Best-effort local file port patch. The binary also honors $PORT after config load
# (needed when PGSTORE/GITSTORE/OBJECTSTORE supplies config.yaml).
if [ -n "${PORT:-}" ] && [ -f config.yaml ]; then
  if grep -qE '^port:[[:space:]]*' config.yaml; then
    sed -i "s/^port:[[:space:]]*.*/port: ${PORT}/" config.yaml
  else
    printf '\nport: %s\n' "${PORT}" >> config.yaml
  fi
fi

exec ./CLIProxyAPI "$@"
