# RealEntity Node Tests

This directory contains various testing utilities and scripts for the RealEntity Node project.

## Directory Structure

```
tests/
├── README.md           # This file - testing documentation
├── console/            # Console-based testing tools
│   ├── test_text_service.go    # Command-line service tester
│   └── interactive_text.go     # Interactive service tester
├── api/                # HTTP API testing tools
│   ├── test_api.bat           # Windows batch script for API testing
│   └── test_api.ps1           # PowerShell script for API testing
└── examples/           # Example test scenarios and data
```

## Available Tests

### 1. Console Tests (`tests/console/`)

#### Command-Line Service Tester
**File:** `test_text_service.go`

Test individual services from the command line:
```bash
cd tests/console
go run test_text_service.go "hello world" uppercase
go run test_text_service.go "HELLO WORLD" lowercase
go run test_text_service.go "hello" reverse
```

#### Interactive Service Tester
**File:** `interactive_text.go`

Run an interactive console to test services:
```bash
cd tests/console
go run interactive_text.go
```

### 2. API Tests (`tests/api/`)

#### Windows Batch Script
**File:** `test_api.bat`

Run comprehensive API tests using curl:
```cmd
cd tests\api
test_api.bat
```

#### PowerShell Script
**File:** `test_api.ps1`

Run API tests with PowerShell (recommended for Windows):
```powershell
cd tests\api
.\test_api.ps1
```

## Prerequisites

### For Console Tests
- Go 1.19+ installed
- RealEntity Node source code

### For API Tests
- Running RealEntity Node server (port 8080)
- For batch script: curl installed
- For PowerShell: PowerShell 5.0+ (included in Windows 10+)

## Quick Start

1. **Start the RealEntity Node server:**
   ```bash
   go run cmd/main.go
   ```

2. **Test services via console:**
   ```bash
   cd tests/console
   go run test_text_service.go "test message" uppercase
   ```

3. **Test services via API:**
   ```powershell
   cd tests/api
   .\test_api.ps1
   ```

## Available Services

### text.process
- **Operations:** `uppercase`, `lowercase`, `reverse`
- **Payload:**
  ```json
  {
    "text": "your text here",
    "operation": "uppercase|lowercase|reverse"
  }
  ```

### echo
- **Purpose:** Simple echo service for connectivity testing
- **Payload:**
  ```json
  {
    "message": "your message here"
  }
  ```

## API Endpoints

- `GET /health` - Health check
- `GET /api/services` - List available services
- `GET /api/node` - Node information
- `GET /api/peers` - Connected peers
- `POST /api/services/execute` - Execute a service

## Example API Request

```bash
curl -X POST http://localhost:8080/api/services/execute \
  -H "Content-Type: application/json" \
  -d '{
    "service": "text.process",
    "payload": {
      "text": "hello world",
      "operation": "uppercase"
    }
  }'
```

## Troubleshooting

### Common Issues

1. **Server not running:** Make sure `go run cmd/main.go` is running
2. **Port conflicts:** Check if port 8080 is already in use
3. **Service not found:** Verify the service name is correct (`text.process`, `echo`)
4. **Invalid operation:** Ensure operation is one of: `uppercase`, `lowercase`, `reverse`

### Checking Server Status

```bash
curl http://localhost:8080/health
```

Should return:
```json
{
  "status": "healthy",
  "services": ["echo", "text.process"],
  ...
}
```
