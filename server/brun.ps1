# Quick build and run in developer mode
$succeeded = &".\build.ps1"
if ($succeeded) {
    .\git-stories-server -ao=false
}