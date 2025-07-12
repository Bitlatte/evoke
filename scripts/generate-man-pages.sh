#!/bin/bash

# Change to the project root directory
cd "$(dirname "$0")/.."

# Create the man directory if it doesn't exist
mkdir -p man

# Find all markdown files in docs/content and convert them to man pages
find docs/content -name "*.md" | while read -r file; do
  # Get the base name of the file
  base_name=$(basename "$file" .md)
  # Convert the file to a man page
  pandoc -s -f markdown -t man "$file" -o "man/evoke-$base_name.1"
done
