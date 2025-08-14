#!/bin/bash

# Run script for Login App in development mode

set -e

echo "Starting Login App in development mode..."

# Set development environment variables
export LOG_LEVEL=debug
export ENVIRONMENT=development
export JWT_SECRET=development-secret-key-change-in-production

# Run the application with hot reload support
echo "Running on http://localhost:8080"
echo "Press Ctrl+C to stop"

go run . -env=development
