#!/bin/bash

# Build script for Login App

set -e

echo "Building Login App..."

# Create bin directory if it doesn't exist
mkdir -p bin

# Build the application
go build -ldflags="-s -w" -o bin/login-app .

echo "Build completed successfully!"
echo "Binary created at: bin/login-app"

# Make the binary executable
chmod +x bin/login-app

echo "To run the application:"
echo "  ./bin/login-app"
echo ""
echo "Or with custom port:"
echo "  ./bin/login-app -port=3000"
