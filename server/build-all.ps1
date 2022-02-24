$ok = &".\build.ps1"
if (-Not $ok) {
    return $false
}
$ok = &".\build-linux.ps1"
if (-Not $ok) {
    return $false
}
Write-Host "Build all server: done"
return $ok