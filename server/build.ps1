$Null = @(go build)
$succeeded = $?
Write-Host "Build git-secretary $succeeded"
return $succeeded