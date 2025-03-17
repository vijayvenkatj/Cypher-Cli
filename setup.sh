#!/bin/bash

# Clone the repository
git clone https://github.com/vijayvenkatj/Cypher-Cli.git
cd Cypher-Cli || exit

# Install dependencies
npm install crypto-js

# Create the hidden config directory
CONFIG_DIR="$HOME/.cypher-cli"
mkdir -p "$CONFIG_DIR"

# Set up environment variable in ~/.cypher-cli/.env
echo 'BACKEND_URL="https://cypher-backend-harshiyers-projects.vercel.app"' > "$CONFIG_DIR/.env"

# Build the binary
go build -o cypher

echo "Cypher CLI setup complete."
echo "Make sure to add the binary to your PATH or run ./cypher from this directory."
