@echo off
REM RealEntity Node Test Runner for Windows
REM Usage: run_tests.bat [test-type]
REM
REM Test types:
REM   console    - Run console-based tests
REM   api        - Run API tests (requires server running)
REM   unit       - Run Go unit tests
REM   all        - Run all tests

setlocal enabledelayedexpansion

echo RealEntity Node Test Runner
echo ==============================

REM Check if Go is installed
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo ERROR: Go is not installed or not in PATH
    exit /b 1
)

set TEST_TYPE=%1
if "%TEST_TYPE%"=="" set TEST_TYPE=all

if "%TEST_TYPE%"=="console" goto :run_console_tests
if "%TEST_TYPE%"=="api" goto :run_api_tests
if "%TEST_TYPE%"=="unit" goto :run_unit_tests
if "%TEST_TYPE%"=="all" goto :run_all_tests

echo Usage: %0 [console^|api^|unit^|all]
echo.
echo Test types:
echo   console  - Run console-based service tests
echo   api      - Run HTTP API tests (requires server running)
echo   unit     - Run Go unit tests
echo   all      - Run unit and console tests
exit /b 1

:run_console_tests
echo Running console tests...
cd tests\console

echo Testing text.process service...
go run test_text_service.go "hello world" uppercase
go run test_text_service.go "TESTING" lowercase
go run test_text_service.go "abcd" reverse

cd ..\..
echo Console tests completed
goto :end

:run_api_tests
echo Running API tests...

REM Check if server is running
curl -s http://localhost:8080/health >nul 2>nul
if %errorlevel% neq 0 (
    echo Server is not running on port 8080
    echo Please start the server with: go run cmd/main.go
    exit /b 1
)

cd tests\api
powershell -ExecutionPolicy Bypass -File test_api.ps1
cd ..\..
echo API tests completed
goto :end

:run_unit_tests
echo Running unit tests...
go test ./cmd -v
go test ./internal/... -v
echo Unit tests completed
goto :end

:run_all_tests
echo Running all tests...
call :run_unit_tests
call :run_console_tests
echo.
echo To run API tests, start the server and run:
echo    run_tests.bat api
goto :end

:end
echo.
echo Testing completed successfully!
