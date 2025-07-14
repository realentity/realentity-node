# Go Scripts

This directory contains Go-based utility programs for the RealEntity Node project.

## Available Tools

### Keygen (`keygen/`)

Generates consistent private keys and peer IDs for bootstrap nodes.

**Usage:**
```bash
cd scripts/go/keygen
go run main.go -generate-key
go run main.go -generate-key -output bootstrap-config.json
```

**Purpose:** Creating bootstrap nodes with known peer IDs for production deployment.

## Adding New Go Scripts

When creating new Go scripts:

1. Create a new subdirectory under `scripts/go/`
2. Include a `main.go` file with proper package structure
3. Use Go modules if external dependencies are needed
4. Include command-line argument parsing
5. Provide clear usage instructions
6. Handle errors gracefully

## Example Structure

```
scripts/go/mytool/
├── main.go
├── go.mod (if dependencies needed)
├── go.sum (if dependencies needed)
└── README.md (optional)
```

## Running Go Scripts

From project root:
```bash
# Direct execution
go run scripts/go/keygen/main.go [args]

# Or navigate to script directory
cd scripts/go/keygen
go run main.go [args]
```

## Building Go Scripts

To create standalone executables:
```bash
cd scripts/go/keygen
go build -o keygen main.go
```
