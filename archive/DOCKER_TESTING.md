# Docker Testing Environment

This directory contains a complete Docker-based testing environment for RealEntity nodes that simulates a distributed production environment locally.

## Architecture

The test environment creates:
- **1 Bootstrap Node** (172.20.0.10) - Entry point for the network
- **4 Peer Nodes** (172.20.0.11-14) - Regular network participants
- **Custom Docker Network** (172.20.0.0/16) - Isolated network environment

## Quick Start

### Prerequisites
- Docker installed and running
- Docker Compose installed

### Start the Test Network

**Linux/macOS:**
```bash
chmod +x docker-test.sh docker-test-services.sh
./docker-test.sh start
```

**Windows:**
```cmd
docker-test.bat start
```

### Monitor the Network

**View all logs in real-time:**
```bash
docker-compose logs -f
```

**Test services between nodes:**
```bash
./docker-test-services.sh all
```

**Monitor peer discovery:**
```bash
./docker-test-services.sh monitor
```

### Stop the Network

```bash
./docker-test.sh stop
```

## What Gets Tested

### 1. Bootstrap Discovery
- Bootstrap node starts and becomes network entry point
- Peer nodes connect to bootstrap automatically
- Peer discovery through bootstrap peerstore

### 2. Service Auto-Testing
- Nodes automatically test services when they discover peers
- Echo service testing with unique node IDs
- Text processing service testing

### 3. Network Resilience
- Connection monitoring and reconnection
- Bootstrap peer health checks
- Network partition recovery

## Network Flow

```
1. Bootstrap starts with empty peer list
2. Peer1 connects to Bootstrap → Bootstrap learns about Peer1
3. Peer2 connects to Bootstrap → Learns about Peer1, connects directly
4. Peer3 connects to Bootstrap → Learns about Peer1&2, connects directly
5. Peer4 connects to Bootstrap → Learns about all peers, connects directly
6. All peers auto-test services on discovered peers
```

## Files Structure

```
configs/
├── bootstrap-config.json   # Bootstrap node configuration
└── peer-config.json       # Peer node configuration template

docker-compose.yml          # Multi-container setup
Dockerfile                  # Node container definition
docker-test.sh             # Linux/macOS test runner
docker-test.bat            # Windows test runner
docker-test-services.sh    # Service testing utilities
```

## Configuration Details

### Bootstrap Node Config
```json
{
  "discovery": {
    "enable_bootstrap": true,
    "bootstrap_peers": [],          // Empty - it's the bootstrap
    "enable_mdns": false           // Disabled for container environment
  },
  "server": {
    "public_ip": "172.20.0.10"     // Fixed container IP
  }
}
```

### Peer Node Config
```json
{
  "discovery": {
    "enable_bootstrap": true,
    "bootstrap_peers": [
      "/ip4/172.20.0.10/tcp/4001/p2p/BOOTSTRAP_ID"  // Points to bootstrap
    ],
    "enable_mdns": false
  }
}
```

## Testing Scenarios

### Basic Connectivity Test
```bash
./docker-test-services.sh connectivity
```

### Service Communication Test
```bash
./docker-test-services.sh services
```

### Network Partition Test
```bash
# Stop bootstrap temporarily
docker stop realentity-bootstrap
sleep 30
docker start realentity-bootstrap
# Observe reconnection in logs
```

### Scale Test
Add more peer nodes by editing `docker-compose.yml`:
```yaml
peer5:
  build: .
  container_name: realentity-peer5
  networks:
    realentity-net:
      ipv4_address: 172.20.0.15
  # ... same config as other peers
```

## Expected Output

When working correctly, you should see logs like:
```
realentity-bootstrap  | Node started with ID: 12D3KooW...
realentity-peer1      | Connected to bootstrap peer: 12D3KooW...
realentity-peer2      | Discovered peer: 12D3KooW...peer1...
realentity-peer2      | Auto-testing services on newly discovered peer
realentity-peer1      | Echo response: Hello from 12D3KooW...peer2...
```

## Troubleshooting

### Bootstrap ID Not Found
If peer nodes can't get bootstrap ID:
```bash
# Check bootstrap logs
docker logs realentity-bootstrap

# Manually get bootstrap ID
docker logs realentity-bootstrap 2>&1 | grep "Node started with ID:"
```

### No Peer Connections
```bash
# Check network connectivity
docker exec realentity-peer1 nc -z bootstrap 4001

# Check discovery logs
docker logs realentity-peer1 | grep -i discovery
```

### Port Conflicts
If port 4001 is busy:
```bash
# Change exposed port in docker-compose.yml
ports:
  - "4002:4001"  # Map to different host port
```

## Production Simulation

This setup accurately simulates:
- **Cross-VPS discovery** through bootstrap nodes
- **Service auto-testing** between distributed nodes  
- **Network resilience** and reconnection scenarios
- **Peer management** and reliability scoring

The containerized environment provides isolation while maintaining realistic network communication patterns that mirror actual VPS deployments.
