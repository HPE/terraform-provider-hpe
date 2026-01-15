#!/usr/bin/env bash

set -e

os="linux"
arch="amd64"
repo="HPE/terraform-provider-hpe"
linux_hpe_dir="${HOME}/.local/share/terraform/plugins/registry.terraform.io/hpe/hpe"

get_latest_release () {
  local release_url="https://api.github.com/repos/${repo}/releases/latest"
  curl -sL "$release_url" | grep -Po '"tag_name": "\K.*?(?=")'
}

download_and_extract () {
  local dest_dir="${linux_hpe_dir}/${version_number}/${os}_${arch}/"
  local hpe_zip="terraform-provider-hpe_${version_number}_${os}_${arch}.zip"
  local hpe_dl_url="https://github.com/${repo}/releases/download/${VERSION}/${hpe_zip}"

  mkdir -p "$dest_dir" && cd "$dest_dir"
  curl -sL "$hpe_dl_url" -o "$hpe_zip" && \
    unzip -u "$hpe_zip" && \
    rm -f "$hpe_zip"
}

VERSION=${VERSION:=$(get_latest_release)}
version_number=${VERSION//v}
download_and_extract
