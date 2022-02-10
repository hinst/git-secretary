go build
$succeeded = $?
Write-Host "Build git-stories-server $succeeded"
return $succeeded