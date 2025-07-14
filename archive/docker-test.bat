@echo off
REM docker-test.bat - Windows version of the Docker test script

echo Setting up Docker test environment for RealEntity nodes...

if "%1"=="stop" goto stop
if "%1"=="logs" goto logs
if "%1"=="status" goto status
if "%1"=="clean" goto clean
if "%1"=="start" goto start
if "%1"=="" goto start

echo Usage: %0 {start^|stop^|logs^|status^|clean}
echo   start  - Start the entire test network
echo   stop   - Stop all nodes
echo   logs   - Show real-time logs from all nodes
echo   status - Show current network status
echo   clean  - Clean up containers and images
goto end

:start
echo [INFO] Checking Docker installation...
docker --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker is not installed or not running
    goto end
)

echo [INFO] Cleaning up existing containers...
docker-compose down --remove-orphans 2>nul

echo [INFO] Building RealEntity Docker image...
docker-compose build

echo [INFO] Starting bootstrap node...
docker-compose up -d bootstrap

echo [INFO] Waiting for bootstrap node to be ready...
timeout /t 5 /nobreak >nul

echo [INFO] Getting bootstrap peer ID...
for /f "tokens=*" %%i in ('docker logs realentity-bootstrap 2^>^&1 ^| findstr "Node started with ID:"') do set BOOTSTRAP_LOG=%%i
for /f "tokens=5" %%i in ("%BOOTSTRAP_LOG%") do set BOOTSTRAP_ID=%%i

if "%BOOTSTRAP_ID%"=="" (
    echo [ERROR] Failed to get bootstrap peer ID
    docker logs realentity-bootstrap
    goto end
)

echo [INFO] Bootstrap node started with ID: %BOOTSTRAP_ID%

echo [INFO] Updating peer configurations...
powershell -Command "(Get-Content configs\peer-config.json) -replace 'BOOTSTRAP_ID_PLACEHOLDER', '%BOOTSTRAP_ID%' | Set-Content configs\peer-config.json"
powershell -Command "(Get-Content docker-compose.yml) -replace 'BOOTSTRAP_ID_PLACEHOLDER', '%BOOTSTRAP_ID%' | Set-Content docker-compose.yml"

echo [INFO] Starting peer nodes...
docker-compose up -d peer1 peer2 peer3 peer4

echo [INFO] Waiting for peer nodes to connect...
timeout /t 10 /nobreak >nul

echo [INFO] Network Status:
docker-compose ps

echo [INFO] You can monitor logs with: docker-compose logs -f
echo [INFO] To test services: docker-test-services.bat
echo [INFO] To stop all nodes: docker-test.bat stop
goto end

:stop
echo [INFO] Stopping all nodes...
docker-compose down
goto end

:logs
docker-compose logs -f
goto end

:status
echo [INFO] Network Status:
docker-compose ps
echo.
echo [INFO] Recent bootstrap logs:
docker logs realentity-bootstrap --tail 10
goto end

:clean
echo [INFO] Cleaning up containers and images...
docker-compose down --remove-orphans
docker system prune -f
goto end

:end
