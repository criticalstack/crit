BIN_DIR := bin

LDFLAGS := -s -w
GOFLAGS = -gcflags "all=-trimpath=$(PWD)" -asmflags "all=-trimpath=$(PWD)"
GO_BUILD_ENV_VARS := GO111MODULE=on CGO_ENABLED=0

tfconfig: ## Build the tfconfig binary
	@$(GO_BUILD_ENV_VARS) go build -o $(BIN_DIR)/tfconfig $(GOFLAGS) -ldflags '$(LDFLAGS)' ./cmd/tfconfig
	@go mod tidy

test:
	@go test ./pkg/...
