#!/bin/bash

# Raspberry Pi Go Package Setup Script for beta-turnstiles
# This script installs Go 1.21+, Make, and downloads the specified GitHub repository

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to get system architecture
get_arch() {
    case $(uname -m) in
        armv6l) echo "armv6l" ;;
        armv7l) echo "armv6l" ;;  # Go uses armv6l for both armv6 and armv7
        aarch64) echo "arm64" ;;
        x86_64) echo "amd64" ;;
        *) echo "unknown" ;;
    esac
}

# Function to install Go
install_go() {
    local go_version="1.22.1"  # Latest stable version >= 1.21
    local arch=$(get_arch)

    if [ "$arch" = "unknown" ]; then
        print_error "Unsupported architecture: $(uname -m)"
        exit 1
    fi

    print_status "Installing Go $go_version for $arch architecture..."

    # Download Go
    local go_tarball="go${go_version}.linux-${arch}.tar.gz"
    local download_url="https://golang.org/dl/${go_tarball}"

    print_status "Downloading Go from $download_url"
    wget -O "/tmp/${go_tarball}" "$download_url"

    # Remove existing Go installation
    if [ -d "/usr/local/go" ]; then
        print_status "Removing existing Go installation..."
        sudo rm -rf /usr/local/go
    fi

    # Extract Go
    print_status "Extracting Go to /usr/local/go..."
    sudo tar -C /usr/local -xzf "/tmp/${go_tarball}"

    # Clean up
    rm "/tmp/${go_tarball}"

    print_success "Go $go_version installed successfully"
}

# Function to setup Go environment
setup_go_environment() {
    print_status "Setting up Go environment..."

    # Add Go to PATH in .bashrc if not already present
    if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        echo 'export GOPATH=$HOME/go' >> ~/.bashrc
        echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
        print_success "Go environment variables added to ~/.bashrc"
    else
        print_warning "Go environment variables already exist in ~/.bashrc"
    fi

    # Export for current session
    export PATH=$PATH:/usr/local/go/bin
    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin

    # Create GOPATH directory
    mkdir -p $HOME/go/{bin,src,pkg}
}

# Main installation process
main() {
    print_status "Starting Raspberry Pi Go setup for beta-turnstiles..."
    print_status "System: $(uname -a)"

    # Update system packages
    print_status "Updating system packages..."
    sudo apt update

    # Install essential tools
    print_status "Installing essential tools..."
    sudo apt install -y wget curl git build-essential

    # Install Make if not present
    if ! command_exists make; then
        print_status "Installing Make..."
        sudo apt install -y make
        print_success "Make installed successfully"
    else
        print_success "Make is already installed: $(make --version | head -n1)"
    fi

    # Check if Go is already installed and version
    if command_exists go; then
        local current_version=$(go version | awk '{print $3}' | sed 's/go//')
        local required_version="1.21.0"

        # Simple version comparison (works for most cases)
        if [ "$(printf '%s\n' "$required_version" "$current_version" | sort -V | head -n1)" = "$required_version" ]; then
            print_success "Go is already installed with acceptable version: $current_version"
        else
            print_warning "Go version $current_version is older than required $required_version"
            install_go
        fi
    else
        print_status "Go is not installed. Installing..."
        install_go
    fi

    # Setup Go environment
    setup_go_environment

    # Verify Go installation
    if /usr/local/go/bin/go version >/dev/null 2>&1; then
        print_success "Go installation verified: $(/usr/local/go/bin/go version)"
    else
        print_error "Go installation verification failed"
        exit 1
    fi

    # Clone the repository
    local repo_url="https://github.com/Chalk-Technologies/beta-turnstiles.git"
    local project_dir="$HOME/beta-turnstiles"

    print_status "Cloning repository from $repo_url..."

    if [ -d "$project_dir" ]; then
        print_warning "Directory $project_dir already exists. Removing it..."
        rm -rf "$project_dir"
    fi

    git clone "$repo_url" "$project_dir"
    cd "$project_dir"

    print_success "Repository cloned to $project_dir"

    # Check if go.mod exists and run go mod download
    if [ -f "go.mod" ]; then
        print_status "Found go.mod file. Downloading Go dependencies..."
        /usr/local/go/bin/go mod download
        /usr/local/go/bin/go mod tidy
        print_success "Go dependencies downloaded"
    else
        print_warning "No go.mod file found. You may need to initialize the module manually."
    fi

    # Check if Makefile exists
    if [ -f "Makefile" ]; then
        print_success "Makefile found. You can now run 'make' commands in $project_dir"
        print_status "Available make targets:"
        make help 2>/dev/null || make -n 2>/dev/null || print_warning "Could not display make targets"
    else
        print_warning "No Makefile found in the repository"
    fi

    # Final instructions
    print_success "Setup completed successfully!"
    echo
    print_status "Next steps:"
    echo "1. Source your bashrc or restart your terminal: source ~/.bashrc"
    echo "2. Navigate to the project directory: cd $project_dir"
    echo "3. Build the project (if Makefile exists): make"
    echo "4. Or build with Go directly: go build"
    echo
    print_status "Current Go version: $(/usr/local/go/bin/go version)"
    print_status "Project location: $project_dir"
}

# Run main function
main "$@"