#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO_URL="https://github.com/talhahasanzia/string-resource-utility.git"
INSTALL_DIR="$HOME/.local/bin"
TEMP_DIR="/tmp/string-resource-utility-install"
EXECUTABLE_NAME="localize"

echo -e "${BLUE}🚀 String Resource Utility Installer${NC}"
echo -e "${BLUE}======================================${NC}"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Error: Go is not installed. Please install Go first: https://go.dev${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Go is installed: $(go version)${NC}"

# Check if git is installed
if ! command -v git &> /dev/null; then
    echo -e "${RED}❌ Error: Git is not installed. Please install Git first.${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Git is installed${NC}"

# Create install directory if it doesn't exist
echo -e "${YELLOW}📁 Creating installation directory: $INSTALL_DIR${NC}"
mkdir -p "$INSTALL_DIR"

# Remove temporary directory if it exists
if [ -d "$TEMP_DIR" ]; then
    echo -e "${YELLOW}🧹 Cleaning up existing temporary directory${NC}"
    rm -rf "$TEMP_DIR"
fi

# Clone the repository
echo -e "${YELLOW}📥 Cloning repository...${NC}"
git clone "$REPO_URL" "$TEMP_DIR"

# Navigate to the cloned directory
cd "$TEMP_DIR"

# Build the executable
echo -e "${YELLOW}🔨 Building executable for macOS...${NC}"
go build -o "$EXECUTABLE_NAME" localize.go

# Copy executable to install directory
echo -e "${YELLOW}📦 Installing executable to $INSTALL_DIR${NC}"
cp "$EXECUTABLE_NAME" "$INSTALL_DIR/"

# Make executable
chmod +x "$INSTALL_DIR/$EXECUTABLE_NAME"

# Clean up temporary directory
echo -e "${YELLOW}🧹 Cleaning up temporary files${NC}"
rm -rf "$TEMP_DIR"

# Add to PATH if not already there
SHELL_PROFILE=""
if [ -f "$HOME/.zshrc" ]; then
    SHELL_PROFILE="$HOME/.zshrc"
elif [ -f "$HOME/.bashrc" ]; then
    SHELL_PROFILE="$HOME/.bashrc"
elif [ -f "$HOME/.bash_profile" ]; then
    SHELL_PROFILE="$HOME/.bash_profile"
fi

if [ -n "$SHELL_PROFILE" ]; then
    if ! grep -q "export PATH=\"\$HOME/.local/bin:\$PATH\"" "$SHELL_PROFILE"; then
        echo -e "${YELLOW}⚙️  Adding $INSTALL_DIR to PATH in $SHELL_PROFILE${NC}"
        echo "" >> "$SHELL_PROFILE"
        echo "# Added by string-resource-utility installer" >> "$SHELL_PROFILE"
        echo "export PATH=\"\$HOME/.local/bin:\$PATH\"" >> "$SHELL_PROFILE"
        echo -e "${GREEN}✅ PATH updated in $SHELL_PROFILE${NC}"
    else
        echo -e "${GREEN}✅ $INSTALL_DIR is already in PATH${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  Could not detect shell profile. Please manually add the following to your shell profile:${NC}"
    echo -e "${BLUE}export PATH=\"\$HOME/.local/bin:\$PATH\"${NC}"
fi

echo -e "${GREEN}🎉 Installation completed successfully!${NC}"
echo -e "${GREEN}📍 Executable installed at: $INSTALL_DIR/$EXECUTABLE_NAME${NC}"
echo ""
echo -e "${BLUE}📖 Usage:${NC}"
echo -e "${BLUE}  To use the command immediately in this session:${NC}"
echo -e "${YELLOW}    export PATH=\"\$HOME/.local/bin:\$PATH\"${NC}"
echo -e "${YELLOW}    localize -file=/path/to/your.csv -platform=flutter${NC}"
echo ""
echo -e "${BLUE}  For future sessions, restart your terminal or run:${NC}"
if [ -n "$SHELL_PROFILE" ]; then
    echo -e "${YELLOW}    source $SHELL_PROFILE${NC}"
fi
echo -e "${YELLOW}    localize -file=absolute/path/to/your.csv -platform=flutter${NC}"
echo ""
echo -e "${GREEN}✨ Happy localizing! ✨${NC}"