cd frontend
$env:PUBLIC_URL="/git-stories/static-files"
$env:REACT_APP_API_URL="/git-stories/api"
npm run build
$built=$?
cd ..
if (-Not $built) {
    exit $built
}

Remove-Item -Recurse .\dist\frontend\
Copy-Item -Recurse -Path .\frontend\build -Destination .\dist\frontend\

Write-Output "Frontend build succeeded"