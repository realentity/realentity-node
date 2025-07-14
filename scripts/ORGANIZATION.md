# Scripts Organization Summary

## New Structure

The scripts directory has been reorganized into a clear, maintainable structure:

```
scripts/
├── go/                     # Go-based utility programs
│   ├── keygen/            # Private key and peer ID generator
│   │   └── main.go
│   └── README.md
├── shell/                  # Shell scripts (bash/batch)
│   ├── generate-tls-cert.sh   # TLS cert generation (Unix)
│   ├── generate-tls-cert.bat  # TLS cert generation (Windows)
│   └── README.md
├── utils/                  # General utility scripts
│   ├── config-validator.py    # Configuration validator example
│   └── README.md
├── run.sh                  # Script runner (Unix)
├── run.bat                 # Script runner (Windows)
└── README.md              # Main documentation
```

## Quick Usage

### Using the Script Runner
```bash
# Unix/Linux/macOS
./scripts/run.sh keygen
./scripts/run.sh generate-tls-cert --domain localhost

# Windows
scripts\run.bat keygen
scripts\run.bat generate-tls-cert
```

### Direct Execution
```bash
# Go scripts
cd scripts/go/keygen && go run main.go -generate-key

# Shell scripts  
cd scripts/shell && ./generate-tls-cert.sh

# Utility scripts
python3 scripts/utils/config-validator.py config.json
```

## Benefits of New Organization

1. **Clear Separation**: Different script types are in separate directories
2. **Scalability**: Easy to add new scripts without clutter
3. **Documentation**: Each directory has specific guidelines
4. **Cross-Platform**: Both Unix and Windows scripts provided
5. **Unified Interface**: Script runner provides consistent access
6. **Best Practices**: Guidelines for adding new scripts

## Migration Changes

- Moved `scripts/keygen/` → `scripts/go/keygen/`
- Moved TLS cert scripts to `scripts/shell/`
- Updated `deploy/universal.sh` to use new paths
- Added comprehensive documentation
- Created script runner utilities

## Adding New Scripts

### Go Scripts
1. Create subdirectory under `scripts/go/`
2. Include `main.go` with proper CLI interface
3. Add to script runner if commonly used

### Shell Scripts  
1. Add to `scripts/shell/`
2. Provide both `.sh` and `.bat` versions
3. Follow naming conventions

### Utility Scripts
1. Add to `scripts/utils/`
2. Use appropriate file extension
3. Include proper documentation

This organization supports the project's growth while maintaining clean separation of concerns and easy discoverability.
