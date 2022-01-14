cd frontend
npm run build
$built=$?
cd ..
if (-Not $built) {
    exit $built
}

Write-Output "Frontend build succeeded"