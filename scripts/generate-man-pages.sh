#!/bin/bash

# Change to the project root directory
cd "$(dirname "$0")/.."

# Create the man directory if it doesn't exist
mkdir -p man

# Create a temporary file to hold all the markdown content
tmp_file=$(mktemp)

# Concatenate all markdown files into the temporary file
find docs/content -name "*.md" -print0 | sort -z | xargs -0 cat > "$tmp_file"

# Convert the temporary file to a single man page
pandoc -s -f markdown -t man "$tmp_file" -o "man/evoke.1"

# Remove the temporary file
rm "$tmp_file"
