#!/bin/bash
# scripts/run.sh - Script runner utility for RealEntity Node

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

print_usage() {
    echo "RealEntity Node Script Runner"
    echo ""
    echo "Usage: $0 <category>/<script> [args...]"
    echo ""
    echo "Available scripts:"
    echo ""
    echo "Go Scripts (scripts/go/):"
    echo "  keygen                 - Generate private keys and peer IDs"
    echo ""
    echo "Shell Scripts (scripts/shell/):"
    echo "  generate-tls-cert      - Generate TLS certificates"
    echo ""
    echo "Examples:"
    echo "  $0 keygen"
    echo "  $0 keygen -output config.json"
    echo "  $0 generate-tls-cert --domain localhost"
    echo ""
    echo "For script-specific help, run:"
    echo "  $0 <script> --help"
}

if [ $# -eq 0 ]; then
    print_usage
    exit 1
fi

SCRIPT_NAME="$1"
shift

case "$SCRIPT_NAME" in
    "keygen")
        cd "$SCRIPT_DIR/go/keygen"
        exec go run main.go "$@"
        ;;
    "generate-tls-cert")
        cd "$SCRIPT_DIR/shell"
        if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
            exec ./generate-tls-cert.bat "$@"
        else
            exec ./generate-tls-cert.sh "$@"
        fi
        ;;
    "--help"|"-h")
        print_usage
        exit 0
        ;;
    *)
        echo "Unknown script: $SCRIPT_NAME"
        echo ""
        print_usage
        exit 1
        ;;
esac
