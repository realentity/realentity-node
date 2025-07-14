#!/bin/bash
# migrate-deployment.sh - Migrate from legacy deployment scripts to unified system

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Legacy scripts to be removed
LEGACY_SCRIPTS=(
    "deploy_bootstrap.sh"
    "deploy_bootstrap.bat"
    "deploy_bootstrap_clean.sh"
    "deploy_peer.sh"
    "deploy_peer.bat" 
    "deploy_peer_clean.sh"
    "quick_deploy.sh"
    "quick_deploy_clean.sh"
)

# Legacy documentation to be archived
LEGACY_DOCS=(
    "VPS_DEPLOYMENT.md"
    "UNIVERSAL_DEPLOYMENT.md"
    "DEPLOYMENT_SUMMARY.md"
)

# Create backup directory
create_backup() {
    local backup_dir="deployment-backup-$(date +%Y%m%d-%H%M%S)"
    mkdir -p "$backup_dir"
    echo "$backup_dir"
}

# Show migration mapping
show_migration_guide() {
    echo "=== Deployment Script Migration Guide ==="
    echo ""
    echo "OLD COMMAND                              → NEW COMMAND"
    echo "────────────────────────────────────────────────────────────────────"
    echo "./deploy_bootstrap.sh                   → ./deploy/universal.sh vps-bootstrap --public-ip <IP>"
    echo "./deploy_peer.sh <bootstrap-addr>       → ./deploy/universal.sh vps-peer --bootstrap-peer <addr>"
    echo "./quick_deploy.sh bootstrap             → ./deploy/universal.sh vps-bootstrap --public-ip <IP>"
    echo "./quick_deploy.sh peer <addr>           → ./deploy/universal.sh vps-peer --bootstrap-peer <addr>"
    echo "./docker-test.sh                        → ./deploy/universal.sh docker"
    echo "go run main.go (multiple instances)     → ./deploy/universal.sh local"
    echo ""
    echo "=== New Features Available ==="
    echo "• Dry-run mode: --dry-run"
    echo "• Clean deployment: --clean"
    echo "• Custom config: --config <file>"
    echo "• Kubernetes deployment: ./deploy/universal.sh k8s"
    echo "• Build automation: make deploy-bootstrap, make deploy-peer"
    echo ""
}

# Check if legacy scripts are being used
check_usage() {
    print_status "Checking for legacy script usage..."
    
    local found_usage=false
    
    # Check for recent executions in shell history
    if command -v history &> /dev/null; then
        for script in "${LEGACY_SCRIPTS[@]}"; do
            if history | grep -q "$script" 2>/dev/null; then
                print_warning "Found recent usage of $script in shell history"
                found_usage=true
            fi
        done
    fi
    
    # Check for references in other scripts
    for script in *.sh; do
        if [ -f "$script" ] && [ "$script" != "migrate-deployment.sh" ]; then
            for legacy in "${LEGACY_SCRIPTS[@]}"; do
                if grep -q "$legacy" "$script" 2>/dev/null; then
                    print_warning "Found reference to $legacy in $script"
                    found_usage=true
                fi
            done
        fi
    done
    
    if [ "$found_usage" = false ]; then
        print_status "No active usage of legacy scripts detected"
    fi
}

# Create compatibility symlinks
create_compatibility_layer() {
    print_status "Creating compatibility symlinks..."
    
    # Create wrapper for common use cases
    cat > deploy_bootstrap_compat.sh << 'EOF'
#!/bin/bash
echo "️  This script has been replaced by the unified deployment system"
echo "New command: ./deploy/universal.sh vps-bootstrap --public-ip <your-ip>"
echo ""
read -p "Do you want to run the new command? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    read -p "Enter your public IP: " public_ip
    exec ./deploy/universal.sh vps-bootstrap --public-ip "$public_ip"
fi
EOF
    chmod +x deploy_bootstrap_compat.sh
    
    cat > deploy_peer_compat.sh << 'EOF'
#!/bin/bash
echo "️  This script has been replaced by the unified deployment system"
echo "New command: ./deploy/universal.sh vps-peer --bootstrap-peer <bootstrap-address>"
echo ""
if [ -n "$1" ]; then
    echo "Detected bootstrap peer: $1"
    read -p "Do you want to run the new command with this peer? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        exec ./deploy/universal.sh vps-peer --bootstrap-peer "$1"
    fi
else
    read -p "Do you want to run the new command? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        read -p "Enter bootstrap peer address: " bootstrap_peer
        exec ./deploy/universal.sh vps-peer --bootstrap-peer "$bootstrap_peer"
    fi
fi
EOF
    chmod +x deploy_peer_compat.sh
    
    print_status "Compatibility wrappers created: deploy_bootstrap_compat.sh, deploy_peer_compat.sh"
}

# Archive legacy files
archive_legacy() {
    local backup_dir="$1"
    local dry_run="$2"
    
    print_status "Archiving legacy deployment files..."
    
    for script in "${LEGACY_SCRIPTS[@]}"; do
        if [ -f "$script" ]; then
            if [ "$dry_run" = "true" ]; then
                echo "Would move $script to $backup_dir/"
            else
                mv "$script" "$backup_dir/"
                print_status "Archived: $script"
            fi
        fi
    done
    
    for doc in "${LEGACY_DOCS[@]}"; do
        if [ -f "$doc" ]; then
            if [ "$dry_run" = "true" ]; then
                echo "Would move $doc to $backup_dir/"
            else
                mv "$doc" "$backup_dir/"
                print_status "Archived: $doc"
            fi
        fi
    done
    
    # Archive old configs that are no longer used
    for config in config-*.json; do
        if [ -f "$config" ] && [ "$config" != "config-local.json" ]; then
            if [ "$dry_run" = "true" ]; then
                echo "Would move $config to $backup_dir/"
            else
                mv "$config" "$backup_dir/"
                print_status "Archived: $config"
            fi
        fi
    done
}

# Update documentation references
update_docs() {
    print_status "Updating documentation references..."
    
    # Update README.md if it exists
    if [ -f "README.md" ]; then
        sed -i.bak 's|deploy_bootstrap\.sh|deploy/universal.sh vps-bootstrap|g' README.md
        sed -i.bak 's|deploy_peer\.sh|deploy/universal.sh vps-peer|g' README.md
        sed -i.bak 's|quick_deploy\.sh|deploy/universal.sh|g' README.md
        print_status "Updated README.md (backup saved as README.md.bak)"
    fi
}

# Main migration function
main() {
    case "${1:-interactive}" in
        "check")
            show_migration_guide
            check_usage
            ;;
        "migrate")
            backup_dir=$(create_backup)
            print_status "Created backup directory: $backup_dir"
            
            check_usage
            archive_legacy "$backup_dir" "false"
            create_compatibility_layer
            update_docs
            
            print_status "Migration completed!"
            print_status "Legacy files backed up to: $backup_dir"
            print_status "Use 'deploy/universal.sh --help' to see new deployment options"
            ;;
        "dry-run")
            backup_dir="backup-preview"
            print_status "DRY RUN - No files will be moved"
            
            show_migration_guide
            check_usage
            archive_legacy "$backup_dir" "true"
            ;;
        "interactive")
            show_migration_guide
            echo ""
            echo "Migration Options:"
            echo "1. Check current setup (safe)"
            echo "2. Dry run (preview changes)"
            echo "3. Full migration (backup + cleanup)"
            echo "4. Exit"
            echo ""
            read -p "Choose option (1-4): " -n 1 -r
            echo ""
            
            case $REPLY in
                1) main "check" ;;
                2) main "dry-run" ;;
                3) main "migrate" ;;
                4) print_status "Exiting without changes" ;;
                *) print_error "Invalid option" ;;
            esac
            ;;
        *)
            echo "Usage: $0 [check|migrate|dry-run|interactive]"
            echo "  check       - Show migration guide and check for legacy usage"
            echo "  migrate     - Perform full migration (backup legacy files)"
            echo "  dry-run     - Preview what would be changed"
            echo "  interactive - Interactive migration wizard (default)"
            ;;
    esac
}

main "$@"
