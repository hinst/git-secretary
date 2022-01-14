cd frontend
npm run build
$built=$?
cd ..
if (-Not $built) {
    exit $built
}

Remove-Item -Recurse .\dist\frontend\
Copy-Item -Recurse -Path .\frontend\build -Destination .\dist\frontend\

Write-Output "Frontend build succeeded"