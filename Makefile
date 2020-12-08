.DEFAULT_GOAL:=help

BIN_DIR       ?= bin
TOOLS_DIR     := hack/tools
TOOLS_BIN_DIR := $(TOOLS_DIR)/bin
GOLANGCI_LINT := $(TOOLS_BIN_DIR)/golangci-lint
GENDOCS       := $(TOOLS_BIN_DIR)/gendocs
GENTMPLS      := $(TOOLS_BIN_DIR)/gentmpls
REPO_ROOT     := github.com/criticalstack/crit

GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD | sed 's/\///g')
GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

VERSION = $(GIT_BRANCH).$(GIT_SHA)
ifneq ($(GIT_TAG),)
	VERSION = $(GIT_TAG)
endif

LDFLAGS := -s -w
LDFLAGS += -X "$(REPO_ROOT)/internal/buildinfo.Date=$(shell date -u +'%Y-%m-%dT%TZ')"
LDFLAGS += -X "$(REPO_ROOT)/internal/buildinfo.GitSHA=$(GIT_SHA)"
LDFLAGS += -X "$(REPO_ROOT)/internal/buildinfo.GitTreeState=$(GIT_DIRTY)"
LDFLAGS += -X "$(REPO_ROOT)/internal/buildinfo.Version=$(VERSION)"
GOFLAGS = -gcflags "all=-trimpath=$(PWD)" -asmflags "all=-trimpath=$(PWD)"

GO_BUILD_ENV_VARS := GO111MODULE=on CGO_ENABLED=0

##@ Building

.PHONY: crit cinder

cinder: ## Build the cinder binary
	@$(GO_BUILD_ENV_VARS) go build -o $(BIN_DIR)/cinder $(GOFLAGS) -ldflags '$(LDFLAGS)' ./cmd/cinder

crit: ## Build the crit binary
	@$(GO_BUILD_ENV_VARS) GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/crit $(GOFLAGS) -ldflags '$(LDFLAGS)' ./cmd/crit

.PHONY: update-codegen update-charts update-docs update-embedded
update-codegen: ## Update generated code (slow)
	@echo "Updating generated code files ..."
	@echo "  *** This can be slow and does not need to run every build ***"
	@hack/tools/update-codegen.sh

update-charts: ## Update helm chart templates
	@echo "Updating helm charts ..."
	@helm template hack/charts/coredns --namespace kube-system \
		--name-template coredns -f hack/charts/coredns-values.yaml > templates/coredns.yaml

update-docs: clean $(GENDOCS)
	@echo "Generating CLI docs ..."
	$(GENDOCS) ./docs/src
	@echo "Building mdbook ..."
	$(MAKE) -s -C ./docs book

update-embedded: $(GENTMPLS) ## Update embedded templates
	@echo "Updating embedded template files ..."
	$(GENTMPLS) ./templates ./pkg/cluster

all: clean update-charts update-embedded update-codegen crit ## Generate all files and build crit binary (slow)

##@ Testing

.PHONY: lint

test: ## Run all tests
	go test ./pkg/...

lint: $(GOLANGCI_LINT) ## Lint codebase
	$(GOLANGCI_LINT) run -v

lint-full: $(GOLANGCI_LINT) ## Run slower linters to detect possible issues
	$(GOLANGCI_LINT) run -v --fast=false

##@ Helpers

.PHONY: help

$(GENDOCS): # Build gendocs from tools folder.
	@cd $(TOOLS_DIR); go build -tags=tools -o bin/gendocs ./gendocs

$(GENTMPLS): # Build gentmpls from tools folder.
	@cd $(TOOLS_DIR); go build -tags=tools -o bin/gentmpls ./gentmpls

$(GOLANGCI_LINT): $(TOOLS_DIR)/go.mod # Build golangci-lint from tools folder.
	cd $(TOOLS_DIR); go build -tags=tools -o bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint

clean: ## Cleanup the project folders
	@rm -f $(BIN_DIR)/*
	@rm -f $(TOOLS_BIN_DIR)/*

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
