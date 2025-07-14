#!/bin/bash
# scripts/generate-tls-cert.sh - Generate self-signed TLS certificates for HTTPS

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
CERTS_DIR="$PROJECT_ROOT/certs"

# Create certs directory if it doesn't exist
mkdir -p "$CERTS_DIR"

# Default values
DOMAIN="localhost"
DAYS=365
COUNTRY="US"
STATE="CA"
CITY="San Francisco"
ORG="RealEntity"
OU="Development"

print_usage() {
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  -d, --domain DOMAIN     Domain name (default: localhost)"
    echo "  --days DAYS            Certificate validity in days (default: 365)"
    echo "  --country COUNTRY      Country code (default: US)"
    echo "  --state STATE          State (default: CA)"
    echo "  --city CITY            City (default: San Francisco)"
    echo "  --org ORG              Organization (default: RealEntity)"
    echo "  --ou OU                Organizational Unit (default: Development)"
    echo "  -h, --help             Show this help"
    echo ""
    echo "Example:"
    echo "  $0 --domain yourdomain.com --days 3650"
    echo ""
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--domain)
            DOMAIN="$2"
            shift 2
            ;;
        --days)
            DAYS="$2"
            shift 2
            ;;
        --country)
            COUNTRY="$2"
            shift 2
            ;;
        --state)
            STATE="$2"
            shift 2
            ;;
        --city)
            CITY="$2"
            shift 2
            ;;
        --org)
            ORG="$2"
            shift 2
            ;;
        --ou)
            OU="$2"
            shift 2
            ;;
        -h|--help)
            print_usage
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            print_usage
            exit 1
            ;;
    esac
done

echo "Generating TLS certificate for domain: $DOMAIN"
echo "Certificate will be valid for $DAYS days"
echo ""

# Generate private key
openssl genrsa -out "$CERTS_DIR/server.key" 2048

# Generate certificate signing request
openssl req -new -key "$CERTS_DIR/server.key" -out "$CERTS_DIR/server.csr" -subj "/C=$COUNTRY/ST=$STATE/L=$CITY/O=$ORG/OU=$OU/CN=$DOMAIN"

# Generate self-signed certificate
openssl x509 -req -in "$CERTS_DIR/server.csr" -signkey "$CERTS_DIR/server.key" -out "$CERTS_DIR/server.crt" -days $DAYS -extensions v3_req -extfile <(cat <<EOF
[v3_req]
keyUsage = keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1 = $DOMAIN
DNS.2 = localhost
DNS.3 = *.localhost
IP.1 = 127.0.0.1
IP.2 = ::1
EOF
)

# Clean up CSR file
rm "$CERTS_DIR/server.csr"

echo " TLS certificate generated successfully!"
echo ""
echo "Files created:"
echo "  Certificate: $CERTS_DIR/server.crt"
echo "  Private Key: $CERTS_DIR/server.key"
echo ""
echo "To use HTTPS, update your config.json:"
echo "{"
echo "  \"server\": {"
echo "    \"https_port\": 8443,"
echo "    \"tls_cert_file\": \"certs/server.crt\","
echo "    \"tls_key_file\": \"certs/server.key\""
echo "  }"
echo "}"
echo ""
echo "️  Note: This is a self-signed certificate for development only."
echo "   For production, use certificates from a trusted CA like Let's Encrypt."
