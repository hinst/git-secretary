cd git-stories-server
$env:GOOS=""
go build -o ../dist/
$built=$?
cd ..
if (-Not $built) {
    exit $built
}

cd git-stories-server
$env:GOOS="linux"
go build -o ../dist/
$built=$?
cd ..
if (-Not $built) {
    exit $built
}

cd story-girls-standard
$env:GOOS=""
go-bindata resources
$built=$?
if ($built) {
    go build -o ../dist/plugins/
    $built=$?
}
cd ..
if (-Not $built) {
    exit $built
}

cd story-girls-standard
$env:GOOS="linux"
go-bindata resources
$built=$?
if ($built) {
    go build -o ../dist/plugins/
    $built=$?
}
cd ..
if (-Not $built) {
    exit $built
}

Copy-Item .\configuration.json .\dist\

Write-Output "Backend build succeeded"