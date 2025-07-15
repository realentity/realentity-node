
### Commands
go mod init github.com/realentity/realentity-node

go clean -modcache
go mod tidy

---

##  Project Definition: **RealEntity Node**

###  Project Name:

`##  Current Implementation Status

###  Completed Features

* **Enhanced Discovery System**: Multi-mechanism peer discovery with mDNS, Bootstrap, and DHT framework
* **Discovery Manager**: Coordinates multiple discovery mechanisms with peer store and reliability scoring  
* **Configuration Management**: JSON-based config with auto-generation and sensible defaults
* **Service Registry**: Dynamic service registration and execution system with extensible framework
* **Protocol Handler**: Enhanced JSON-based request/response messaging with unique request IDs
* **Peer Store**: Metadata management with connection status, reliability scoring, and automatic cleanup
* **Bootstrap Discovery**: Global peer discovery through configurable bootstrap nodes
* **Enhanced mDNS**: Local network discovery with improved connection management
* **Example Services**: Echo and text processing services demonstrating the framework
* **Auto-Testing**: Automatic service testing and validation when peers are discovered
* **Client Utilities**: Service calling, peer management, and discovery utilities
* **Connection Management**: Automatic peer connection, status tracking, and reliability scoring

### ‍️ Current Capabilities

1. **Multi-Mechanism Discovery**: Find peers locally (mDNS) and globally (Bootstrap)
2. **Service Registration**: Dynamic registration and execution of distributed services
3. **Automatic Networking**: Nodes automatically find, connect, and test each other
4. **Configuration**: Easy setup via JSON config with auto-generated defaults
5. **Peer Management**: Track peer reliability, connection status, and service capabilities
6. **Real-time Monitoring**: Discovery statistics and peer status logging
7. **Manual Control**: Utilities for manual peer connection and service testing
---

##  Purpose

Build a **decentralized peer-to-peer (P2P) node** that runs on **servers (VPS)** and later on **mobile devices**, enabling communication, service execution, and distributed task routing **without centralized control**.

---

##  Core Objectives

*  Deployable **P2P nodes** using Go and libp2p
*  Dynamic peer **discovery** (LAN and public)
*  Secure **message exchange** over custom protocol
*  Local **service registry** per node (e.g., `resize.image`, `ai.infer`)
*  Optional **central bootstrap node** for global connectivity

---

## ️ Architecture Overview

```
+------------------------+           +------------------------+
|      Node A (VPS)      |<--------->|      Node B (Mobile)   |
| - Service: ai.infer    |           | - Service: resize.img  |
| - Discovery: DHT/mDNS  |           | - Secure stream        |
+------------------------+           +------------------------+

        ^                                ^
        |                                |
        +-------> Bootstrap node <-------+ 
                  (Optional VPS)
```

Each node:

* Registers local services
* Discovers other peers
* Accepts and processes incoming tasks
* Sends tasks to capable peers when needed

---

## ️ Tech Stack

| Layer         | Stack                                            |
| ------------- | ------------------------------------------------ |
| Language      | Go 1.22+                                         |
| P2P           | [go-libp2p](https://github.com/libp2p/go-libp2p) |
| Transport     | TCP, QUIC                                        |
| Discovery     | mDNS (local), DHT or Bootstrap (remote)          |
| Task Protocol | JSON over libp2p stream                          |
| Future        | Mobile embedding via `gomobile`, WASM support    |

---

## ️ Project Structure

```
realentity-node/
├── main.go                 # Entry point with discovery manager
├── config.json             # Auto-generated configuration file
├── config/                 # Configuration management
│   └── config.go          # JSON config loading and defaults
├── node/                   # P2P host management
│   └── host.go            # libp2p host creation and management
├── discovery/              # Multi-mechanism peer discovery
│   ├── manager.go         # Discovery coordinator and peer store
│   ├── mdns.go            # Enhanced mDNS integration
│   ├── mdns_discovery.go  # mDNS discovery mechanism
│   └── bootstrap.go       # Bootstrap peer discovery
├── protocol/               # Communication protocol
│   └── protocol.go        # JSON request/response over libp2p
├── services/               # Service registry and execution
│   ├── registry.go        # Dynamic service registration
│   └── examples.go        # Echo and text processing services
├── utils/                  # Utilities and helpers
│   ├── client.go          # Service client for remote calls
│   └── discovery.go       # Discovery management utilities
└── DISCOVERY.md           # Detailed discovery architecture guide
```

---

##  Getting Started

### Quick Start

```bash
# Clone the repository
git clone https://github.com/realentity/realentity-node.git
cd realentity-node

# Install dependencies
go mod tidy

# Build the project
go build

# Run a node (auto-creates config.json on first run)
go run .
```

### Configuration

The node automatically creates a `config.json` file on first run with sensible defaults:

```json
{
  "discovery": {
    "enable_mdns": true,
    "enable_bootstrap": true,
    "enable_dht": false,
    "mdns_service_tag": "realentity-mdns",
    "bootstrap_peers": [],
    "dht_rendezvous": "realentity-dht"
  },
  "log_level": "info"
}
```

**Discovery Mechanisms:**
* **mDNS**: Automatic local network discovery (enabled by default)
* **Bootstrap**: Global discovery via known peers (enabled, add peers to config)
* **DHT**: Distributed discovery (framework ready, disabled by default)

### Testing Multi-Node Communication

1. **Start first node:**

   ```bash
   go run .
   ```

2. **Start second node in another terminal:**

   ```bash
   go run .
   ```

3. **Watch automatic discovery and testing:**

   * Nodes will automatically discover each other via mDNS
   * Upon discovery, they'll automatically test each other's services
   * You'll see echo and text processing tests executed

### Adding Bootstrap Peers for Global Discovery

To connect nodes across the internet, add bootstrap peer addresses to your `config.json`:

```json
{
  "discovery": {
    "enable_bootstrap": true,
    "bootstrap_peers": [
      "/ip4/YOUR_VPS_IP/tcp/4001/p2p/12D3KooW...",
      "/ip4/ANOTHER_PEER_IP/tcp/4001/p2p/12D3KooW..."
    ]
  }
}
```

### Current Services Available

| Service | Description | Example Request |
|---------|-------------|-----------------|
| `echo` | Simple echo service | `{"message": "Hello World"}` |
| `text.process` | Text transformation | `{"text": "hello", "operation": "uppercase"}` |

---

## Current Implementation Status

###  Completed Features

- **P2P Host Creation**: Automatic Ed25519 key generation and libp2p host setup
- **Service Registry**: Dynamic service registration and execution system
- **mDNS Discovery**: Automatic peer discovery on local networks
- **Protocol Handler**: JSON-based request/response messaging
- **Example Services**: Echo and text processing services
- **Auto-Testing**: Automatic service testing when peers are discovered
- **Client Utilities**: Service calling and testing utilities

### ‍️ Current Capabilities

1. **Automatic Peer Discovery**: Nodes find each other on the same network
2. **Service Registration**: Each node exposes available services
3. **Service Execution**: Nodes can call services on remote peers
4. **Real-time Testing**: Automatic service validation between peers

---

## Upcoming Features

| Feature                        | Status        |
| ------------------------------ | ------------- |
| Modular Go structure           |  Complete    |
| JSON-based task payloads       |  Complete    |
| Peer-to-peer service execution |  Complete    |
| Service discovery protocol     | ⏳ In Progress |
| DHT-based global discovery     | ⏳ Planned     |
| Secure bootstrap discovery     | ⏳ Planned     |
| Service marketplace/registry   | ⏳ Planned     |
| Mobile-ready runtime           | ⏳ Later       |
| Metrics & debugging console    | ⏳ Optional    |

---

## Example Output

```
2025/07/12 Host created: 12D3K...abc123
2025/07/12 Listening on: /ip4/192.168.1.100/tcp/39157
2025/07/12 Node started with ID: 12D3K...abc123
2025/07/12 Services initialized: [echo text.process]
2025/07/12 Starting mDNS discovery with service tag: realentity-mdns
2025/07/12 Protocol handler registered for: /realentity/1.0.0
2025/07/12 Node is ready! Registered services: [echo text.process]
2025/07/12 Waiting for connections...
2025/07/12 Discovered peer: QmXY...def456
2025/07/12 Connected to peer: QmXY...def456
2025/07/12 Testing echo service on peer QmXY...def456
2025/07/12 Echo response: Hello from 12D3K...abc123 (from node: QmXY...def456)
```

---

## Primary Use Cases (Initial)

* Send and receive AI or image-processing tasks
* Create a scalable distributed service network
* Prepare infrastructure for future mobile clients
* Demonstrate decentralized computing capabilities

---

##  Next Steps for Development

1. **Add more service types** (image processing, AI inference)
2. **Implement service discovery protocol** (list available services)
3. **Add DHT-based discovery** for global peer finding
4. **Create configuration system** for bootstrap nodes
5. **Implement service authentication** and rate limiting
6. **Add metrics and monitoring** capabilities

---

