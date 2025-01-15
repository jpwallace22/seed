#!/usr/bin/env bash

set -e

VERSION=${1#v} # Remove 'v' prefix if present

if [ ! -f dist/checksums.txt ]; then
  echo "Checksums file not found"
  exit 1
fi

# Get macOS checksums
DARWIN_ARM64_SHA=$(cat dist/checksums.txt | grep darwin_arm64.tar.gz | awk '{print $1}')
DARWIN_AMD64_SHA=$(cat dist/checksums.txt | grep darwin_amd64.tar.gz | awk '{print $1}')

# Get Linux checksums
LINUX_ARM64_SHA=$(cat dist/checksums.txt | grep linux_arm64.tar.gz | awk '{print $1}')
LINUX_AMD64_SHA=$(cat dist/checksums.txt | grep linux_amd64.tar.gz | awk '{print $1}')

# Clone and update formula
git clone https://x-access-token:${GITHUB_TOKEN}@github.com/jpwallace22/homebrew-seed.git
cd homebrew-seed

# Update version
sed -i.bak "s/version \".*\"/version \"${VERSION}\"/" Formula/seed.rb

# Update all checksums
sed -i.bak "s/sha256 \".*\".*darwin_arm64/sha256 \"${DARWIN_ARM64_SHA}\"/" Formula/seed.rb
sed -i.bak "s/sha256 \".*\".*darwin_amd64/sha256 \"${DARWIN_AMD64_SHA}\"/" Formula/seed.rb
sed -i.bak "s/sha256 \".*\".*linux_arm64/sha256 \"${LINUX_ARM64_SHA}\"/" Formula/seed.rb
sed -i.bak "s/sha256 \".*\".*linux_amd64/sha256 \"${LINUX_AMD64_SHA}\"/" Formula/seed.rb

# Verify and commit
git config user.name 'GitHub Action'
git config user.email 'action@github.com'
git add Formula/seed.rb
git commit -m "Update formula for version ${VERSION}"
git push
