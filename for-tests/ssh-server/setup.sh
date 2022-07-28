#!/usr/bin/env bash

set -e

ROOT_DIR=$(git rev-parse --show-toplevel)
CURRENT_DIR="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
source "${ROOT_DIR}/for-tests/scripts/common.sh"
source "${CURRENT_DIR}/common.sh"

image_name=ssh-server
${DOCKER} build -f "${CURRENT_DIR}/Dockerfile" -t ${image_name} "${CURRENT_DIR}"

${DOCKER} rm -f "${SSH_SERVER_CONTAINER_NAME}"
LOG INFO "using public key ${SSH_PUBLIC_KEY}"
${DOCKER} run -d --name="${SSH_SERVER_CONTAINER_NAME}" -e PUBLIC_KEY="$(cat -e "${SSH_PUBLIC_KEY}")" -p "${SSH_SERVER_PORT}":2222 ${image_name}

# shellcheck disable=SC2046
LOG INFO $(${DOCKER} inspect -f '{{.Id}} {{.Name}} {{.NetworkSettings.IPAddress}} {{.NetworkSettings.Ports}}' "${SSH_SERVER_CONTAINER_NAME}")