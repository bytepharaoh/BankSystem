#!/usr/bin/env bash
set -euo pipefail

# Installs local development requirements for this project.
# Supported package managers:
# - Homebrew (macOS/Linux)
# - apt-get (Ubuntu/Debian)

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

get_linux_codename() {
  if need_cmd lsb_release; then
    lsb_release -sc
    return
  fi

  if [[ -f /etc/os-release ]]; then
    # shellcheck disable=SC1091
    . /etc/os-release
    if [[ -n "${VERSION_CODENAME:-}" ]]; then
      echo "$VERSION_CODENAME"
      return
    fi
  fi

  echo "stable"
}

install_migrate_with_apt_repo() {
  local codename
  codename="$(get_linux_codename)"

  echo "==> Configuring migrate apt repository (packagecloud)"
  sudo mkdir -p /etc/apt/keyrings
  curl -fsSL https://packagecloud.io/golang-migrate/migrate/gpgkey \
    | sudo gpg --dearmor -o /etc/apt/keyrings/migrate.gpg
  echo "deb [signed-by=/etc/apt/keyrings/migrate.gpg] https://packagecloud.io/golang-migrate/migrate/ubuntu/ ${codename} main" \
    | sudo tee /etc/apt/sources.list.d/migrate.list >/dev/null
  sudo apt-get update
  install_with_apt migrate
}

install_migrate() {
  if need_cmd migrate; then
    echo "[skip] migrate is already installed"
    return
  fi

  if need_cmd brew; then
    install_with_brew golang-migrate
    return
  fi

  if need_cmd apt-get; then
    # First, try direct package install (in case repo already exists).
    set +e
    sudo apt-get update
    sudo apt-get install -y migrate
    local apt_status=$?
    set -e
    if [[ $apt_status -ne 0 ]]; then
      install_migrate_with_apt_repo
    fi
    return
  fi

  echo "[warn] Package manager install for migrate is unavailable. Falling back to go install."
  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
}

echo "==> Checking base tools"

if ! need_cmd go; then
  if need_cmd brew; then
    install_with_brew go
  elif need_cmd apt-get; then
    sudo apt-get update
    install_with_apt golang-go
  else
    echo "[error] Go is missing and no supported package manager was found."
    echo "[hint] Install Go 1.25+ manually from https://go.dev/dl/"
    exit 1
  fi
else
  echo "[skip] go is already installed"
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

install_migrate

echo "==> Installing Go-based CLI tools"
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install go.uber.org/mock/mockgen@latest

echo "==> Ensuring Go bin path is available"
if [[ ":$PATH:" != *":$(go env GOPATH)/bin:"* ]]; then
  echo "[warn] Add this to your shell profile:"
  echo "export PATH=\"$(go env GOPATH)/bin:\$PATH\""
fi

echo "Done. Requirements installation completed."
