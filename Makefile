SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Build

.PHONY: build
build: ## Compiles the source code.
	go build ./...

.PHONY: generate
generate: wire ## Compile-time Dependency Injection using code generation.
	$(WIRE) gen -tags "goverter wireinject" ./cmd/app/

.PHONY: lint
lint: golangci-lint ## Analyze and report style, formatting, and syntax issues in the source code.
	$(GOLANGCI_LINT) run ./...

##@ Tool Binaries

WIRE = $(shell pwd)/bin/wire
.PHONY: wire
wire: ## Checks for wire installation and downloads it if not found.
	$(call go-get-tool,$(WIRE),github.com/google/wire/cmd/wire@v0.6.0)

GOLANGCI_LINT = $(shell pwd)/bin/golangci-lint
.PHONY: golangci-lint
golangci-lint: ## Checks for golangci-lint installation and downloads it if not found.
	$(call go-get-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef
