#!/usr/bin/env python3
"""
scripts/utils/config-validator.py - Configuration file validator

This script validates RealEntity Node configuration files for common issues
and provides suggestions for improvements.

Usage:
    python3 scripts/utils/config-validator.py config.json
    python3 scripts/utils/config-validator.py --help

Requirements:
    - Python 3.6+
    - No external dependencies (uses only standard library)
"""

import json
import sys
import argparse
from pathlib import Path


def validate_ports(config):
    """Validate port configurations"""
    issues = []
    server = config.get('server', {})
    
    p2p_port = server.get('port', 0)
    http_port = server.get('http_port', 0)
    https_port = server.get('https_port', 0)
    
    if not (1 <= p2p_port <= 65535):
        issues.append(f"Invalid P2P port: {p2p_port} (must be 1-65535)")
    
    if http_port and not (1 <= http_port <= 65535):
        issues.append(f"Invalid HTTP port: {http_port} (must be 1-65535)")
    
    if https_port and not (1 <= https_port <= 65535):
        issues.append(f"Invalid HTTPS port: {https_port} (must be 1-65535)")
    
    # Check for port conflicts
    ports = [p for p in [p2p_port, http_port, https_port] if p > 0]
    if len(ports) != len(set(ports)):
        issues.append("Port conflict detected: multiple services using same port")
    
    return issues


def validate_discovery(config):
    """Validate discovery configurations"""
    issues = []
    discovery = config.get('discovery', {})
    
    enable_mdns = discovery.get('enable_mdns', False)
    enable_bootstrap = discovery.get('enable_bootstrap', False)
    enable_dht = discovery.get('enable_dht', False)
    
    if not (enable_mdns or enable_bootstrap or enable_dht):
        issues.append("No discovery mechanisms enabled - node will be isolated")
    
    if enable_bootstrap:
        bootstrap_peers = discovery.get('bootstrap_peers', [])
        if not bootstrap_peers:
            issues.append("Bootstrap discovery enabled but no bootstrap peers configured")
        
        for peer in bootstrap_peers:
            if not peer.startswith(('/ip4/', '/ip6/')):
                issues.append(f"Invalid bootstrap peer format: {peer}")
    
    return issues


def validate_tls(config):
    """Validate TLS configurations"""
    issues = []
    server = config.get('server', {})
    
    https_port = server.get('https_port', 0)
    cert_file = server.get('tls_cert_file', '')
    key_file = server.get('tls_key_file', '')
    
    if https_port > 0:
        if not cert_file:
            issues.append("HTTPS enabled but no certificate file specified")
        if not key_file:
            issues.append("HTTPS enabled but no private key file specified")
        
        if cert_file and not Path(cert_file).is_file():
            issues.append(f"Certificate file not found: {cert_file}")
        if key_file and not Path(key_file).is_file():
            issues.append(f"Private key file not found: {key_file}")
    
    return issues


def validate_config(config_path):
    """Main validation function"""
    try:
        with open(config_path, 'r') as f:
            config = json.load(f)
    except FileNotFoundError:
        return [f"Configuration file not found: {config_path}"]
    except json.JSONDecodeError as e:
        return [f"Invalid JSON format: {e}"]
    
    all_issues = []
    all_issues.extend(validate_ports(config))
    all_issues.extend(validate_discovery(config))
    all_issues.extend(validate_tls(config))
    
    return all_issues


def main():
    parser = argparse.ArgumentParser(
        description="Validate RealEntity Node configuration files",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
    python3 config-validator.py config.json
    python3 config-validator.py production-config.json
        """
    )
    parser.add_argument('config_file', help='Path to configuration file')
    parser.add_argument('--quiet', '-q', action='store_true', 
                       help='Only show errors, not success messages')
    
    args = parser.parse_args()
    
    issues = validate_config(args.config_file)
    
    if issues:
        print(f" Configuration validation failed for {args.config_file}")
        print()
        for i, issue in enumerate(issues, 1):
            print(f"{i}. {issue}")
        sys.exit(1)
    else:
        if not args.quiet:
            print(f" Configuration {args.config_file} is valid")
        sys.exit(0)


if __name__ == '__main__':
    main()
