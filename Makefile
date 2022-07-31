SHELL=bash
GOCMD=go
export GO111MODULE=on
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
DOCKER=docker
DOCKERFILE=Dockerfile
TEST_DOCKERFILE=test -f $(DOCKERFILE)
NO_DOCKERFILE=echo No $(DOCKERFILE) available
REG=$(reg)
HADOLINT=dockerregistry-v2.vih.infineon.com/hadolint/hadolint:2.10.0-alpine
GOLINT=golangci/golangci-lint:v1.47-alpine
YAMLLINT=cytopia/yamllint:1.26
MAKE=make --no-print-directory
CONF=config
OUT=out
COV_HTML_REP=$(OUT)/coverage_html_report
BINARY_NAME=projects
BINARY_PATH=$(OUT)/bin/$(BINARY_NAME)
VERSION?=0.0.0
SERVICE_PORT?=3000
DOCKER_REGISTRY?= #if set it should finished by /

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all test build tidy vendor

all: help

## Run:
run: vendor ## Run the current module, r is likely better
	@go run .

r: build ## Run the built program with all given arguments
# Simpler solution: just prepend out/bin/ to PATH and easily call binary
	@printf "Run: $(BINARY_NAME) " && IFS=' ' read -r -a args && echo "INFO Arguments used: " $${args[@]} && echo && $(BINARY_PATH) $${args[@]}

debug: ## Start a headless debugger
	@printf "Debug: $(BINARY_NAME) " && IFS=' ' read -r -a args && echo "INFO Arguments used: " $${args[@]} && echo && dlv debug . --headless --output=out/bin/debug_bin --listen=:5678 --api-version=2 -- $${args[@]}

## Build:
build: vendor ## Build your project and put the output binary in out/bin/
	@mkdir -p $(OUT)/bin
	$(GOCMD) build -mod vendor -o $(BINARY_PATH) .

watch-build: ## Run a build everytime a file changes (requires entr, does not rerun if there is a new file)
	@find . -path ./vendor -prune -o -name '*_test.go' -prune -o -name '*.go' -print | entr -c $(MAKE) build

clean: ## Remove build related files
	rm -rf $(OUT)

tidy: ## Tidy go.mod
	@$(GOCMD) mod tidy

vendor: tidy ## Copy of all packages needed to support builds and tests in the vendor directory
	@$(GOCMD) mod vendor

## Test:
test: ## Run the tests of the project
	$(GOTEST) -v -race ./...

watch-tests: ## Run the tests everytime go files change (requires entr, does not rerun if there is a new file)
	@find . -path ./vendor -prune -o -type f -name '*.go' -print | entr -c $(MAKE) test

coverage: ## Run the tests of the project and export the coverage
	$(GOTEST) -cover -coverprofile=$(OUT)/coverage.out ./...
	mkdir -p $(COV_HTML_REP)
	$(GOCMD) tool cover -html $(OUT)/coverage.out -o $(COV_HTML_REP)/index.html
	@echo You can view the coverage in $(COV_HTML_REP)/index.html

## Lint:
lint: lint-go lint-dockerfile lint-yaml ## Run all available linters

lint-dockerfile: ## Lint your Dockerfile
	$(TEST_DOCKERFILE) && $(DOCKER) run --rm -i $(REG)/$(HADOLINT) < $(DOCKERFILE) || $(NO_DOCKERFILE)

lint-go: ## Use golintci-lint on your project
	$(DOCKER) run --rm -v $(shell pwd):/code -w /code $(REG)/$(GOLINT) golangci-lint run --config config/.golangci.yaml --deadline=60s 

lint-yaml: ## Use yamllint on the yaml file of your projects
	$(DOCKER) run --rm -it -v $(shell pwd):/data -w /data $(REG)/$(YAMLLINT) --format colored --config-file $(CONF)/.yamllint.yaml $(shell find -name '*.yml' -o -name '*.yaml')

## Docker:
docker-build: ## Use the dockerfile to build the container
	$(TEST_DOCKERFILE) && $(DOCKER) build --rm --tag $(BINARY_NAME) . || $(NO_DOCKERFILE)

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