#!/usr/bin/env bash

set -euo pipefail

log() {
  echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1"
}

if [ $# -ne 1 ]; then
  log "Error: Version argument required"
  echo "Usage: $0 <version>"
  exit 1
fi

VERSION=${1#v}
log "Processing version: ${VERSION}"

get_checksum() {
  local platform=$1
  local checksum
  checksum=$(grep "seed_${VERSION}_${platform}.tar.gz" dist/checksums.txt | head -n1 | awk '{print $1}')
  echo "$checksum"
}

# Get checksums with proper escaping
DARWIN_ARM64=$(get_checksum "darwin_arm64" | sed 's/[\/&]/\\&/g')
DARWIN_AMD64=$(get_checksum "darwin_amd64" | sed 's/[\/&]/\\&/g')
LINUX_ARM64=$(get_checksum "linux_arm64" | sed 's/[\/&]/\\&/g')
LINUX_AMD64=$(get_checksum "linux_amd64" | sed 's/[\/&]/\\&/g')

# Clone repository
log "Cloning homebrew-seed repository..."
git clone https://x-access-token:${GITHUB_TOKEN}@github.com/jpwallace22/homebrew-seed.git
cd homebrew-seed

log "Updating formula..."

# Update each SHA separately to avoid issues with escaping
sed -i.bak "s/version \".*\"/version \"${VERSION}\"/" Formula/seed.rb
sed -i.bak "/# darwin_arm64/ s/sha256 \"[^\"]*\"/sha256 \"${DARWIN_ARM64}\"/" Formula/seed.rb
sed -i.bak "/# darwin_amd64/ s/sha256 \"[^\"]*\"/sha256 \"${DARWIN_AMD64}\"/" Formula/seed.rb
sed -i.bak "/# linux_arm64/ s/sha256 \"[^\"]*\"/sha256 \"${LINUX_ARM64}\"/" Formula/seed.rb
sed -i.bak "/# linux_amd64/ s/sha256 \"[^\"]*\"/sha256 \"${LINUX_AMD64}\"/" Formula/seed.rb

log "Cleaning up..."
rm Formula/seed.rb.bak

# Verify the changes
log "Verifying changes..."
if ! grep -q "# darwin_arm64" Formula/seed.rb ||
  ! grep -q "# darwin_amd64" Formula/seed.rb ||
  ! grep -q "# linux_arm64" Formula/seed.rb ||
  ! grep -q "# linux_amd64" Formula/seed.rb; then
  log "Error: Platform comments were removed"
  exit 1
fi

# Configure git and commit
log "Committing changes..."
git config user.name "GitHub Action"
git config user.email "action@github.com"
git add Formula/seed.rb
git commit -m "Update formula to version ${VERSION}"
git push

log "Successfully updated Homebrew formula to version ${VERSION}"
