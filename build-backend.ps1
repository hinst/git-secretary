cd git-stories-server
go build -o ../dist/
$built=$?
cd ..
if (-Not $built) {
    exit $built
}

cd story-girls-standard
go-bindata resources
$built=$?
if ($built) {
    go build -o ../dist/plugins/story-girls-standard.exe
    $built=$?
}
cd ..
if (-Not $built) {
    exit $built
}

Copy-Item .\configuration.json .\dist\

Write-Output "Backend build succeeded"