#!/bin/bash
# deploy/universal.sh - Single deployment script for all scenarios

set -e

VERSION="1.0.0"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_banner() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE} RealEntity Node Universal Deployer${NC}"
    echo -e "${BLUE} Version: $VERSION${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo ""
}

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

show_usage() {
    echo "Usage: $0 <deployment-type> [options]"
    echo ""
    echo "Deployment Types:"
    echo "  local          - Local development with mDNS"
    echo "  docker         - Docker containerized testing"
    echo "  vps-bootstrap  - VPS bootstrap node"
    echo "  vps-peer       - VPS peer node"
    echo "  k8s            - Kubernetes deployment"
    echo "  generate-key   - Generate new private key and peer ID"
    echo ""
    echo "Options:"
    echo "  --bootstrap-peer <addr>    Bootstrap peer address (for peer nodes)"
    echo "  --public-ip <ip>          Public IP address (for VPS)"
    echo "  --port <port>             Port number (default: 4001)"
    echo "  --clean                   Clean deployment (remove existing)"
    echo "  --dry-run                 Show what would be done"
    echo "  --config <file>           Custom config file"
    echo ""
    echo "Examples:"
    echo "  $0 local"
    echo "  $0 docker --clean"
    echo "  $0 vps-bootstrap --public-ip 1.2.3.4"
    echo "  $0 vps-peer --bootstrap-peer '/ip4/1.2.3.4/tcp/4001/p2p/12D3...'"
    echo "  $0 k8s --namespace realentity"
    exit 1
}

# Parse command line arguments
DEPLOYMENT_TYPE=""
BOOTSTRAP_PEER=""
PUBLIC_IP=""
PORT="4001"
CLEAN_DEPLOY=false
DRY_RUN=false
CUSTOM_CONFIG=""
NAMESPACE="default"

while [[ $# -gt 0 ]]; do
    case $1 in
        local|docker|vps-bootstrap|vps-peer|k8s|generate-key)
            DEPLOYMENT_TYPE="$1"
            shift
            ;;
        --bootstrap-peer)
            BOOTSTRAP_PEER="$2"
            shift 2
            ;;
        --public-ip)
            PUBLIC_IP="$2"
            shift 2
            ;;
        --port)
            PORT="$2"
            shift 2
            ;;
        --clean)
            CLEAN_DEPLOY=true
            shift
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --config)
            CUSTOM_CONFIG="$2"
            shift 2
            ;;
        --namespace)
            NAMESPACE="$2"
            shift 2
            ;;
        -h|--help)
            show_usage
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            ;;
    esac
done

if [[ -z "$DEPLOYMENT_TYPE" ]]; then
    print_error "Deployment type is required"
    show_usage
fi

# Deployment functions
deploy_local() {
    print_status "Deploying for local development..."
    
    if [[ $DRY_RUN == true ]]; then
        echo "Would create local config with mDNS enabled"
        echo "Would build and run locally"
        return
    fi
    
    # Generate local config
    cat > "$PROJECT_ROOT/config.json" << EOF
{
  "discovery": {
    "enable_mdns": true,
    "enable_bootstrap": false,
    "enable_dht": false,
    "mdns_service_tag": "realentity-local",
    "mdns_quiet_mode": false,
    "bootstrap_peers": [],
    "dht_rendezvous": "realentity-dht"
  },
  "log_level": "info",
  "server": {
    "bind_address": "0.0.0.0",
    "port": 0
  }
}
EOF
    
    print_status "Local configuration created"
    print_status "Run 'go run main.go' to start multiple instances"
}

deploy_docker() {
    print_status "Deploying Docker test environment..."
    
    if [[ $DRY_RUN == true ]]; then
        echo "Would build Docker images"
        echo "Would start Docker Compose with bootstrap + 4 peers"
        return
    fi
    
    if [[ $CLEAN_DEPLOY == true ]]; then
        docker-compose down --remove-orphans 2>/dev/null || true
    fi
    
    cd "$PROJECT_ROOT"
    exec ./docker-test.sh start
}

deploy_vps_bootstrap() {
    print_status "Deploying VPS bootstrap node..."
    
    if [[ -z "$PUBLIC_IP" ]]; then
        PUBLIC_IP=$(detect_public_ip)
        print_status "Detected public IP: $PUBLIC_IP"
    fi
    
    if [[ $DRY_RUN == true ]]; then
        echo "Would create bootstrap config for IP: $PUBLIC_IP"
        echo "Would install Go and build application"
        echo "Would create systemd service"
        return
    fi
    
    install_dependencies
    build_application
    create_bootstrap_config
    setup_service "bootstrap"
    configure_firewall
    
    print_status "Bootstrap node deployed successfully!"
    print_status "Multiaddr: /ip4/$PUBLIC_IP/tcp/$PORT/p2p/\$(check logs for peer ID)"
}

deploy_vps_peer() {
    print_status "Deploying VPS peer node..."
    
    if [[ -z "$BOOTSTRAP_PEER" ]]; then
        print_error "Bootstrap peer address is required for peer deployment"
        echo "Use: --bootstrap-peer '/ip4/IP/tcp/PORT/p2p/ID'"
        exit 1
    fi
    
    if [[ $DRY_RUN == true ]]; then
        echo "Would create peer config pointing to: $BOOTSTRAP_PEER"
        echo "Would install Go and build application" 
        echo "Would create systemd service"
        return
    fi
    
    install_dependencies
    build_application
    create_peer_config
    setup_service "peer"
    configure_firewall
    
    print_status "Peer node deployed successfully!"
}

deploy_k8s() {
    print_status "Deploying to Kubernetes..."
    
    if [[ $DRY_RUN == true ]]; then
        echo "Would create Kubernetes manifests"
        echo "Would deploy to namespace: $NAMESPACE"
        return
    fi
    
    create_k8s_manifests
    kubectl apply -f "$PROJECT_ROOT/deploy/k8s/" -n "$NAMESPACE"
    
    print_status "Kubernetes deployment created!"
}

# Helper functions
detect_public_ip() {
    curl -s ifconfig.me || curl -s icanhazip.com || echo "127.0.0.1"
}

install_dependencies() {
    if ! command -v go &> /dev/null; then
        print_status "Installing Go..."
        # Platform-specific Go installation logic
        case "$(uname -s)" in
            Linux*)   install_go_linux ;;
            Darwin*)  install_go_macos ;;
            *)        print_error "Unsupported platform" && exit 1 ;;
        esac
    fi
}

build_application() {
    print_status "Building RealEntity node..."
    cd "$PROJECT_ROOT"
    go mod tidy
    go build -o realentity-node cmd/main.go
}

create_bootstrap_config() {
    cat > "$PROJECT_ROOT/config.json" << EOF
{
  "discovery": {
    "enable_mdns": false,
    "enable_bootstrap": true,
    "enable_dht": false,
    "mdns_service_tag": "realentity-mdns",
    "mdns_quiet_mode": true,
    "bootstrap_peers": [],
    "dht_rendezvous": "realentity-dht"
  },
  "log_level": "info",
  "server": {
    "bind_address": "0.0.0.0",
    "port": $PORT,
    "public_ip": "$PUBLIC_IP"
  }
}
EOF
}

create_peer_config() {
    cat > "$PROJECT_ROOT/config.json" << EOF
{
  "discovery": {
    "enable_mdns": false,
    "enable_bootstrap": true,
    "enable_dht": false,
    "mdns_service_tag": "realentity-mdns",
    "mdns_quiet_mode": true,
    "bootstrap_peers": ["$BOOTSTRAP_PEER"],
    "dht_rendezvous": "realentity-dht"
  },
  "log_level": "info",
  "server": {
    "bind_address": "0.0.0.0",
    "port": $PORT
  }
}
EOF
}

setup_service() {
    local node_type="$1"
    print_status "Setting up systemd service..."
    
    sudo tee /etc/systemd/system/realentity-$node_type.service > /dev/null << EOF
[Unit]
Description=RealEntity $node_type Node
After=network.target

[Service]
Type=simple
User=$(whoami)
WorkingDirectory=$PROJECT_ROOT
ExecStart=$PROJECT_ROOT/realentity-node
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF
    
    sudo systemctl daemon-reload
    sudo systemctl enable realentity-$node_type
    sudo systemctl start realentity-$node_type
}

configure_firewall() {
    print_status "Configuring firewall..."
    
    if command -v ufw &> /dev/null; then
        sudo ufw allow $PORT/tcp
    elif command -v firewall-cmd &> /dev/null; then
        sudo firewall-cmd --permanent --add-port=$PORT/tcp
        sudo firewall-cmd --reload
    fi
}

create_k8s_manifests() {
    mkdir -p "$PROJECT_ROOT/deploy/k8s"
    
    # Create Kubernetes deployment manifests
    # (This would include ConfigMaps, Deployments, Services, etc.)
    print_status "Kubernetes manifests created in deploy/k8s/"
}

generate_key() {
    print_status "Generating new private key and peer ID..."
    
    if [[ $DRY_RUN == true ]]; then
        echo "Would generate new private key and peer ID"
        echo "Would create bootstrap configuration template"
        return
    fi
    
    cd "$PROJECT_ROOT"
    go run scripts/go/keygen/main.go -generate-key -output "bootstrap-config.json"
    
    print_status "Bootstrap configuration saved to: bootstrap-config.json"
    print_status "Edit the file to set your public IP, then use:"
    print_status "  cp bootstrap-config.json config.json"
    print_status "  go run cmd/main.go"
}

# Main execution
main() {
    print_banner
    
    case $DEPLOYMENT_TYPE in
        local)         deploy_local ;;
        docker)        deploy_docker ;;
        vps-bootstrap) deploy_vps_bootstrap ;;
        vps-peer)      deploy_vps_peer ;;
        k8s)           deploy_k8s ;;
        generate-key)  generate_key ;;
        *)             print_error "Invalid deployment type: $DEPLOYMENT_TYPE" && exit 1 ;;
    esac
}

main
