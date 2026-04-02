#!/usr/bin/env bash

set -e

os="darwin"

if [[ `uname -m` == 'x86_64' ]] 
then
  echo 'Intel Architecture detected'
  arch="amd64"
elif [[ `uname -m` == 'arm64' ]] 
then
  echo 'Apple Silicon Architecture detected'
  arch="arm64"
fi

repo="HPE/terraform-provider-hpe"
install_dir="${INSTALL_DIR:-/usr/local/bin}"
binary_name="tfmigrator"

get_latest_release () {
  local release_url="https://api.github.com/repos/${repo}/releases/latest"
  curl -sL "$release_url" | perl -nle 'print $1 if /"tag_name":\s*"(.*?)"/'
}

download_and_install () {
  local tmp_dir
  tmp_dir=$(mktemp -d)
  
  local archive_name="migration_tool_${version_number}_${os}_${arch}.zip"
  local download_url="https://github.com/${repo}/releases/download/${VERSION}/${archive_name}"
  
  cd "$tmp_dir"
  curl -sL "$download_url" -o "$archive_name" && \
    unzip -q "$archive_name"
  
  # Binary is named migration_tool_v${version_number}
  local binary_path="migration_tool_v${version_number}"
  if [[ ! -f "$binary_path" ]]; then
    echo "Error: Could not find binary ${binary_path} in archive"
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
