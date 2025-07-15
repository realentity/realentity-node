@echo off
echo Starting RealEntity Bootstrap Node
echo ===================================

echo Starting bootstrap node on port 4001...
echo This will be the first node that others connect to.
echo.

go run ../../cmd/main.go -config config-bootstrap.json

pause
