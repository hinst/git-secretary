.\build-frontend.ps1
$built=$?
if (-Not $built) {
    exit $built
}

.\build-backend.ps1
$built=$?
if (-Not $built) {
    exit $built
}

Compress-Archive -Path .\dist\* -DestinationPath .\git-stories.zip

Write-Output "Build all succeeded"