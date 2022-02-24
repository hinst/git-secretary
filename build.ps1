Set-Location frontend
$Null = &".\build.ps1"
$built=$?
Write-Host "Build frontend $built"
Set-Location ..
if (-Not $built) {
    exit $built
}

Set-Location server
$Null = &".\build-all.ps1"
$built=$?
Write-Host "Build server $built"
Set-Location ..
if (-Not $built) {
    exit $built
}

Write-Host "Build all: succeeded"

return
Compress-Archive -Path .\dist\* -DestinationPath .\git-stories.zip -Force