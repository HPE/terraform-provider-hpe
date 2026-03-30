param ($VERSION, $INSTALL_DIR)

$os="windows"
$arch="amd64"
$repo="HPE/terraform-provider-hpe"
$binary_name="tfmigrator.exe"

if (!$INSTALL_DIR) {
    $INSTALL_DIR="$env:LOCALAPPDATA\tfmigrator"
}

$users_pwd = Get-Location

function get_latest_release {
    Write-Host "Getting latest release"
    $release_url="https://api.github.com/repos/${repo}/releases/latest"
    $tag = (Invoke-WebRequest $release_url | ConvertFrom-Json)[0].tag_name
    $VERSION=${tag}
    
    $VERSION
}

if (!$VERSION) {
    $VERSION=get_latest_release
}

$version_number=$VERSION -replace 'v'  

$tmp_dir="${env:TEMP}\tfmigrator-install-$(Get-Random)"
$archive_name="migration_tool_${version_number}_${os}_${arch}.zip"
$archive_path="${tmp_dir}\${archive_name}"
$download_url="https://github.com/${repo}/releases/download/${VERSION}/${archive_name}"

New-Item -ItemType Directory -Path "$tmp_dir" -Force | Out-Null
Set-Location "$tmp_dir"

try {
    Invoke-WebRequest $download_url -Out $archive_path     
}
catch {
    Write-Host "Error: The version that was specified does not exist."

    Set-Location "${users_pwd}"
    Remove-Item -Path "${tmp_dir}" -Recurse -Force -ErrorAction SilentlyContinue 

    Write-Host "Exiting..."
    Return 
}

Write-Host "Extracting release files"
Expand-Archive $archive_path -Force

# Find the binary (might be named 'tfmigrator.exe' or 'migration_tool.exe')
$extract_dir = $archive_name -replace '.zip'
$binary_path = $null
if (Test-Path "${extract_dir}\tfmigrator.exe") {
    $binary_path = "${extract_dir}\tfmigrator.exe"
}
elseif (Test-Path "${extract_dir}\migration_tool.exe") {
    $binary_path = "${extract_dir}\migration_tool.exe"
}
elseif (Test-Path "tfmigrator.exe") {
    $binary_path = "tfmigrator.exe"
}
elseif (Test-Path "migration_tool.exe") {
    $binary_path = "migration_tool.exe"
}
else {
    Write-Host "Error: Could not find binary in archive"
    Set-Location "${users_pwd}"
    Remove-Item -Path "${tmp_dir}" -Recurse -Force -ErrorAction SilentlyContinue 
    Return
}

New-Item -ItemType Directory -Path "$INSTALL_DIR" -Force | Out-Null

$install_path = "${INSTALL_DIR}\${binary_name}"
Get-ChildItem -Path $binary_path -File | Move-Item -Destination $install_path -Force

Remove-Item -Path "${tmp_dir}" -Recurse -Force -ErrorAction SilentlyContinue 
Write-Host "Complete"

Set-Location "${users_pwd}"
