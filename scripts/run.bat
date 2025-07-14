@echo off
REM scripts\run.bat - Script runner utility for RealEntity Node (Windows)

setlocal enabledelayedexpansion

set "SCRIPT_DIR=%~dp0"
set "PROJECT_ROOT=%SCRIPT_DIR%\.."

if "%~1"=="" (
    call :print_usage
    exit /b 1
)

set "SCRIPT_NAME=%~1"
shift /1

if "%SCRIPT_NAME%"=="keygen" (
    cd /d "%SCRIPT_DIR%\go\keygen"
    go run main.go %*
    goto :eof
)

if "%SCRIPT_NAME%"=="generate-tls-cert" (
    cd /d "%SCRIPT_DIR%\shell"
    call generate-tls-cert.bat %*
    goto :eof
)

if "%SCRIPT_NAME%"=="--help" goto :print_usage
if "%SCRIPT_NAME%"=="-h" goto :print_usage

echo Unknown script: %SCRIPT_NAME%
echo.
call :print_usage
exit /b 1

:print_usage
echo RealEntity Node Script Runner
echo.
echo Usage: %~nx0 ^<category^>/^<script^> [args...]
echo.
echo Available scripts:
echo.
echo Go Scripts (scripts/go/):
echo   keygen                 - Generate private keys and peer IDs
echo.
echo Shell Scripts (scripts/shell/):
echo   generate-tls-cert      - Generate TLS certificates
echo.
echo Examples:
echo   %~nx0 keygen
echo   %~nx0 keygen -output config.json
echo   %~nx0 generate-tls-cert
echo.
echo For script-specific help, run:
echo   %~nx0 ^<script^> --help
goto :eof
