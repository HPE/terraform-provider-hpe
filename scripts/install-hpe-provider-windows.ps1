param ($VERSION)

$os="windows"
$arch="amd64"
$repo="HPE/terraform-provider-hpe"
$windows_hpe_dir="$env:appdata\terraform.d\plugins\registry.terraform.io\hpe\hpe"

$users_pwd = Get-Location

function get_latest_release {
    Write-Host Getting latest release
    $release_url="https://api.github.com/repos/${repo}/releases/latest"
    $tag = (Invoke-WebRequest $release_url | ConvertFrom-Json)[0].tag_name
    $VERSION=${tag}
    
    $VERSION
}

if (!$VERSION) {
    $VERSION=get_latest_release
}

$version_number=$VERSION -replace 'v'  

$dest_dir="${windows_hpe_dir}\${version_number}\${os}_${arch}\"
$hpe_zip="terraform-provider-hpe_${version_number}_${os}_${arch}.zip"
$hpe=$hpe_zip -replace '.zip'
$hpe_dl_url="https://github.com/${repo}/releases/download/${VERSION}/${hpe_zip}"

New-Item -ItemType Directory -Path "$dest_dir" -Force | Out-Null
Set-Location "$dest_dir"

try {
    Invoke-WebRequest $hpe_dl_url -Out $hpe_zip     
}
catch {
    Write-Host "Error: The version that was specified does not exist."

    Set-Location "${users_pwd}"
    Remove-Item -Path "${windows_hpe_dir}\${version_number}" -Recurse -Force -ErrorAction SilentlyContinue 

    Write-Host "Exiting..."
    Return 
}

Write-Host Extracting release files
Expand-Archive $hpe_zip -Force

Get-ChildItem -Path $hpe -Recurse -File | Move-Item -Destination $dest_dir -Force

Remove-Item $hpe_zip -Recurse -Force -ErrorAction SilentlyContinue 
Remove-Item $hpe -Recurse -Force -ErrorAction SilentlyContinue 
Write-Host Complete

Set-Location "${users_pwd}"
