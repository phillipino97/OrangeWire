setlocal enabledelayedexpansion
set port=2000
for /l %%x in (1, 1, 20) do (
    set /A port=!port!+1
    echo !port!
    start cmd.exe /k go run .\SecureFileSharing.go -serverport=!port! -port=2000 -filepath=!port!
)