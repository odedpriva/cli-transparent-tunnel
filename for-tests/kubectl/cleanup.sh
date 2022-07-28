#!/usr/bin/env bash

set -e

ROOT_DIR=$(git rev-parse --show-toplevel)
CURRENT_DIR="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
source "${ROOT_DIR}/for-tests/scripts/common.sh"
source "${CURRENT_DIR}/common.sh"

LOG INFO "downloading kind from ${KIND_GO_URL}"
GO_GET "${KIND_GO_URL}" kind
${KIND} delete cluster --name "${K8S_CLUSTER_NAME}"