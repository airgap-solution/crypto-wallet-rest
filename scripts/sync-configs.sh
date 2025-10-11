#!/usr/bin/env bash
set -euo pipefail
REPO_NAME="infra-common"
COMMON_REPO="https://github.com/airgap-solution/${REPO_NAME}.git"
COMMON_DIR=".${REPO_NAME}"
if [ ! -d "$COMMON_DIR" ]; then
  echo "Cloning ${REPO_NAME}..."
  git clone "$COMMON_REPO" "$COMMON_DIR"
else
  echo "Updating ${REPO_NAME}..."
  git -C "$COMMON_DIR" pull --quiet
fi

echo "Syncing workflows..."
mkdir -p .github/workflows
cp -r "$COMMON_DIR/workflows/go/"* .github/workflows/
cp -r "$COMMON_DIR/workflows/ts/"* .github/workflows/
