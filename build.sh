#!/usr/bin/env bash

set -x
set -euo pipefail

export CGO_ENABLED=0
export GOOS=windows

go build -v -tags=release -o ~/Downloads/ebda-cli.exe .

# APP_ID="earth.ebda.tools"
# APP_NAME="ebda"
# # IMAGE="ghcr.io/fyne-io/fyne-cross:latest"

# PLATFORMS=(
#   linux
#   windows
#   # darwin
# )

# build() {
#   local platform="$1"

#   doas fyne-cross "$platform" \
#     --app-id "$APP_ID" \
#     --name "$APP_NAME" 
#     # --image "$IMAGE"
# }

# for platform in "${PLATFORMS[@]}"; do
#   build "$platform"
# done

# for platform in "${PLATFORMS[@]}"; do
#   build "$platform" &
# done
# wait