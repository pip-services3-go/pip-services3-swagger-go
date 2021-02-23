#!/usr/bin/env pwsh

# Get component data and set necessary variables
$component = Get-Content -Path "component.json" | ConvertFrom-Json
$buildImage="$($component.registry)/$($component.name):$($component.version)-$($component.build)-build"
$testImage="$($component.registry)/$($component.name):$($component.version)-$($component.build)-test"
$rcImage="$($component.registry)/$($component.name):$($component.version)-$($component.build)-rc"

# Clean up build directories
Get-ChildItem -Path "." -Include "obj" -Recurse | foreach($_) { Remove-Item -Force -Recurse $_.FullName }

# Remove docker images
docker rmi $buildImage --force
docker rmi $testImage --force
docker rmi $rcImage --force
docker image prune --force

# Remove existed containers
docker ps -a | Select-String -Pattern "Exit" | foreach($_) { docker rm $_.ToString().Split(" ")[0] }