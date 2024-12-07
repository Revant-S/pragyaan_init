#!/bin/bash

# Define source and destination directories
HOOKS_DIR=".githooks"
DEST_DIR=".git/hooks"

# Check if the .git directory exists
if [ ! -d ".git" ]; then
  echo "Error: This script must be run from the root of a Git repository."
  exit 1
fi

# Check if the hooks directory exists
if [ ! -d "$HOOKS_DIR" ]; then
  echo "Error: Hooks directory '$HOOKS_DIR' not found."
  exit 1
fi

# Copy hooks to the .git/hooks directory
echo "Copying hooks from '$HOOKS_DIR' to '$DEST_DIR'..."
cp -r "$HOOKS_DIR/"* "$DEST_DIR/"

# Make hooks executable
echo "Making hooks executable..."
chmod +x "$DEST_DIR/"*

echo "Git hooks have been set up successfully!"
