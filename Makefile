APP_NAME := iching-api
BUILD_DIR := bin
MAIN_PKG := ./cmd/api
GO := go
GOFMT := gofmt
GOLANGCI_LINT := $(CURDIR)/bin/golangci-lint
GOLANGCI_LINT_VERSION ?= v2.12.2
COMPOSE ?= docker compose

.PHONY: help build build-windows build-windows-arm64 run format fmt lint tidy clean test vet audit docker-up docker-down docker-logs

help:
	@echo "Available targets:"
	@echo "  make build       - build binary to ./bin/$(APP_NAME)"
	@echo "  make build-windows - build Windows binary"
	@echo "  make build-windows-arm64 - build Windows ARM64 binary"
	@echo "  make run         - run app from $(MAIN_PKG)"
	@echo "  make format      - format Go files with gofmt"
	@echo "  make fmt         - alias for format"
	@echo "  make lint        - run golangci-lint"
	@echo "  make tidy        - run go mod tidy"
	@echo "  make test        - run go test ./..."
	@echo "  make vet         - run go vet ./..."
	@echo "  make audit       - run format check stack: tidy, vet, lint, test"
	@echo "  make docker-up   - start PostgreSQL via docker compose"
	@echo "  make docker-down - stop PostgreSQL via docker compose"
	@echo "  make docker-logs - follow docker compose logs"
	@echo "  make clean       - remove build artifacts"

build:
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PKG)

build-windows:
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME).exe $(MAIN_PKG)

build-windows-arm64:
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME)-arm64.exe $(MAIN_PKG)
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
