#!/usr/bin/env bash

set -e

ROOT_DIR=$(git rev-parse --show-toplevel)
CURRENT_DIR="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
source "${ROOT_DIR}/for-tests/scripts/common.sh"
source "${CURRENT_DIR}/common.sh"

LOG INFO "setting up ${COMMAND}"


KUBECTL_VERSION="v1.24.0"
KUBECTL_DOWNLOAD_URL="https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/${OS}/${ARCH}/kubectl"
LOG INFO "downloading kubectl from ${KUBECTL_DOWNLOAD_URL}"
CURL_GET "${KUBECTL_DOWNLOAD_URL}" kubectl

LOG INFO "downloading kind from ${KIND_GO_URL}"
GO_GET "${KIND_GO_URL}" kind

LOG INFO "deleting kind cluster ${K8S_CLUSTER_NAME}"
${KIND} delete cluster --name "${K8S_CLUSTER_NAME}"


kind_config_file=${TEST_CONFIG_DIR}/kind-config.yaml
LOG INFO "templating king config file to ${kind_config_file}"
${GO_TEMPLATE} --file "${CURRENT_DIR}/kind-config.yaml.tpl" > "${kind_config_file}"
LOG INFO "creating kind cluster name: ${K8S_CLUSTER_NAME} kubeconfig: ${KUBECONFIG}"
${KIND} create cluster --name "${K8S_CLUSTER_NAME}" --kubeconfig "${KUBECONFIG}" \
  --config "${kind_config_file}"

exit 0