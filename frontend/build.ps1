$env:PUBLIC_URL="/git-stories/static-files"
$env:REACT_APP_API_URL="/git-stories/api"
npm run build | Out-Host
return $?
