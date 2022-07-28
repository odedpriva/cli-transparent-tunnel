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
export SSH_TUNNEL_SERVER=${SSH_SERVER_USER}@localhost:${SSH_SERVER_PORT}
export SSH_PRIVATE_KEY_PATH="${SSH_PRIVATE_KEY_PATH}"
# shellcheck disable=SC2155
export K8S_ENDPOINT=${LOCAL_IP}:$(docker inspect --format='{{range $p, $conf := .NetworkSettings.Ports}}{{(index $conf 0).HostPort}} {{end}}' "${KIND_CONTAINER_NAME}")
export K8S_CLUSTER_NAME=${K8S_CONTEXT_NAME}

LOG INFO templating ctt tunnel config file to "${TEST_CONFIG_DIR}/ctt-tunnel-configuration.yaml"
${GO_TEMPLATE} -f "${CURRENT_DIR}/ctt-tunnel-configuration.yaml.tpl" > "${TEST_CONFIG_DIR}/ctt-tunnel-configuration.yaml"