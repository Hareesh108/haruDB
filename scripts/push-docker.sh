#!/bin/bash
#
# Build and push HaruDB Docker image with both server & CLI binaries.
# Usage: ./push-docker.sh v0.0.4
#

set -euo pipefail

VERSION="${1:-}"
IMAGE_NAME="hareesh108/harudb"

if [[ -z "$VERSION" ]]; then
  echo "❌ Usage: $0 <version-tag>  (e.g. $0 v0.0.4)"
  exit 1
fi

echo "🚀 Building HaruDB Docker image for $VERSION ..."

# Build the versioned image
docker build  -t "$IMAGE_NAME:$VERSION" ../

# Tag as latest
docker tag "$IMAGE_NAME:$VERSION" "$IMAGE_NAME:latest"

echo "📤 Pushing to Docker Hub ..."
docker push "$IMAGE_NAME:$VERSION"
docker push "$IMAGE_NAME:latest"

echo "✅ Done: pushed $IMAGE_NAME:$VERSION and latest"
