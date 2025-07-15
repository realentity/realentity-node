@echo off
echo Starting 3 RealEntity nodes for peer discovery testing
echo =====================================================

echo Starting Node 1 (port 4001, HTTP 8081)...
start "Node 1" cmd /c "go run ../../cmd/main.go -config config-node1.json"
timeout /t 3

echo Starting Node 2 (port 4002, HTTP 8082)...
start "Node 2" cmd /c "go run ../../cmd/main.go -config config-node2.json"
timeout /t 3

echo Starting Node 3 (port 4003, HTTP 8083)...
start "Node 3" cmd /c "go run ../../cmd/main.go -config config-node3.json"

echo.
echo All nodes started! 
echo Check the console windows for peer discovery messages.
echo.
echo API endpoints:
echo - Node 1: http://localhost:8081/health
echo - Node 2: http://localhost:8082/health  
echo - Node 3: http://localhost:8083/health
echo.
echo Press any key to exit and stop all nodes...
pause >nul

echo Stopping all nodes...
taskkill /f /fi "WindowTitle eq Node 1*" >nul 2>nul
taskkill /f /fi "WindowTitle eq Node 2*" >nul 2>nul
taskkill /f /fi "WindowTitle eq Node 3*" >nul 2>nul
echo All nodes stopped.
