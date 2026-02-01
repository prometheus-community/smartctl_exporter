param(
    [string]$Version = $(Get-Content -Raw -Path "$PSScriptRoot/../../VERSION").Trim()
)

$ErrorActionPreference = "Stop"

$root = Resolve-Path "$PSScriptRoot/../.."
$buildDir = Join-Path $PSScriptRoot "build"
$distDir = Join-Path $root "dist"
$vendorDir = Join-Path $PSScriptRoot "vendor"

New-Item -ItemType Directory -Force -Path $buildDir, $distDir, $vendorDir | Out-Null

$exePath = Join-Path $distDir "smartctl_exporter.exe"
$configPath = Join-Path $distDir "smartctl_exporter.yml"

Write-Host "Building smartctl_exporter.exe..."
Push-Location $root
go build -o $exePath
Pop-Location

Copy-Item -Force -Path (Join-Path $PSScriptRoot "smartctl_exporter.yml") -Destination $configPath

$msiPath = Join-Path $buildDir "smartctl_exporter.msi"

Write-Host "Building MSI..."
wix build "$PSScriptRoot/smartctl_exporter.wxs" `
    -dProductVersion=$Version `
    -dExeSource=$exePath `
    -dConfigSource=$configPath `
    -o $msiPath

Write-Host "Done. Artifacts:"
Write-Host " - $msiPath"
