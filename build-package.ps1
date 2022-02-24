if (Test-Path .\dist) {
    Remove-Item -Recurse .\dist
}
mkdir .\dist

$Null = &".\build.ps1"

Copy-Item -Recurse -Path .\frontend\build -Destination .\dist\frontend\
Copy-Item configuration.json .\dist\configuration.json

if (Test-Path .\git-secretary-windows) {
    Remove-Item -Recurse .\git-secretary-windows
}
Copy-Item -Recurse -Path .\dist .\git-secretary-windows
Copy-Item .\server\git-secretary.exe .\git-secretary-windows\
Compress-Archive -Path .\git-secretary-windows -DestinationPath .\git-secretary-windows.zip -Force

if (Test-Path .\dist-linux) {
    Remove-Item -Recurse .\dist-linux
}
Copy-Item -Recurse -Path .\dist .\git-secretary-linux
Copy-Item .\server\git-secretary .\git-secretary-linux\
Compress-Archive -Path .\git-secretary-linux -DestinationPath .\git-secretary-linux.zip -Force