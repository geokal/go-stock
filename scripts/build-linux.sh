#!/bin/bash

# Linux Platform Build Script
# This script must be run on a Linux environment

echo -e "======================================"
echo -e "  go-stock Linux Build Script"
echo -e "======================================"

# Check if running on Linux
if [ "$(uname -s)" != "Linux" ]; then
    echo -e "\nError: This script can only run on Linux"
    echo -e "Use Docker, WSL, or a Linux VM for building"
    echo -e "\nExample Docker build command:"
    echo -e "  docker run --rm -v \$(pwd):/app -w /app golang:1.21 ./scripts/build-linux.sh"
    exit 1
fi

# Enter project root directory
cd "$(dirname "$0")/.."

echo -e "\n[1/5] Checking dependencies..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "Error: Go compiler not found, please install Go 1.21+ first"
    exit 1
fi
echo -e "  OK Go version：$(go version)"

# Check if Node.js is installed
if ! command -v npm &> /dev/null; then
    echo -e "Error: Node.js/npm not found, please install Node.js 18+ first"
    exit 1
fi
echo -e "  OK Node.js version：$(node --version)"

# Check if Wails is installed
if ! command -v wails &> /dev/null; then
    echo -e "  Installing Wails..."
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
    if ! command -v wails &> /dev/null; then
        echo -e "Error: Wails installation failed"
        exit 1
    fi
fi
echo -e "  OK Wails version：$(wails version)"

# Check for libgtk-3-dev and libwebkit2gtk-4.0-dev (Ubuntu/Debian)
if command -v apt-get &> /dev/null; then
    if ! dpkg -l | grep -q libgtk-3-dev || ! dpkg -l | grep -q libwebkit2gtk-4.0-dev; then
        echo -e "\n  Warning: Missing required system dependencies"
        echo -e "  Please run the following command to install dependencies:"
        echo -e "    sudo apt-get update && sudo apt-get install -y \\"
        echo -e "      libgtk-3-dev \\"
        echo -e "      libwebkit2gtk-4.0-dev \\"
        echo -e "      gcc \\"
        echo -e "      libglib2.0-dev"
        exit 1
    fi
fi

echo -e "\n[2/5] Cleaning old build files..."
rm -rf build/bin/linux

echo -e "\n[3/5] Installing frontend dependencies..."
cd frontend
npm install
cd ..

echo -e "\n[4/5] Building frontend..."
cd frontend
npm run build
cd ..

echo -e "\n[5/5] Building Linux application..."
wails build --platform linux/amd64 --clean

echo -e "\n======================================"
echo -e "  Build complete!"
echo -e "======================================"
echo -e "\n  Executable location：build/bin/linux/go-stock"
echo -e "\nRun with:"
echo -e "  ./build/bin/linux/go-stock"
echo -e "\n======================================"
