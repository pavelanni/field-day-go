#!/bin/bash
# Check Markdown files for linting errors using markdownlint-cli2

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "Checking Markdown files for linting errors..."
echo

# Check if markdownlint-cli2 is installed
if ! command -v markdownlint-cli2 &> /dev/null; then
    echo -e "${RED}Error: markdownlint-cli2 not installed${NC}"
    echo "Install with: npm install -g markdownlint-cli2"
    exit 1
fi

# Find all Markdown files
MD_FILES=$(find . -name "*.md" -not -path "./node_modules/*" -not -path "./.git/*")

if [ -z "$MD_FILES" ]; then
    echo -e "${YELLOW}No Markdown files found${NC}"
    exit 0
fi

# Run markdownlint-cli2
if markdownlint-cli2 "**/*.md" "#node_modules" "#.git"; then
    echo -e "${GREEN}✓ All Markdown files pass linting${NC}"
    exit 0
else
    echo -e "${RED}✗ Markdown linting failed${NC}"
    echo
    echo "To see detailed errors, run:"
    echo "  markdownlint-cli2 specs/**/*.md"
    echo
    echo "Common fixes:"
    echo "  - Add blank lines before/after headings"
    echo "  - Add blank lines before/after code blocks"
    echo "  - Add blank lines before/after lists"
    echo "  - Add language to code fences (e.g., \`\`\`go)"
    echo "  - Use '1.' for all ordered list items"
    echo "  - Use '-' for unordered lists"
    exit 1
fi
