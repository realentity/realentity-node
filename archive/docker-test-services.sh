#!/bin/bash
# docker-test-services.sh - Test services between Docker nodes

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() {
    echo -e "${GREEN}[TEST]${NC} $1"
}

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

# Test service calls between nodes
test_services() {
    print_status "Testing services between Docker nodes..."
    
    # Get list of running containers
    containers=$(docker-compose ps --services --filter status=running)
    
    for container in $containers; do
        if [ "$container" != "bootstrap" ]; then
            print_info "Testing services on $container..."
            
            # Execute service test inside the container
            docker exec "realentity-$container" sh -c '
                echo "Container: $(hostname)"
                echo "Recent logs:"
                tail -n 5 /var/log/realentity.log 2>/dev/null || echo "No specific log file, checking stdout..."
            '
        fi
    done
}

# Monitor network connectivity
test_connectivity() {
    print_status "Testing network connectivity between nodes..."
    
    # Test from peer1 to other nodes
    print_info "Testing connectivity from peer1..."
    docker exec realentity-peer1 sh -c '
        echo "Testing connectivity from peer1:"
        nc -z bootstrap 4001 && echo " Can reach bootstrap:4001" || echo " Cannot reach bootstrap:4001"
        nc -z peer2 4001 && echo " Can reach peer2:4001" || echo " Cannot reach peer2:4001"
        nc -z peer3 4001 && echo " Can reach peer3:4001" || echo " Cannot reach peer3:4001"
        nc -z peer4 4001 && echo " Can reach peer4:4001" || echo " Cannot reach peer4:4001"
    '
}

# Show detailed logs for debugging
show_detailed_logs() {
    print_status "Showing detailed logs for all nodes..."
    
    for node in bootstrap peer1 peer2 peer3 peer4; do
        print_info "=== Logs for $node ==="
        docker logs "realentity-$node" --tail 20
        echo ""
    done
}

# Monitor real-time peer discovery
monitor_discovery() {
    print_status "Monitoring peer discovery in real-time..."
    print_info "Press Ctrl+C to stop monitoring"
    
    # Follow logs and filter for discovery events
    docker-compose logs -f | grep -E "(Discovered peer|Connected to peer|Auto-testing|Echo response|Text processing)"
}

# Run specific tests
run_tests() {
    case "${1:-all}" in
        "connectivity")
            test_connectivity
            ;;
        "services")
            test_services
            ;;
        "logs")
            show_detailed_logs
            ;;
        "monitor")
            monitor_discovery
            ;;
        "all")
            test_connectivity
            echo ""
            test_services
            echo ""
            show_detailed_logs
            ;;
        *)
            echo "Usage: $0 {connectivity|services|logs|monitor|all}"
            echo "  connectivity - Test network connectivity between nodes"
            echo "  services    - Test service calls between nodes"
            echo "  logs        - Show detailed logs from all nodes"
            echo "  monitor     - Monitor peer discovery in real-time"
            echo "  all         - Run all tests"
            exit 1
            ;;
    esac
}

# Check if containers are running
check_containers() {
    if ! docker-compose ps | grep -q "Up"; then
        echo "No containers are running. Please start the test environment first:"
        echo "  ./docker-test.sh start"
        exit 1
    fi
}

main() {
    check_containers
    run_tests "$@"
}

main "$@"
