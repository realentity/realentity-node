package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

// DiscoveryConfig holds configuration for discovery mechanisms
type DiscoveryConfig struct {
	EnableMDNS      bool     `json:"enable_mdns"`
	EnableBootstrap bool     `json:"enable_bootstrap"`
	EnableDHT       bool     `json:"enable_dht"`
	MDNSServiceTag  string   `json:"mdns_service_tag"`
	MDNSQuietMode   bool     `json:"mdns_quiet_mode"`
	BootstrapPeers  []string `json:"bootstrap_peers"`
	DHTRendezvous   string   `json:"dht_rendezvous"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	BindAddress string `json:"bind_address"`
	Port        int    `json:"port"`          // P2P port
	HTTPPort    int    `json:"http_port"`     // HTTP API port
	HTTPSPort   int    `json:"https_port"`    // HTTPS API port (0 = disabled)
	TLSCertFile string `json:"tls_cert_file"` // Path to TLS certificate file
	TLSKeyFile  string `json:"tls_key_file"`  // Path to TLS private key file
	PublicIP    string `json:"public_ip"`
}

// NodeConfig holds all node configuration
type NodeConfig struct {
	Discovery  DiscoveryConfig `json:"discovery"`
	Server     ServerConfig    `json:"server"`
	LogLevel   string          `json:"log_level"`
	PrivateKey string          `json:"private_key,omitempty"` // Optional base64-encoded private key for consistent peer ID
}

// DefaultConfig returns a default configuration
func DefaultConfig() *NodeConfig {
	return &NodeConfig{
		Discovery: DiscoveryConfig{
			EnableMDNS:      true,
			EnableBootstrap: true,
			EnableDHT:       false, // Disabled by default until dependencies are added
			MDNSServiceTag:  "realentity-mdns",
			MDNSQuietMode:   true, // Suppress mDNS warnings by default
			BootstrapPeers:  []string{
				// Add some default bootstrap peers here when available
				// "/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
			},
			DHTRendezvous: "realentity-dht",
		},
		Server: ServerConfig{
			BindAddress: "0.0.0.0", // Listen on all interfaces for VPS
			Port:        4001,      // Standard libp2p port
			HTTPPort:    8080,      // HTTP API port
			HTTPSPort:   0,         // HTTPS disabled by default
			TLSCertFile: "",
			TLSKeyFile:  "",
			PublicIP:    "", // Auto-detect if empty
		},
		LogLevel: "info",
	}
}

// LoadConfig loads configuration from a file
func LoadConfig(filename string) (*NodeConfig, error) {
	// If file doesn't exist, create it with default config
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		config := DefaultConfig()
		if err := SaveConfig(config, filename); err != nil {
			return nil, fmt.Errorf("failed to create default config: %v", err)
		}
		log.Printf("Created default config file: %s\n", filename)
		return config, nil
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config NodeConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

// SaveConfig saves configuration to a file
func SaveConfig(config *NodeConfig, filename string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// AddBootstrapPeer adds a bootstrap peer to the config
func (dc *DiscoveryConfig) AddBootstrapPeer(addr string) {
	for _, existing := range dc.BootstrapPeers {
		if existing == addr {
			return // Already exists
		}
	}
	dc.BootstrapPeers = append(dc.BootstrapPeers, addr)
}

// LoadConfigWithEnv loads configuration from file and applies environment overrides
func LoadConfigWithEnv(filename string) (*NodeConfig, error) {
	cfg, err := LoadConfig(filename)
	if err != nil {
		return nil, err
	}

	// Apply environment variable overrides
	applyEnvOverrides(cfg)

	// Validate configuration
	if err := ValidateConfig(cfg); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %v", err)
	}

	return cfg, nil
}

// applyEnvOverrides applies environment variable overrides to configuration
func applyEnvOverrides(cfg *NodeConfig) {
	// Server configuration
	if publicIP := os.Getenv("REALENTITY_PUBLIC_IP"); publicIP != "" {
		cfg.Server.PublicIP = publicIP
	}
	if port := os.Getenv("REALENTITY_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Server.Port = p
		}
	}
	if bindAddr := os.Getenv("REALENTITY_BIND_ADDRESS"); bindAddr != "" {
		cfg.Server.BindAddress = bindAddr
	}

	// Discovery configuration
	if enableMDNS := os.Getenv("REALENTITY_ENABLE_MDNS"); enableMDNS != "" {
		cfg.Discovery.EnableMDNS = strings.ToLower(enableMDNS) == "true"
	}
	if enableBootstrap := os.Getenv("REALENTITY_ENABLE_BOOTSTRAP"); enableBootstrap != "" {
		cfg.Discovery.EnableBootstrap = strings.ToLower(enableBootstrap) == "true"
	}
	if enableDHT := os.Getenv("REALENTITY_ENABLE_DHT"); enableDHT != "" {
		cfg.Discovery.EnableDHT = strings.ToLower(enableDHT) == "true"
	}
	if bootstrapPeers := os.Getenv("REALENTITY_BOOTSTRAP_PEERS"); bootstrapPeers != "" {
		cfg.Discovery.BootstrapPeers = strings.Split(bootstrapPeers, ",")
	}
	if logLevel := os.Getenv("REALENTITY_LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}
}

// ValidateConfig validates the configuration for common issues
func ValidateConfig(cfg *NodeConfig) error {
	// Validate server configuration
	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", cfg.Server.Port)
	}

	// Validate discovery configuration
	if !cfg.Discovery.EnableMDNS && !cfg.Discovery.EnableBootstrap && !cfg.Discovery.EnableDHT {
		return fmt.Errorf("at least one discovery mechanism must be enabled")
	}

	// If bootstrap is enabled for a peer, ensure bootstrap peers are configured
	if cfg.Discovery.EnableBootstrap && len(cfg.Discovery.BootstrapPeers) == 0 && cfg.Server.PublicIP != "" {
		log.Printf("Warning: Bootstrap discovery enabled but no bootstrap peers configured. This might be a bootstrap node.")
	}

	// Validate bootstrap peer addresses
	for _, peer := range cfg.Discovery.BootstrapPeers {
		if !strings.HasPrefix(peer, "/ip4/") && !strings.HasPrefix(peer, "/ip6/") {
			return fmt.Errorf("invalid bootstrap peer address format: %s", peer)
		}
	}

	return nil
}

// GetConfigForDeployment returns a configuration template for a deployment type
func GetConfigForDeployment(deploymentType string, options map[string]string) (*NodeConfig, error) {
	var cfg *NodeConfig

	switch deploymentType {
	case "local":
		cfg = &NodeConfig{
			Discovery: DiscoveryConfig{
				EnableMDNS:      true,
				EnableBootstrap: false,
				EnableDHT:       false,
				MDNSServiceTag:  "realentity-local",
				MDNSQuietMode:   false,
				BootstrapPeers:  []string{},
				DHTRendezvous:   "realentity-dht",
			},
			Server: ServerConfig{
				BindAddress: "0.0.0.0",
				Port:        0, // Random port for local
			},
			LogLevel: "debug",
		}
	case "vps-bootstrap":
		cfg = &NodeConfig{
			Discovery: DiscoveryConfig{
				EnableMDNS:      false,
				EnableBootstrap: true,
				EnableDHT:       true,
				MDNSServiceTag:  "realentity-mdns",
				MDNSQuietMode:   true,
				BootstrapPeers:  []string{},
				DHTRendezvous:   "realentity-dht",
			},
			Server: ServerConfig{
				BindAddress: "0.0.0.0",
				Port:        4001,
				PublicIP:    options["public_ip"],
			},
			LogLevel: "info",
		}
	case "vps-peer":
		bootstrapPeers := []string{}
		if peer := options["bootstrap_peer"]; peer != "" {
			bootstrapPeers = []string{peer}
		}
		cfg = &NodeConfig{
			Discovery: DiscoveryConfig{
				EnableMDNS:      false,
				EnableBootstrap: true,
				EnableDHT:       true,
				MDNSServiceTag:  "realentity-mdns",
				MDNSQuietMode:   true,
				BootstrapPeers:  bootstrapPeers,
				DHTRendezvous:   "realentity-dht",
			},
			Server: ServerConfig{
				BindAddress: "0.0.0.0",
				Port:        4001,
			},
			LogLevel: "info",
		}
	case "docker":
		bootstrapPeers := []string{}
		if peer := options["bootstrap_peer"]; peer != "" {
			bootstrapPeers = []string{peer}
		}
		cfg = &NodeConfig{
			Discovery: DiscoveryConfig{
				EnableMDNS:      false,
				EnableBootstrap: true,
				EnableDHT:       false,
				MDNSServiceTag:  "realentity-docker",
				MDNSQuietMode:   true,
				BootstrapPeers:  bootstrapPeers,
				DHTRendezvous:   "realentity-dht",
			},
			Server: ServerConfig{
				BindAddress: "0.0.0.0",
				Port:        4001,
			},
			LogLevel: "info",
		}
	default:
		return nil, fmt.Errorf("unknown deployment type: %s", deploymentType)
	}

	return cfg, nil
}
