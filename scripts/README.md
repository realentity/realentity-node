# Scripts Directory

This directory contains various utility scripts organized by type and purpose.

## Directory Structure

```
scripts/
├── go/                 # Go-based utility programs
│   └── keygen/         # Private key and peer ID generator
├── shell/              # Shell scripts (bash/batch)
│   ├── generate-tls-cert.sh   # TLS certificate generation (Unix)
│   └── generate-tls-cert.bat  # TLS certificate generation (Windows)
├── utils/              # General utility scripts
└── README.md           # This file
```

## Go Scripts (`scripts/go/`)

### Keygen Tool
- **Location**: `scripts/go/keygen/`
- **Purpose**: Generate consistent private keys and peer IDs for bootstrap nodes
- **Usage**: 
  ```bash
  cd scripts/go/keygen
  go run main.go -generate-key
  go run main.go -generate-key -output bootstrap-config.json
  ```

## Shell Scripts (`scripts/shell/`)

### TLS Certificate Generation
- **Files**: `generate-tls-cert.sh` (Unix), `generate-tls-cert.bat` (Windows)
- **Purpose**: Generate self-signed TLS certificates for HTTPS development
- **Usage**:
  ```bash
  # Unix/Linux/macOS
  cd scripts/shell
  chmod +x generate-tls-cert.sh
  ./generate-tls-cert.sh --domain yourdomain.com
  
  # Windows
  cd scripts\shell
  generate-tls-cert.bat
  ```

## Utility Scripts (`scripts/utils/`)

Directory for general purpose utility scripts that don't fit into other categories.

## Script Development Guidelines

### For Go Scripts:
1. Create a subdirectory under `scripts/go/`
2. Include a `main.go` file as the entry point
3. Use Go modules if the script has dependencies
4. Include usage documentation in comments
5. Example structure:
   ```text
   scripts/go/mytool/
   ├── main.go
   ├── go.mod (if needed)
   └── README.md (optional)
   ```

### For Shell Scripts:
1. Place scripts in `scripts/shell/`
2. Provide both Unix (.sh) and Windows (.bat) versions when possible
3. Include usage documentation at the top of the script
4. Make scripts executable: `chmod +x script.sh`

### For Utility Scripts:
1. Place in `scripts/utils/`
2. Use appropriate file extensions (.py, .js, .rb, etc.)
3. Include shebang lines for interpreted scripts
4. Document dependencies and usage

## Best Practices

1. **Naming**: Use kebab-case for script names (e.g., `generate-tls-cert.sh`)
2. **Documentation**: Include usage examples and parameter descriptions
3. **Error Handling**: Implement proper error checking and user feedback
4. **Platform Support**: Consider cross-platform compatibility
5. **Dependencies**: Document any external dependencies clearly
6. **Testing**: Test scripts on target platforms before committing

## Adding New Scripts

To add new helper utilities:

1. Create a new directory: `scripts/your-utility/`
2. Add `main.go` with your utility logic
3. Update this README with usage instructions
4. Add any references to deployment scripts if needed

## Integration with Main Application

- Go scripts can be built and distributed with the application
- Shell scripts are typically used for deployment and development tasks
- All scripts should be callable from the project root directory

## Examples of Future Scripts

- **`dbmigrate/`** - Database migration utilities
- **`benchmark/`** - Performance testing scripts  
- **`monitor/`** - Network monitoring tools
- **`deploy-helper/`** - Advanced deployment utilities
- **`config-validator/`** - Configuration validation tools

## Usage in Deployment

The deployment system (`deploy/universal.sh`) can reference these scripts:

```bash
# Example: Universal script using keygen
./deploy/universal.sh generate-key
```
