#!/usr/bin/env bash

ROOT_DIR=$(git rev-parse --show-toplevel)
source "${ROOT_DIR}/for-tests/scripts/common.sh"

export COMMAND=kubectl
export KUBECONFIG=${TEST_CONFIG_DIR}/k8s-config.yaml
export KIND="${BIN_DIR}/kind"
K8S_CLUSTER_NAME=${TOOL_NAME}

export K8S_CONTEXT_NAME=kind-${K8S_CLUSTER_NAME}
export KIND_CONTAINER_NAME=${K8S_CLUSTER_NAME}-control-plane
export KIND_GO_URL=sigs.k8s.io/kind@v0.14.0