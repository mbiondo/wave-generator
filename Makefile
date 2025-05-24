# Colors
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RESET := \033[0m

# Variables
APP_NAME := wave-generator
GO_FILES := $(shell find . -name '*.go')
DOCKER_IMAGE := $(APP_NAME):latest
COVERAGE_FILE := coverage.out

# Targets
.PHONY: all build run clean test coverage coverage-html docker-build docker-run help

help: ## Show this help message
	@echo '${CYAN}Usage:${RESET}'
	@echo '  make ${YELLOW}<target>${RESET}'
	@echo ''
	@echo '${CYAN}Targets:${RESET}'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  ${YELLOW}%-15s${RESET} %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: build ## Build the application (default)

build: $(GO_FILES) ## Build the Go application
	@echo "${GREEN}Building the Go application...${RESET}"
	@go build -o $(APP_NAME) main.go

run: build ## Run the application
	@echo "${GREEN}Running the application...${RESET}"
	@./$(APP_NAME)

test: ## Run tests
	@echo "${GREEN}Running tests...${RESET}"
	@go test -v ./...

coverage: ## Run tests with coverage
	@echo "${GREEN}Running tests with coverage...${RESET}"
	@go test -v -coverprofile=$(COVERAGE_FILE) ./...
	@go tool cover -func=$(COVERAGE_FILE)

coverage-html: coverage ## Show coverage report in browser
	@echo "${GREEN}Opening coverage report in browser...${RESET}"
	@go tool cover -html=$(COVERAGE_FILE)

clean: ## Clean up build artifacts
	@echo "${GREEN}Cleaning up...${RESET}"
	@rm -f $(APP_NAME)
	@rm -f $(COVERAGE_FILE)

docker-build: ## Build Docker image
	@echo "${GREEN}Building Docker image...${RESET}"
	@docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run Docker container
	@echo "${GREEN}Running Docker container...${RESET}"
	@docker run -p 1155:1155 $(DOCKER_IMAGE)