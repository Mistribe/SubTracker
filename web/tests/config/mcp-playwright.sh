#!/bin/bash

# MCP Playwright Server Launcher
# This script ensures NVM is loaded and runs the TypeScript MCP server

# Load NVM if available
export NVM_DIR="$HOME/.nvm"
if [ -s "$NVM_DIR/nvm.sh" ]; then
  . "$NVM_DIR/nvm.sh"
fi

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Change to project root
cd "$PROJECT_ROOT"

# Run the MCP server with ts-node
exec npx ts-node "$SCRIPT_DIR/mcp-playwright.ts"
