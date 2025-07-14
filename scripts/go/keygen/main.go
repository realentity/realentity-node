package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/realentity/realentity-node/internal/node"
)

func main() {
	var (
		generateKey = flag.Bool("generate-key", false, "Generate a new private key and peer ID")
		outputFile  = flag.String("output", "", "Output file for configuration (optional)")
	)
	flag.Parse()

	if *generateKey {
		privKey, peerID, err := node.GeneratePrivateKeyBase64()
		if err != nil {
			log.Fatalf("Failed to generate private key: %v", err)
		}

		fmt.Printf("Generated new identity:\n")
		fmt.Printf("Peer ID: %s\n", peerID.String())
		fmt.Printf("Private Key (base64): %s\n\n", privKey)

		configTemplate := fmt.Sprintf(`{
  "private_key": "%s",
  "discovery": {
    "enable_mdns": false,
    "enable_bootstrap": true,
    "enable_dht": false,
    "mdns_service_tag": "realentity-mdns",
    "mdns_quiet_mode": true,
    "bootstrap_peers": [],
    "dht_rendezvous": "realentity-dht"
  },
  "server": {
    "bind_address": "0.0.0.0",
    "port": 4001,
    "public_ip": "YOUR_PUBLIC_IP_HERE"
  },
  "log_level": "info"
}`, privKey)

		fmt.Printf("Bootstrap node configuration template:\n%s\n", configTemplate)

		if *outputFile != "" {
			err := os.WriteFile(*outputFile, []byte(configTemplate), 0644)
			if err != nil {
				log.Fatalf("Failed to write config file: %v", err)
			}
			fmt.Printf("Configuration saved to: %s\n", *outputFile)
		}

		fmt.Printf("\nBootstrap multiaddr (replace YOUR_PUBLIC_IP_HERE):\n")
		fmt.Printf("/ip4/YOUR_PUBLIC_IP_HERE/tcp/4001/p2p/%s\n", peerID.String())
		return
	}

	fmt.Println("RealEntity Node Key Generator")
	fmt.Println("Usage: go run scripts/keygen/main.go -generate-key [-output config.json]")
}
