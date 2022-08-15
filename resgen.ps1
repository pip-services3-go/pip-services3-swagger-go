#!/usr/bin/env pwsh

Set-StrictMode -Version latest
$ErrorActionPreference = "Stop"

# Generate image and container names using the data in the "component.json" file
$component = Get-Content -Path "component.json" | ConvertFrom-Json

$resImage = "$($component.registry)/$($component.name):$($component.version)-$($component.build)-res"
$container = $component.name

# Remove build files
if (Test-Path "$PSScriptRoot/resources") {
    Remove-Item -Recurse -Force -Path "$PSScriptRoot/resources/*"
} else {
    $null = New-Item -ItemType Directory -Force -Path "$PSScriptRoot/resources"
}
if (Test-Path "$PSScriptRoot/example/resources") {
    Remove-Item -Recurse -Force -Path "$PSScriptRoot/example/resources/*"
} else {
    $null = New-Item -ItemType Directory -Force -Path "$PSScriptRoot/example/resources"
}

# Build docker image
docker build -f "$PSScriptRoot/docker/Dockerfile.resgen" -t $resImage "$PSScriptRoot/."

# Run resgen container
docker run -d --name $container $resImage

# Copy resources from container
docker cp "$($container):/app/src/example/resources" "$PSScriptRoot/example"
docker cp "$($container):/app/src/resources" "$PSScriptRoot/."
# Remove docgen container
docker rm $container --force

# Verify resources
if (-not (Test-Path "$PSScriptRoot/resources/*.go")) {
    Write-Error "resources folder doesn't exist in root dir. Watch logs above."
}
