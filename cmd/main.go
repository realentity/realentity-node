package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/realentity/realentity-node/internal/api"
	"github.com/realentity/realentity-node/internal/config"
	"github.com/realentity/realentity-node/internal/discovery"
	"github.com/realentity/realentity-node/internal/node"
	"github.com/realentity/realentity-node/internal/protocol"
	"github.com/realentity/realentity-node/internal/services"
	"github.com/realentity/realentity-node/internal/utils"

	// Import service implementations to trigger their init() functions
	_ "github.com/realentity/realentity-node/internal/services/impl"
)

func main() {
	ctx := context.Background()

	// Parse command line flags
	configFile := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create the libp2p host with VPS optimization
	var host host.Host

	if cfg.PrivateKey != "" {
		// Use hardcoded private key for consistent peer ID
		log.Println("Using hardcoded private key for consistent peer ID")
		priv, err := node.DecodePrivateKeyFromBase64(cfg.PrivateKey)
		if err != nil {
			log.Fatalf("Failed to decode private key: %v", err)
		}

		hostConfig := node.DefaultHostConfig()
		if cfg.Server.PublicIP != "" {
			hostConfig.ExternalIP = cfg.Server.PublicIP
		}
		if cfg.Server.Port != 0 {
			hostConfig.ListenPort = cfg.Server.Port
		}

		host, err = node.CreateHostWithPrivateKey(ctx, hostConfig, priv)
		if err != nil {
			log.Fatalf("Failed to create host with private key: %v", err)
		}
	} else if cfg.Server.PublicIP != "" {
		// VPS mode with known public IP
		host, err = node.CreateVPSHost(ctx, cfg.Server.PublicIP, cfg.Server.Port)
	} else {
		// Default mode
		host, err = node.CreateHost(ctx)
	}
	if err != nil {
		log.Fatalln("Host creation failed:", err)
	}

	log.Println("Node started with ID:", utils.FormatPeerID(host.ID()))

	// Initialize services
	initializeServices(host.ID().String())

	// Set up enhanced discovery
	dm := discovery.NewDiscoveryManager(host)

	// Add discovery mechanisms based on config
	if cfg.Discovery.EnableMDNS {
		err = discovery.SetupEnhancedMDNS(ctx, host, cfg.Discovery.MDNSServiceTag, dm)
		if err != nil {
			log.Printf("mDNS setup failed: %v\n", err)
		}
	}

	if cfg.Discovery.EnableBootstrap && len(cfg.Discovery.BootstrapPeers) > 0 {
		bootstrapDisc, err := discovery.NewBootstrapDiscovery(host, cfg.Discovery.BootstrapPeers)
		if err != nil {
			log.Printf("Bootstrap discovery setup failed: %v\n", err)
		} else {
			dm.AddMechanism(bootstrapDisc)
		}
	}

	// Start discovery
	if err := dm.Start(); err != nil {
		log.Printf("Failed to start discovery manager: %v\n", err)
	}

	// Register protocol handler
	protocol.RegisterHandler(host, "/realentity/1.0.0")

	// Start HTTP API server
	log.Printf("Starting HTTP API server on port %d\n", cfg.Server.HTTPPort)
	if cfg.Server.HTTPSPort > 0 && cfg.Server.TLSCertFile != "" && cfg.Server.TLSKeyFile != "" {
		log.Printf("HTTPS will be available on port %d\n", cfg.Server.HTTPSPort)
	}
	apiServer := api.NewServer(host, dm, cfg.Server.HTTPPort, cfg.Server.HTTPSPort, cfg.Server.TLSCertFile, cfg.Server.TLSKeyFile)
	go func() {
		if err := apiServer.Start(); err != nil {
			log.Printf("HTTP API server failed: %v\n", err)
		}
	}()

	log.Printf("Node is ready! Registered services: %v\n", services.GlobalRegistry.ListServices())
	log.Println("Discovery mechanisms active:")
	log.Printf("- mDNS: %v\n", cfg.Discovery.EnableMDNS)
	log.Printf("- Bootstrap: %v (%d peers)\n", cfg.Discovery.EnableBootstrap, len(cfg.Discovery.BootstrapPeers))
	log.Printf("- DHT: %v\n", cfg.Discovery.EnableDHT)
	log.Printf("- HTTP API: http://localhost:%d/health\n", cfg.Server.HTTPPort)
	if cfg.Server.HTTPSPort > 0 && cfg.Server.TLSCertFile != "" && cfg.Server.TLSKeyFile != "" {
		log.Printf("- HTTPS API: https://localhost:%d/health\n", cfg.Server.HTTPSPort)
	}
	log.Println("Waiting for connections...")
	log.Printf("Start another instance to see automatic peer discovery and service testing!\n")

	// Info about mDNS warnings on Windows
	if cfg.Discovery.EnableMDNS {
		log.Printf("Note: mDNS multicast warnings are normal on Windows and don't affect functionality\n")
	}

	// Periodically log discovery stats
	go logDiscoveryStats(dm)

	select {} // keep alive
}

func initializeServices(nodeID string) {
	// Use dynamic service loading - automatically discover and register all services
	log.Printf("Discovering and initializing services for node %s", nodeID)

	// Get all available services from the registry
	availableServices := services.GlobalServiceRegistry.ListAvailableServices()
	log.Printf("Available services: %v", availableServices)

	// Create and register all services dynamically
	serviceInstances, err := services.GlobalServiceRegistry.CreateAllServices(nodeID)
	if err != nil {
		log.Printf("Failed to create services: %v", err)
		return
	}

	// Register all created services with the global registry
	registeredCount := 0
	for _, service := range serviceInstances {
		if err := services.GlobalRegistry.RegisterService(service); err != nil {
			log.Printf("Failed to register service %s: %v", service.Name, err)
		} else {
			registeredCount++
			log.Printf("Successfully registered service: %s (v%s) - %s",
				service.Name, service.Version, service.Description)
		}
	}

	log.Printf("Service initialization complete: %d/%d services registered successfully",
		registeredCount, len(serviceInstances))
	log.Printf("Active services: %v", services.GlobalRegistry.ListServices())
}

func logDiscoveryStats(dm *discovery.DiscoveryManager) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		peers := dm.GetPeers()
		connectable := dm.GetConnectablePeers()
		log.Printf("Discovery stats: %d total peers, %d connectable\n", len(peers), len(connectable))
	}
}
