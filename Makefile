GIT_COMMIT:=$(shell git describe --dirty --always)
GIT_TAG:=$(shell git describe --dirty --always --tags)
PKG:=github.com/odedpriva/cli-transparent-tunnel/version

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Find current user version of Go, set golangci-lint version accordingly
GOLANGCI_VER= 1.43.0
GO_VER = $(shell go version | awk '{ print $$3 }' | awk -F '.' '{ print $$2 }')

ifeq ($(shell expr $(GO_VER) \> 17), 1)
GOLANGCI_VER = 1.45.2
else
GOLANGCI_VER = 1.43.0
endif

SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

all: build

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

integration-setup: build ## Setup the integration test openssh-server, kind
	@./for-tests/ssh-server/setup.sh
	@./for-tests/scripts/commands-setup.sh
	@./for-tests/scripts/setup-ctt-cli-config.sh
	@./for-tests/scripts/setup-ctt-tunnel-config.sh


integration-cleanup: ## Cleanup the integration
	@./for-tests/ssh-server/cleanup.sh
	@./for-tests/scripts/commands-cleanup.sh

integration-build:  ## Build integration tests
	go test -c for-tests/*/*.go -o bin/integration-tests

integration-tests: integration-setup integration-build
	./bin/integration-tests -test.v

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

mod:
	go mod download

tidy:
	go mod tidy

golangci-lint: ## Download golangci-lint
	@mkdir -p $(shell pwd)/bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell pwd)/bin v$(GOLANGCI_VER)

.PHONY: lint
lint: golangci-lint ## Run linter
	$(shell pwd)/bin/golangci-lint run

.PHONY: goimports
goimports: ## Run goimports updating files in place
	goimports -w .

.PHONY: unit-test
unit-test: ## Run Unit tests
	go test ./...

##@ Build

.DEFAULT: build
build:  fmt vet lint ## Build manager binary.
	go build -ldflags="-s -w -X ${PKG}.GitVersion=${GIT_TAG} -X ${PKG}.GitCommit=${GIT_COMMIT}" -o bin/ctt cmd/main.go


##@ Release

release: goreleaser ## release binary
	$(GORELEASER) release --rm-dist


GORELEASER = $(shell pwd)/bin/goreleaser
goreleaser: ## Download kustomize locally if necessary.
	$(call go-get-tool,$(GORELEASER),github.com/goreleaser/goreleaser@latest)


GOGET_CMD = "install"
ifeq ($(shell expr $(GO_VER) \< 16), 1)
GOGET_CMD = "get"
endif

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go $(GOGET_CMD) $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef