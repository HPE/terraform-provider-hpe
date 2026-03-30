#!/usr/bin/env bash

set -e

os="linux"
arch="amd64"
repo="HPE/terraform-provider-hpe"
install_dir="${INSTALL_DIR:-/usr/local/bin}"
binary_name="tfmigrator"

get_latest_release () {
  local release_url="https://api.github.com/repos/${repo}/releases/latest"
  curl -sL "$release_url" | grep -Po '"tag_name": "\K.*?(?=")'
}

download_and_install () {
  local tmp_dir
  tmp_dir=$(mktemp -d)
  
  local archive_name="migration_tool_${version_number}_${os}_${arch}.zip"
  local download_url="https://github.com/${repo}/releases/download/${VERSION}/${archive_name}"
  
  cd "$tmp_dir"
  curl -sL "$download_url" -o "$archive_name" && \
    unzip -q "$archive_name"
  
  # Find the binary (might be named 'tfmigrator' or 'migration_tool')
  local binary_path=""
  if [[ -f "tfmigrator" ]]; then
    binary_path="tfmigrator"
  elif [[ -f "migration_tool" ]]; then
    binary_path="migration_tool"
  else
    echo "Error: Could not find binary in archive"
    rm -rf "$tmp_dir"
    exit 1
  fi
  
  mkdir -p "$install_dir"
  mv "$binary_path" "${install_dir}/${binary_name}"
  chmod +x "${install_dir}/${binary_name}"
  rm -rf "$tmp_dir"
}

VERSION=${VERSION:=$(get_latest_release)}
version_number=${VERSION//v}
download_and_install
