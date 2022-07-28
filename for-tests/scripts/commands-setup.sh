#!/usr/bin/env bash

set -e

ROOT_DIR=$(git rev-parse --show-toplevel)
CURRENT_DIR="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
source "${ROOT_DIR}/for-tests/scripts/common.sh"
source "${CURRENT_DIR}/common.sh"

GO_GET github.com/hairyhenderson/gomplate/v3/cmd/gomplate@latest gomplate

for command in "${commands_list[@]}"; do
   "${TEST_DIR}/${command}/setup.sh"
done