.\build-backend.ps1
$built=$?
if ($built) {
    .\dist\git-stories-server -wd .\dist
}
