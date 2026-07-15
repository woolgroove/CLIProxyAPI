#!/bin/sh
set -e

cd /CLIProxyAPI

# Image ships config.example.yaml only; create a runtime config if missing.
if [ ! -f config.yaml ]; then
  cp config.example.yaml config.yaml
fi

# Render (and similar platforms) inject PORT; CLIProxyAPI only reads config.yaml.
if [ -n "${PORT:-}" ]; then
  if grep -qE '^port:[[:space:]]*' config.yaml; then
    sed -i "s/^port:[[:space:]]*.*/port: ${PORT}/" config.yaml
  else
    printf '\nport: %s\n' "${PORT}" >> config.yaml
  fi
fi

exec ./CLIProxyAPI "$@"
