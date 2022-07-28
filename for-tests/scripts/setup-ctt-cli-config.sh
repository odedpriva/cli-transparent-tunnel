#!/usr/bin/env bash

set -e

ROOT_DIR=$(git rev-parse --show-toplevel)
CURRENT_DIR="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
source "${ROOT_DIR}/for-tests/scripts/common.sh"
source "${CURRENT_DIR}/common.sh"
source "${CURRENT_DIR}/../kubectl/common.sh"

GO_GET github.com/hairyhenderson/gomplate/v3/cmd/gomplate@latest gomplate


# shellcheck disable=SC2155
export KUBECTL_LOCATION="${BIN_DIR}/kubectl"

LOG INFO templating ctt tunnel config file to "${TEST_CONFIG_DIR}/ctt-cli-configuration.yaml"
${GO_TEMPLATE} -f "${CURRENT_DIR}/ctt-cli-configuration.yaml.tpl" > "${TEST_CONFIG_DIR}/ctt-cli-configuration.yaml"