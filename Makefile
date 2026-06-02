APP_NAME := iching-api
BUILD_DIR := bin
DIST_DIR := dist
MAIN_PKG := ./cmd/api/main.go
GO := go
GOFMT := gofmt
GOLANGCI_LINT := $(CURDIR)/bin/golangci-lint
GOLANGCI_LINT_VERSION ?= v2.12.2
COMPOSE ?= docker compose
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

.PHONY: help build build-linux build-windows build-windows-arm64 build-darwin build-all release-build run format fmt lint tidy clean test vet audit docker-up docker-down docker-logs

help:
	@echo "Available targets:"
	@echo "  make build              - build native binary to ./bin/$(APP_NAME)"
	@echo "  make build-linux        - build Linux x86_64 binary"
	@echo "  make build-linux-arm64  - build Linux ARM64 binary"
	@echo "  make build-windows      - build Windows x86_64 binary"
	@echo "  make build-windows-arm64 - build Windows ARM64 binary"
	@echo "  make build-darwin       - build macOS binary"
	@echo "  make build-all          - build all platforms"
	@echo "  make release-build      - prepare release with versioning"
	@echo "  make run                - run app from $(MAIN_PKG)"
	@echo "  make format             - format Go files with gofmt"
	@echo "  make fmt                - alias for format"
	@echo "  make lint               - run golangci-lint"
	@echo "  make tidy               - run go mod tidy"
	@echo "  make test               - run go test ./..."
	@echo "  make vet                - run go vet ./..."
	@echo "  make audit              - run format check stack: tidy, vet, lint, test"
	@echo "  make docker-up          - start PostgreSQL via docker compose"
	@echo "  make docker-down        - stop PostgreSQL via docker compose"
	@echo "  make docker-logs        - follow docker compose logs"
	@echo "  make clean              - remove build artifacts"
	@echo ""
	@echo "Release workflow:"
	@echo "  1. make audit           - run tests and lint"
	@echo "  2. git tag -a v0.1.0 -m 'Release v0.1.0'"
	@echo "  3. git push origin v0.1.0"
	@echo "  4. make release-build   - builds all platforms and creates dist/"

build:
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PKG)

build-linux:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(MAIN_PKG)

build-linux-arm64:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64 $(MAIN_PKG)

build-windows:
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(MAIN_PKG)

build-windows-arm64:
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME)-windows-arm64.exe $(MAIN_PKG)

build-darwin:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(MAIN_PKG)

build-all: build-linux build-linux-arm64 build-windows build-windows-arm64 build-darwin
	@echo "✓ All platforms built successfully"
	@ls -lh $(BUILD_DIR)/$(APP_NAME)-*

release-build: clean audit build-all
	@mkdir -p $(DIST_DIR)
	@echo "Version: $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo ""
	@echo "Preparing distribution packages..."
	@# Linux x86_64
	@tar -czf $(DIST_DIR)/$(APP_NAME)-linux-amd64-$(VERSION).tar.gz -C $(BUILD_DIR) $(APP_NAME)-linux-amd64 -C ../cmd/api static/
	@echo "✓ $(DIST_DIR)/$(APP_NAME)-linux-amd64-$(VERSION).tar.gz"
	@# Linux ARM64
	@tar -czf $(DIST_DIR)/$(APP_NAME)-linux-arm64-$(VERSION).tar.gz -C $(BUILD_DIR) $(APP_NAME)-linux-arm64 -C ../cmd/api static/
	@echo "✓ $(DIST_DIR)/$(APP_NAME)-linux-arm64-$(VERSION).tar.gz"
	@# Windows x86_64
	@cd $(BUILD_DIR) && zip -j ../$(DIST_DIR)/$(APP_NAME)-windows-amd64-$(VERSION).zip $(APP_NAME)-windows-amd64.exe && cd ..
	@cp -r cmd/api/static $(DIST_DIR)/$(APP_NAME)-windows-amd64-$(VERSION)-static/ 2>/dev/null || true
	@echo "✓ $(DIST_DIR)/$(APP_NAME)-windows-amd64-$(VERSION).zip"
	@# Windows ARM64
	@cd $(BUILD_DIR) && zip -j ../$(DIST_DIR)/$(APP_NAME)-windows-arm64-$(VERSION).zip $(APP_NAME)-windows-arm64.exe && cd ..
	@echo "✓ $(DIST_DIR)/$(APP_NAME)-windows-arm64-$(VERSION).zip"
	@# macOS x86_64
	@tar -czf $(DIST_DIR)/$(APP_NAME)-darwin-amd64-$(VERSION).tar.gz -C $(BUILD_DIR) $(APP_NAME)-darwin-amd64 -C ../cmd/api static/
	@echo "✓ $(DIST_DIR)/$(APP_NAME)-darwin-amd64-$(VERSION).tar.gz"
	@echo ""
	@echo "✓ Release packages ready in $(DIST_DIR)/"
run:
	$(GO) run $(MAIN_PKG)

format:
	$(GOFMT) -w $$(find . -type f -name '*.go' -not -path './bin/*')

fmt: format

$(GOLANGCI_LINT):
	@mkdir -p $(dir $(GOLANGCI_LINT))
	@command -v curl >/dev/null 2>&1 && \
		curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(dir $(GOLANGCI_LINT)) $(GOLANGCI_LINT_VERSION) || \
		wget -O- -nv https://golangci-lint.run/install.sh | sh -s -- -b $(dir $(GOLANGCI_LINT)) $(GOLANGCI_LINT_VERSION)

lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run ./...

tidy:
	$(GO) mod tidy

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

audit: tidy vet lint test

docker-up:
	$(COMPOSE) up -d

docker-down:
	$(COMPOSE) down

docker-logs:
	$(COMPOSE) logs -f

clean:
	rm -rf $(BUILD_DIR)
