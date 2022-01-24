.\build-backend.ps1
$built = $LASTEXITCODE
if (-Not $built) {
    exit $built
}

.\dist\git-stories-server --wd C:\Dev\git-stories\dist