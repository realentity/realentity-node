# Shell Scripts

This directory contains shell scripts for development and deployment tasks.

## Available Scripts

### TLS Certificate Generation
- `generate-tls-cert.sh` - Unix/Linux/macOS version
- `generate-tls-cert.bat` - Windows version

Generates self-signed TLS certificates for HTTPS development and testing.

**Usage:**
```bash
# Unix/Linux/macOS
chmod +x generate-tls-cert.sh
./generate-tls-cert.sh --domain localhost
./generate-tls-cert.sh --domain yourdomain.com --days 3650

# Windows
generate-tls-cert.bat
```

## Script Guidelines

### Cross-Platform Support
- Provide both `.sh` and `.bat` versions when possible
- Use portable commands when feasible
- Document platform-specific requirements

### Script Structure
```bash
#!/bin/bash
# script-name.sh - Brief description

set -e  # Exit on error

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Functions
print_usage() {
    echo "Usage: $0 [options]"
    # ... usage details
}

# Main script logic
# ...
```

### Error Handling
- Use `set -e` to exit on errors
- Provide meaningful error messages
- Check for required dependencies
- Validate input parameters

### Documentation
- Include usage information at the top
- Document all command-line options
- Provide examples
- Note any dependencies

## Adding New Shell Scripts

1. Create both `.sh` and `.bat` versions when possible
2. Follow naming convention: `kebab-case-name.sh`
3. Include proper shebang line: `#!/bin/bash`
4. Add executable permissions: `chmod +x script.sh`
5. Test on target platforms
6. Update this README

## Dependencies

Common dependencies for shell scripts:
- `openssl` - For TLS certificate generation
- `curl` - For HTTP requests
- `jq` - For JSON processing
- `docker` - For container operations

Check and document dependencies in each script.
