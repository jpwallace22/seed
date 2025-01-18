#!/usr/bin/env bash

# Exit on any error
set -euo pipefail

# Check if version argument is provided
if [ $# -ne 1 ]; then
  echo "Usage: $0 <version>"
  exit 1
fi

# Remove 'v' prefix if present
VERSION=${1#v}

# Check for checksums file
if [ ! -f dist/checksums.txt ]; then
  echo "Error: checksums.txt not found in dist directory"
  exit 1
fi

# Function to extract SHA for a specific platform
get_sha() {
  local platform=$1
  local sha=$(grep "${platform}.tar.gz" dist/checksums.txt | awk '{print $1}')
  if [ -z "$sha" ]; then
    echo "Error: Could not find SHA for ${platform}" >&2
    exit 1
  fi
  echo "$sha"
}

# Get all checksums
DARWIN_ARM64_SHA=$(get_sha "darwin_arm64")
DARWIN_AMD64_SHA=$(get_sha "darwin_amd64")
LINUX_ARM64_SHA=$(get_sha "linux_arm64")
LINUX_AMD64_SHA=$(get_sha "linux_amd64")

# Clone repository
echo "Cloning homebrew-seed repository..."
git clone https://x-access-token:${GITHUB_TOKEN}@github.com/jpwallace22/homebrew-seed.git
cd homebrew-seed

# Create temporary file
TMP_FILE=$(mktemp)
FORMULA_FILE="Formula/seed.rb"

# Read the formula file line by line and make replacements
while IFS= read -r line; do
  # Update version
  if [[ $line =~ ^[[:space:]]*version[[:space:]] ]]; then
    echo "  version \"${VERSION}\"" >>"$TMP_FILE"
  # Update SHA for darwin_arm64
  elif [[ $line =~ .*darwin_arm64.*sha256.* ]]; then
    echo "      sha256 \"${DARWIN_ARM64_SHA}\"" >>"$TMP_FILE"
  # Update SHA for darwin_amd64
  elif [[ $line =~ .*darwin_amd64.*sha256.* ]]; then
    echo "      sha256 \"${DARWIN_AMD64_SHA}\"" >>"$TMP_FILE"
  # Update SHA for linux_arm64
  elif [[ $line =~ .*linux_arm64.*sha256.* ]]; then
    echo "      sha256 \"${LINUX_ARM64_SHA}\"" >>"$TMP_FILE"
  # Update SHA for linux_amd64
  elif [[ $line =~ .*linux_amd64.*sha256.* ]]; then
    echo "      sha256 \"${LINUX_AMD64_SHA}\"" >>"$TMP_FILE"
  else
    echo "$line" >>"$TMP_FILE"
  fi
done <"$FORMULA_FILE"

# Replace original file with updated content
mv "$TMP_FILE" "$FORMULA_FILE"

# Verify the changes
echo "Verifying changes..."
if ! grep -q "version \"${VERSION}\"" "$FORMULA_FILE"; then
  echo "Error: Version update failed"
  exit 1
fi

if ! grep -q "$DARWIN_ARM64_SHA" "$FORMULA_FILE"; then
  echo "Error: Darwin ARM64 SHA update failed"
  exit 1
fi

# Configure git
git config user.name 'GitHub Action'
git config user.email 'action@github.com'

# Commit and push changes
git add "$FORMULA_FILE"
git commit -m "Update formula for version ${VERSION}"
git push

echo "Successfully updated formula to version ${VERSION}"
