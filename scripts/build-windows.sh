#!/bin/bash

# Windows Platform Build Script

echo -e "======================================"
echo -e "  go-stock Windows Build Script"
echo -e "======================================"

# Enter project root directory
cd "$(dirname "$0")/.."

echo -e "\n[1/4] Checking dependencies..."

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

echo -e "\n[2/4] Cleaning old build files..."
rm -rf build/bin/windows

echo -e "\n[3/4] Installing frontend dependencies..."
cd frontend
npm install
cd ..

echo -e "\n[4/4] Building Windows application..."
wails build --clean --platform windows/amd64

echo -e "\n======================================"
echo -e "  Build complete!"
echo -e "======================================"
echo -e "\n  Executable location：build/bin/windows/go-stock.exe"
echo -e "\nRun with:"
echo -e "  ./build/bin/windows/go-stock.exe"
echo -e "\n======================================"
