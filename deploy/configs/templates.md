# Configuration Templates

## local.json - Local development
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
  "log_level": "debug",
  "server": {
    "bind_address": "0.0.0.0",
    "port": 0
  }
}

## vps-bootstrap.json - VPS Bootstrap node
{
  "discovery": {
    "enable_mdns": false,
    "enable_bootstrap": true,
    "enable_dht": true,
    "mdns_service_tag": "realentity-mdns",
    "mdns_quiet_mode": true,
    "bootstrap_peers": [],
    "dht_rendezvous": "realentity-dht"
  },
  "log_level": "info",
  "server": {
    "bind_address": "0.0.0.0",
    "port": 4001,
    "public_ip": "${PUBLIC_IP}"
  }
}

## vps-peer.json - VPS Peer node
{
  "discovery": {
    "enable_mdns": false,
    "enable_bootstrap": true,
    "enable_dht": true,
    "mdns_service_tag": "realentity-mdns",
    "mdns_quiet_mode": true,
    "bootstrap_peers": ["${BOOTSTRAP_PEER}"],
    "dht_rendezvous": "realentity-dht"
  },
  "log_level": "info",
  "server": {
    "bind_address": "0.0.0.0",
    "port": 4001
  }
}

## docker.json - Docker testing
{
  "discovery": {
    "enable_mdns": false,
    "enable_bootstrap": true,
    "enable_dht": false,
    "mdns_service_tag": "realentity-docker",
    "mdns_quiet_mode": true,
    "bootstrap_peers": ["${BOOTSTRAP_PEER}"],
    "dht_rendezvous": "realentity-dht"
  },
  "log_level": "info",
  "server": {
    "bind_address": "0.0.0.0",
    "port": 4001
  }
}
