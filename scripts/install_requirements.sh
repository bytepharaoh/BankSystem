#!/usr/bin/env bash
set -euo pipefail

# Installs local development requirements for this project.
# Supported package managers: Homebrew (macOS/Linux), apt-get (Linux, partial).

need_cmd() {
  command -v "$1" >/dev/null 2>&1
}

install_with_brew() {
  local pkg="$1"
  if ! brew list "$pkg" >/dev/null 2>&1; then
    echo "[install] brew install $pkg"
    brew install "$pkg"
  else
    echo "[skip] $pkg is already installed"
  fi
}

install_with_apt() {
  local pkg="$1"
  echo "[install] apt-get install -y $pkg"
  sudo apt-get install -y "$pkg"
}

echo "==> Checking base tools"

if ! need_cmd go; then
  echo "[error] Go is not installed. Install Go 1.25+ first: https://go.dev/dl/"
  exit 1
fi

if ! need_cmd docker; then
  if need_cmd brew; then
    install_with_brew docker
  elif need_cmd apt-get; then
    sudo apt-get update
    install_with_apt docker.io
  else
    echo "[error] docker is missing and no supported package manager was found"
    exit 1
  fi
else
  echo "[skip] docker is already installed"
fi

if ! need_cmd migrate; then
  if need_cmd brew; then
    install_with_brew golang-migrate
  elif need_cmd apt-get; then
    sudo apt-get update
    install_with_apt golang-migrate
  else
    echo "[warn] migrate is missing. Please install manually: https://github.com/golang-migrate/migrate"
  fi
else
  echo "[skip] migrate is already installed"
fi

echo "==> Installing Go-based CLI tools"
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install go.uber.org/mock/mockgen@latest

echo "==> Ensuring Go bin path is available"
if [[ ":$PATH:" != *":$(go env GOPATH)/bin:"* ]]; then
  echo "[warn] Add this to your shell profile:"
  echo "export PATH=\"$(go env GOPATH)/bin:\$PATH\""
fi

echo "Done. Requirements installation completed."
