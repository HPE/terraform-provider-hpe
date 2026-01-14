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
macos_hpe_dir="${HOME}/.terraform.d/plugins/registry.terraform.io/hpe/hpe"

get_latest_release () {
  local release_url="https://api.github.com/repos/${repo}/releases/latest"
  curl -sL "$release_url" | perl -nle 'print $1 if /"tag_name":\s*"(.*?)"/'
}

download_and_extract () {
  local dest_dir="${macos_hpe_dir}/${version_number}/${os}_${arch}/"
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
