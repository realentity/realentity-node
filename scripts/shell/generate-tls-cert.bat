@echo off
REM scripts\generate-tls-cert.bat - Generate self-signed TLS certificates for HTTPS (Windows)

setlocal enabledelayedexpansion

set "SCRIPT_DIR=%~dp0"
set "PROJECT_ROOT=%SCRIPT_DIR%\.."
set "CERTS_DIR=%PROJECT_ROOT%\certs"

REM Default values
set "DOMAIN=localhost"
set "DAYS=365"
set "COUNTRY=US"
set "STATE=CA"
set "CITY=San Francisco"
set "ORG=RealEntity"
set "OU=Development"

REM Create certs directory if it doesn't exist
if not exist "%CERTS_DIR%" mkdir "%CERTS_DIR%"

echo Generating TLS certificate for domain: %DOMAIN%
echo Certificate will be valid for %DAYS% days
echo.

REM Check if OpenSSL is available
where openssl >nul 2>&1
if errorlevel 1 (
    echo Error: OpenSSL is not installed or not in PATH
    echo Please install OpenSSL from: https://slproweb.com/products/Win32OpenSSL.html
    echo Or use Git Bash which includes OpenSSL
    pause
    exit /b 1
)

REM Generate private key
openssl genrsa -out "%CERTS_DIR%\server.key" 2048

REM Generate certificate signing request
openssl req -new -key "%CERTS_DIR%\server.key" -out "%CERTS_DIR%\server.csr" -subj "/C=%COUNTRY%/ST=%STATE%/L=%CITY%/O=%ORG%/OU=%OU%/CN=%DOMAIN%"

REM Create config file for extensions
(
echo [v3_req]
echo keyUsage = keyEncipherment, dataEncipherment
echo extendedKeyUsage = serverAuth
echo subjectAltName = @alt_names
echo.
echo [alt_names]
echo DNS.1 = %DOMAIN%
echo DNS.2 = localhost
echo DNS.3 = *.localhost
echo IP.1 = 127.0.0.1
echo IP.2 = ::1
) > "%CERTS_DIR%\cert_extensions.conf"

REM Generate self-signed certificate
openssl x509 -req -in "%CERTS_DIR%\server.csr" -signkey "%CERTS_DIR%\server.key" -out "%CERTS_DIR%\server.crt" -days %DAYS% -extensions v3_req -extfile "%CERTS_DIR%\cert_extensions.conf"

REM Clean up temporary files
del "%CERTS_DIR%\server.csr"
del "%CERTS_DIR%\cert_extensions.conf"

echo.
echo  TLS certificate generated successfully!
echo.
echo Files created:
echo   Certificate: %CERTS_DIR%\server.crt
echo   Private Key: %CERTS_DIR%\server.key
echo.
echo To use HTTPS, update your config.json:
echo {
echo   "server": {
echo     "https_port": 8443,
echo     "tls_cert_file": "certs/server.crt",
echo     "tls_key_file": "certs/server.key"
echo   }
echo }
echo.
echo ️  Note: This is a self-signed certificate for development only.
echo    For production, use certificates from a trusted CA like Let's Encrypt.
echo.
pause
