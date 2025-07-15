# RealEntity Node

A distributed peer-to-peer networking node built with Go and libp2p, featuring dynamic service discovery and cross-platform deployment capabilities.

## Features

- **Multi-Discovery Architecture**: mDNS for local development, Bootstrap nodes for VPS deployment, DHT for distributed networks
- **Dynamic Service Registry**: Auto-discovery and testing of services between connected peers
- **Cross-Platform Deployment**: Supports local development, Docker testing, VPS deployment, and Kubernetes orchestration
- **Unified Deployment System**: Single script handles all deployment scenarios with dry-run capabilities

## Quick Start

### Local Development

```bash
# Build and run locally with mDNS discovery
./deploy/universal.sh local

# Or manually
go run cmd/main.go
```

### Docker Testing

```bash
# Start multi-node test environment
./deploy/universal.sh docker

# View logs
docker-compose logs -f
```

### VPS Deployment

```bash
# Generate a bootstrap node identity (recommended for production)
./deploy/universal.sh generate-key

# Deploy bootstrap node with generated identity
./deploy/universal.sh vps-bootstrap --public-ip YOUR_PUBLIC_IP

# Deploy peer node
./deploy/universal.sh vps-peer --bootstrap-peer /ip4/IP/tcp/PORT/p2p/ID
```

## Project Structure

```text
cmd/
├── main.go                      # Application entry point
└── main_test.go                 # Application tests
internal/                        # Private application packages
├── config/                # Configuration management
├── discovery/             # Peer discovery mechanisms
├── node/                  # LibP2P host creation and management
├── protocol/              # P2P protocol handlers
├── services/              # Service framework and examples
└── utils/                 # Utility functions
scripts/                         # Helper utilities and tools
├── keygen/                # Private key generator
└── README.md              # Scripts documentation
deploy/                    # Deployment scripts and configurations
├── universal.sh          # Unified deployment script
├── k8s/                  # Kubernetes manifests
└── templates/            # Configuration templates
docs/                     # Documentation
├── README.md            # Detailed documentation
├── DEPLOYMENT_SUMMARY.md # Deployment overview
├── DISCOVERY.md         # Discovery mechanism details
└── *.md                 # Other documentation
tests/                   # Testing utilities and scripts
├── README.md           # Testing documentation
├── console/            # Console-based testing tools
├── api/                # HTTP API testing scripts
└── examples/           # Example test cases and data
```

## Testing

### Quick Test Run

```bash
# Windows
run_tests.bat

# Linux/macOS
./run_tests.sh
```

### Test Types

- **Unit Tests**: `go test ./cmd -v` or `run_tests.bat unit`
- **Console Tests**: `run_tests.bat console` - Test services via command line
- **API Tests**: `run_tests.bat api` - Test HTTP API endpoints (requires running server)

See [`tests/README.md`](tests/README.md) for detailed testing documentation.

## Development

### Building

```bash
# Build the application
go build -o realentity-node cmd/main.go

# Or use Makefile
make build
```

### Testing

```bash
# Run tests
go test ./...

# Or use Makefile
make test
```

### Code Quality

```bash
# Run linting
make lint

# Format code
make fmt
```

## Configuration

The application uses JSON configuration with environment variable overrides:

```json
{
  "private_key": "base64-encoded-private-key",
  "discovery": {
    "enable_mdns": true,
    "enable_bootstrap": false,
    "bootstrap_peers": []
  },
  "server": {
    "bind_address": "0.0.0.0",
    "port": 4001,
    "public_ip": "1.2.3.4"
  }
}
```

### Hardcoded Identity (Bootstrap Nodes)

For bootstrap nodes, you can generate a consistent peer ID:

```bash
# Generate new identity
go run scripts/keygen/main.go -generate-key -output bootstrap-config.json

# Use the generated config
cp bootstrap-config.json config.json
go run cmd/main.go
```

Environment variables override config values using the pattern `REALENTITY_SECTION_KEY`.

## Documentation

Detailed documentation is available in the [`docs/`](docs/) directory:

- [Complete README](docs/README.md) - Comprehensive project documentation
- [Deployment Guide](docs/DEPLOYMENT_SUMMARY.md) - All deployment options
- [Discovery Architecture](docs/DISCOVERY.md) - Peer discovery mechanisms
- [Universal Deployment](docs/UNIVERSAL_DEPLOYMENT.md) - Unified deployment system

## License

[License information]

## Contributing

[Contributing guidelines]
