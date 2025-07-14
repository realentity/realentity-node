#  Universal VPS Deployment Guide

Deploy RealEntity Nodes on **any VPS** or **Ubuntu server** worldwide with automatic platform detection and setup.

##  Quick Start (One Command)

### Deploy Bootstrap Node
```bash
curl -fsSL https://raw.githubusercontent.com/your-repo/realentity-node/main/quick_deploy.sh | bash -s bootstrap
```

### Deploy Peer Node
```bash
curl -fsSL https://raw.githubusercontent.com/your-repo/realentity-node/main/quick_deploy.sh | bash -s peer "/ip4/YOUR_BOOTSTRAP_IP/tcp/4001/p2p/PEER_ID"
```

##  Supported Platforms

 **Linux Distributions:**
- Ubuntu (18.04, 20.04, 22.04, 24.04)
- Debian (10, 11, 12)
- CentOS (7, 8, 9)
- RHEL (7, 8, 9)
- Fedora (35+)
- Alpine Linux

 **Architectures:**
- x86_64 (AMD64)
- ARM64 (AArch64)
- ARMv7

 **Cloud Providers:**
- AWS EC2
- Google Cloud Platform
- Azure Virtual Machines
- DigitalOcean Droplets
- Linode
- Vultr
- Hetzner
- OVH
- Any VPS provider

## ️ Manual Deployment

### Step 1: Download Deployment Scripts

```bash
# Download all scripts
wget https://raw.githubusercontent.com/your-repo/realentity-node/main/deploy_bootstrap.sh
wget https://raw.githubusercontent.com/your-repo/realentity-node/main/deploy_peer.sh
chmod +x deploy_bootstrap.sh deploy_peer.sh
```

### Step 2: Deploy Bootstrap Node (First VPS)

```bash
./deploy_bootstrap.sh
```

**What it does:**
-  Detects platform and architecture
-  Installs Go if not present
-  Creates optimized VPS configuration
-  Sets up systemd service
-  Configures firewall rules
-  Provides multiaddr for peer connections

### Step 3: Deploy Peer Nodes (Additional VPS)

```bash
./deploy_peer.sh "/ip4/BOOTSTRAP_IP/tcp/4001/p2p/BOOTSTRAP_PEER_ID"
```

**Example:**
```bash
./deploy_peer.sh "/ip4/198.51.100.1/tcp/4001/p2p/12D3KooWExample123..."
```

##  VPS Configuration Requirements

### Minimum System Requirements
- **RAM:** 512 MB (1 GB recommended)
- **CPU:** 1 vCPU
- **Storage:** 2 GB free space
- **Network:** Public IP address

### Required Ports
- **4001/tcp:** P2P communication
- **8080/tcp:** Health check endpoint

### Firewall Configuration
Scripts automatically configure:
- **UFW** (Ubuntu/Debian)
- **firewalld** (CentOS/RHEL/Fedora)
- **iptables** (fallback)

Manual configuration:
```bash
# Ubuntu/Debian
sudo ufw allow 4001/tcp
sudo ufw allow 8080/tcp

# CentOS/RHEL/Fedora
sudo firewall-cmd --permanent --add-port=4001/tcp
sudo firewall-cmd --permanent --add-port=8080/tcp
sudo firewall-cmd --reload
```

##  Service Management

### Start/Stop Nodes
```bash
# Manual start
./start.sh

# Manual stop
./stop.sh

# As systemd service
sudo systemctl start realentity-node
sudo systemctl stop realentity-node
sudo systemctl restart realentity-node
```

### Monitor Logs
```bash
# Service logs
sudo journalctl -f -u realentity-node

# Manual logs
tail -f ~/realentity-node/logs.txt
```

### Health Checks
```bash
# Local health check
curl http://localhost:8080/health

# Remote health check
curl http://YOUR_VPS_IP:8080/health
```

##  Cloud Provider Specific Instructions

### AWS EC2
```bash
# Security Group: Open ports 4001, 8080
# Instance type: t3.micro or larger
# AMI: Ubuntu 22.04 LTS

# Connect and deploy
ssh -i your-key.pem ubuntu@ec2-ip-address.compute.amazonaws.com
curl -fsSL https://your-repo/quick_deploy.sh | bash -s bootstrap
```

### Google Cloud Platform
```bash
# Firewall rules
gcloud compute firewall-rules create realentity-p2p --allow tcp:4001,tcp:8080

# Deploy on instance
gcloud compute ssh your-instance
curl -fsSL https://your-repo/quick_deploy.sh | bash -s bootstrap
```

### DigitalOcean
```bash
# Create droplet with Ubuntu 22.04
# Add firewall rules for ports 4001, 8080

ssh root@your-droplet-ip
curl -fsSL https://your-repo/quick_deploy.sh | bash -s bootstrap
```

### Azure
```bash
# Network Security Group: Allow ports 4001, 8080
# VM size: Standard_B1s or larger

ssh azureuser@your-vm-ip
curl -fsSL https://your-repo/quick_deploy.sh | bash -s bootstrap
```

##  Multi-Node Network Setup

### 3-Node Network Example

**Node 1 (Bootstrap) - US East:**
```bash
# VPS: 198.51.100.1
./deploy_bootstrap.sh
./start.sh
# Note the multiaddr from logs
```

**Node 2 (Peer) - Europe:**
```bash
# VPS: 203.0.113.5
./deploy_peer.sh "/ip4/198.51.100.1/tcp/4001/p2p/12D3KooW..."
./start.sh
```

**Node 3 (Peer) - Asia:**
```bash
# VPS: 192.0.2.10
./deploy_peer.sh "/ip4/198.51.100.1/tcp/4001/p2p/12D3KooW..."
./start.sh
```

### Network Topology
```
    Bootstrap (US)
   /              \
Peer (EU) ←→ Peer (Asia)
```

All nodes automatically discover each other and form a mesh network.

##  Troubleshooting

### Common Issues

**Connection Refused:**
- Check firewall rules
- Verify public IP accessibility
- Ensure correct ports are open

**Go Installation Failed:**
- Check internet connectivity
- Verify architecture support
- Try manual Go installation

**Permission Denied:**
- Ensure user has sudo access
- Check file permissions
- Run scripts with proper user

**Bootstrap Connection Failed:**
- Verify bootstrap node is running
- Check multiaddr format
- Confirm network connectivity

### Debug Commands
```bash
# Check node status
curl http://localhost:8080/health

# Check listening ports
netstat -tlnp | grep :4001

# Check firewall status
sudo ufw status  # Ubuntu/Debian
sudo firewall-cmd --list-all  # CentOS/RHEL

# Check service status
sudo systemctl status realentity-node

# View full logs
sudo journalctl -u realentity-node --no-pager
```

### Log Analysis
```bash
# Connection issues
grep "Failed to connect" logs.txt

# Discovery issues
grep "Discovery" logs.txt

# Service issues
grep "Service" logs.txt
```

##  Scaling Considerations

### Load Balancing
- Use multiple bootstrap nodes for redundancy
- Distribute peers across regions
- Implement health checks

### Performance Optimization
- Use SSD storage for better I/O
- Increase RAM for large networks
- Use dedicated CPU instances

### Security
- Configure VPN for sensitive data
- Use private networks where possible
- Implement authentication layers

##  Production Deployment

### Recommended Setup
1. **3 Bootstrap nodes** across different regions
2. **Multiple peer nodes** in each region
3. **Load balancer** for HTTP endpoints
4. **Monitoring** with Prometheus/Grafana
5. **Alerting** for node failures

### Monitoring Stack
```bash
# Deploy monitoring (optional)
curl -fsSL https://your-repo/deploy_monitoring.sh | bash
```

This creates a comprehensive deployment system that works on any VPS or Ubuntu server worldwide! 
