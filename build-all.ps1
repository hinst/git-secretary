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

Write-Output "Build all succeeded"