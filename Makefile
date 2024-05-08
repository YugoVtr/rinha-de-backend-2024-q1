# Makefile based on:
# - https://gist.github.com/thomaspoignant/5b72d579bd5f311904d973652180c705
# - https://raw.githubusercontent.com/asyncapi/template-for-go-projects/master/Makefile

GOCMD=go
GOTEST=$(GOCMD) test
DOCKER=docker
BUILDX=$(DOCKER) buildx
PROJECT_NAME := $(shell basename "$(PWD)")
BINARY_NAME?=$(PROJECT_NAME)
BIN_DIR?=bin
TOOLS_DIR?=bin/tools

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

GOLANGCILINT_VERSION = 1.55.2

.PHONY: all test build

all: help

## Build:
build: ## Build your project and put the output binary in bin/out
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOCMD) build -o bin/out/$(BINARY_NAME) .

release: ## Build and push the Docker image
	$(BUILDX) build -t yugovtr/$(PROJECT_NAME):latest --push .

## Linting:
lint: $(TOOLS_DIR)/golangci-lint ## Run linters
	$(TOOLS_DIR)/golangci-lint run

## Test:
test: ## Run the tests of the project
	$(GOTEST) -v -race -count=1 ./...

integration-test: ## Run the integration tests of the project
	./scripts/integration-test.sh

load-test: ## Run the load tests of the project
	./scripts/load-test.sh

coverage: ## Run the tests of the project and export the coverage
	$(GOTEST) -cover -covermode=count -coverprofile=profile.cov ./...
	$(GOCMD) tool cover -func profile.cov

## Clean:
restart: ## Restart apps and truncate database
	./scripts/restart.sh

## App:
tilt: ## Run Tilt
	tmuxp load ./tmuxp.yml

download: ## Download the dependencies
	$(GOCMD) mod vendor
	./scripts/download-dependencies.sh

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

$(TOOLS_DIR):
	mkdir -p $(TOOLS_DIR)

$(TOOLS_DIR)/golangci-lint: $(TOOLS_DIR)
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | BINDIR=$(@D) sh -s -- -b $(TOOLS_DIR) v$(GOLANGCILINT_VERSION) > /dev/null 2>&1
