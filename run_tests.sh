#!/bin/bash

# RealEntity Node Test Runner
# Usage: ./run_tests.sh [test-type]
# 
# Test types:
#   console    - Run console-based tests
#   api        - Run API tests (requires server running)
#   unit       - Run Go unit tests
#   all        - Run all tests

set -e

echo "RealEntity Node Test Runner"
echo "=============================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed or not in PATH"
    exit 1
fi

# Function to run console tests
run_console_tests() {
    echo "Running console tests..."
    cd tests/console
    
    echo "Testing text.process service..."
    go run test_text_service.go "hello world" uppercase
    go run test_text_service.go "TESTING" lowercase
    go run test_text_service.go "abcd" reverse
    
    cd ../..
    echo "Console tests completed"
}

# Function to run API tests
run_api_tests() {
    echo "Running API tests..."
    
    # Check if server is running
    if ! curl -s http://localhost:8080/health > /dev/null; then
        echo "Server is not running on port 8080"
        echo "Please start the server with: go run cmd/main.go"
        exit 1
    fi
    
    cd tests/api
    
    # Check if PowerShell is available (Windows)
    if command -v powershell &> /dev/null; then
        powershell -ExecutionPolicy Bypass -File test_api.ps1
    elif command -v curl &> /dev/null; then
        # Run batch file on Windows or use curl directly on Unix
        if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
            ./test_api.bat
        else
            echo "Running curl tests..."
            curl -X POST http://localhost:8080/api/services/execute \
                -H "Content-Type: application/json" \
                -d '{"service":"text.process","payload":{"text":"hello world","operation":"uppercase"}}'
        fi
    else
        echo "Neither PowerShell nor curl is available"
        exit 1
    fi
    
    cd ../..
    echo "API tests completed"
}

# Function to run unit tests
run_unit_tests() {
    echo "Running unit tests..."
    go test ./cmd -v
    go test ./internal/... -v
    echo "Unit tests completed"
}

# Main logic
case "${1:-all}" in
    console)
        run_console_tests
        ;;
    api)
        run_api_tests
        ;;
    unit)
        run_unit_tests
        ;;
    all)
        echo "Running all tests..."
        run_unit_tests
        run_console_tests
        echo ""
        echo "To run API tests, start the server and run:"
        echo "   ./run_tests.sh api"
        ;;
    *)
        echo "Usage: $0 [console|api|unit|all]"
        echo ""
        echo "Test types:"
        echo "  console  - Run console-based service tests"
        echo "  api      - Run HTTP API tests (requires server running)"
        echo "  unit     - Run Go unit tests"
        echo "  all      - Run unit and console tests"
        exit 1
        ;;
esac

echo ""
echo "Testing completed successfully!"
