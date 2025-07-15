@echo off
echo Starting RealEntity Peer Node
echo =============================

if "%1"=="" (
    echo Usage: start-peer.bat BOOTSTRAP_ADDRESS
    echo.
    echo Example: start-peer.bat "/ip4/127.0.0.1/tcp/4001/p2p/12D3KooW..."
    echo.
    echo First, start the bootstrap node and copy its peer address from the console output.
    pause
    exit /b 1
)

echo Connecting to bootstrap node: %1
echo.

REM Update the config file with the bootstrap peer address
powershell -Command "(Get-Content config-peer.json) -replace 'BOOTSTRAP_PEER_ADDRESS_HERE', '%1' | Set-Content config-peer-temp.json"

echo Starting peer node on port 4002...
go run ../../cmd/main.go -config config-peer-temp.json

REM Clean up temp file
del config-peer-temp.json 2>nul

pause
