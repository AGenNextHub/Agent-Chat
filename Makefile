# AGenNext Chat — developer workflow.
# Mirrors the CI gates so "it builds on my machine" matches the pipeline.

BINARY      := agennextd
PKG         := ./...
IMAGE       ?= ghcr.io/agennext/agent-chat
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

.DEFAULT_GOAL := check

.PHONY: check
check: tidy vet lint test ## Run the full local gate (tidy, vet, lint, test)

.PHONY: build
build: ## Build the binary
	CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o bin/$(BINARY) ./cmd/agennextd

.PHONY: run
run: ## Run the end-to-end demo
	go run ./cmd/agennextd

.PHONY: test
test: ## Run tests with race detector and coverage
	go test -race -covermode=atomic -coverprofile=coverage.out $(PKG)

.PHONY: cover
cover: test ## Show coverage summary
	go tool cover -func=coverage.out | tail -1

.PHONY: vet
vet: ## go vet
	go vet $(PKG)

.PHONY: lint
lint: ## golangci-lint
	golangci-lint run $(PKG)

.PHONY: tidy
tidy: ## Ensure go.mod is tidy
	go mod tidy
	git diff --exit-code go.mod

.PHONY: vuln
vuln: ## Supply-chain: vulnerability scan (govulncheck)
	go run golang.org/x/vuln/cmd/govulncheck@latest $(PKG)

.PHONY: sbom
sbom: ## Supply-chain: generate a CycloneDX SBOM (requires syft)
	syft dir:. -o cyclonedx-json=sbom.cdx.json

.PHONY: docker
docker: ## Build the container image
	docker build -t $(IMAGE):$(VERSION) .

.PHONY: helm-lint
helm-lint: ## Lint the Helm chart
	helm lint deploy/helm/agennext-chat

.PHONY: help
help: ## List targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'
