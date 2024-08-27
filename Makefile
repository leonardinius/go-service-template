############################################################################################################
# Few rules to follow when editing this file:
#
# 1. Shell commands must be indented with a tab
# 2. Before each target add ".PHONY: target_name" to disable default file target
# 3. Add target description prefixed with "##" on the same line as target definition for "help" target to work
# 4. Be aware that each make command is executed in separate shell
#
# Tips:
#
# * Add an @ sign to suppress output of the command that is executed
# * Define variable like: VAR := value
# * Reference variable like: $(VAR)
# * Reference environment variables like: $(ENV_VARIABLE)
#
#############################################################################################################
.DELETE_ON_ERROR:
# TOOLS
GOLANGCILINT_VERSION 			= v1.60.1
GOFUMPT_VERSION		 			= v0.6.0
COMPILEDAEMON_VERSION			= v1.4.0
BUFF_VERSION					= v1.32.0
PROTOC_GEN_CONNECT_GO_VERSION 	= v1.16.2
PROTOC_GEN_VALIDATE_GO_VERSION 	= v1.1.0
PROTOC_GEN_GO_VERSION 			= v1.34.2
# SHELL
.SHELLFLAGS 	:= -eu -o pipefail -c
SHELL			= /bin/bash
BIN 			:= .bin
BUILDOUT		= ./bin
MAKEFLAGS 		+= --warn-undefined-variables
MAKEFLAGS 		+= --no-builtin-rules
MAKEFLAGS		+= --no-print-directory
BUILDTIME		?= $(shell date -u '+%Y%m%d%H%M%S')
REFNAME			?= $(shell git rev-parse --abbrev-ref HEAD)
COMMIT			?= $(shell [[ `git status --porcelain` ]] && echo dirty || git rev-parse --short HEAD)
GOPATH			?= ${shell go env GOPATH}
GOOS			?= $(shell go env GOHOSTOS)
GOARCH			?= $(shell go env GOHOSTARCH)
CGO_ENABLED 	?= 0
export GOBIN 	:= $(abspath $(BIN))

SERVICE_NAME	= service-template
# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS 	= -ldflags "\
 -X github.com/leonardinius/go-service-template/internal/services/version.ServiceName=${SERVICE_NAME} \
 -X github.com/leonardinius/go-service-template/internal/services/version.BuildTime=${BUILDTIME} \
 -X github.com/leonardinius/go-service-template/internal/services/version.RefName=${REFNAME} \
 -X github.com/leonardinius/go-service-template/internal/services/version.Commit=${COMMIT}"

BOLD = \033[1m
CLEAR = \033[0m
CYAN = \033[36m
GREEN = \033[32m

##@: Default
.PHONY: help
help: ## Display this help
	@awk '\
		BEGIN {FS = ":.*##"; printf "Usage: make $(CYAN)<target>$(CLEAR)\n"} \
		/^[0-9a-zA-Z_\-\/]+?:[^#]*?## .*$$/ { printf "\t$(CYAN)%-20s$(CLEAR) %s\n", $$1, $$2 } \
		/^##@/ { printf "$(BOLD)%s$(CLEAR)\n", substr($$0, 5); }' \
		$(MAKEFILE_LIST)

all: clean go/tidy go/format lint test e2e go/build ## Formats, builds and tests the codebase

##@: Build/Run
.PHONY: clean ## Cleans the build artifacts
clean: ## Format all go files
	@echo -e "$(CYAN)--- clean...$(CLEAR)"
	@go clean
	@rm -rf ${BUILDOUT}

.PHONY: gen
gen: api/generate ## Runs all codegen and docs tasks

.PHONY: test
test: clean go/test ## Runs all tests (excluding e2e)

.PHONY: e2e
e2e: clean go/e2e ## Runs e2e tests

.PHONY: lint
lint: go/lint api/lint ## Runs all linters

.PHONY: build
build: go/build ## Builds all artifacts

.PHONY: run
run: ## Runs service locally. Use ARGS="" make run to pass arguments
	@echo -e "$(CYAN)--- run ...$(CLEAR)"
	CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} GOARCH=${GOARCH} go run ${LDFLAGS} ./app/ $(ARGS)

.PHONY: watch
watch: $(BIN)/CompileDaemon ## Runs in watch mode. Example: `make watch ARGS="http"`
	@echo -e "$(CYAN)--- watch $(ARGS) ...$(CLEAR)"
	@echo -e "$(GREEN)Running 'CompileDaemon -build=\"make go/build\" -command=\"./bin/${SERVICE_NAME}-${GOOS}-${GOARCH}\"'...$(CLEAR)"
	@$(BIN)/CompileDaemon \
		-log-prefix=false \
		-graceful-kill=true \
		-color=true \
		-directory=. \
		-exclude-dir=.git \
		-exclude-dir=bin \
		-exclude-dir=docs \
		-exclude-dir=test \
		-include="*.go" -include=".env" \
		-build="$(MAKE) go/build" \
		-command="${BUILDOUT}/${SERVICE_NAME}-${GOOS}-${GOARCH} $(ARGS)"

###@: Go
.PHONY: go/format
go/format: $(BIN)/gofumpt ### Format all go files
	@echo -e "$(CYAN)--- format go files...$(CLEAR)"
	$(BIN)/gofumpt -w app/ internal/ teste2e/

go/tidy: go.mod go.sum ### Tidy all Go dependencies
	@echo -e "$(CYAN)--- tidy go dependencies...$(CLEAR)"
	go mod tidy -v -x

.PHONY: go/lint
go/lint: $(BIN)/golangci-lint ### Lints the codebase using golangci-lint
	@echo -e "$(CYAN)--- lint codebase...$(CLEAR)"
	$(BIN)/golangci-lint run --modules-download-mode=readonly --config .golangci.yml

.PHONY: go/build
go/build: ### Builds service
	@echo -e "$(CYAN)--- build ...$(CLEAR)"
	@mkdir -p ${BUILDOUT}
	CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BUILDOUT}/${SERVICE_NAME}-${GOOS}-${GOARCH} ./app/
	@echo -e "$(CYAN)--- binary at ${BUILDOUT}/${SERVICE_NAME}-${GOOS}-${GOARCH}$(CLEAR)"

.PHONY: go/test
go/test: ### Runs all tests
	@echo -e "$(CYAN)--- go test ...$(CLEAR)"
	go test -shuffle=on -race -cover -timeout=60s -count 1 -parallel 3 -v ./app/... ./internal/...

.PHONY: go/e2e
go/e2e: ### Runs e2e tests
	@echo -e "$(CYAN)--- go e2e test ...$(CLEAR)"
	go test -shuffle=on -race -cover -timeout=60s -count 1 -parallel 3 -v ./teste2e/...

###@: API spec
.PHONY: api/lint
api/lint: $(BIN)/buf ### Generates API spec
	@echo -e "$(CYAN)--- lint API spec...$(CLEAR)"
	$(BIN)/buf lint ./api/proto

.PHONY: api/breaking
api/breaking: $(BIN)/buf ### Checks for breaking changes in API spec
	@echo -e "$(CYAN)--- check for breaking changes in API spec...$(CLEAR)"
	$(BIN)/buf breaking /api/proto \
	  --against "https://github.com/leonardinius/go-service-template.git#branch=main,subdir=/api/proto" \

.PHONY: api/generate
api/generate: api/lint $(BIN)/buf $(BIN)/protoc-gen-go $(BIN)/protoc-gen-validate-go $(BIN)/protoc-gen-connect-go ### Generate go code from API spec
	@echo -e "$(CYAN)--- api generate go code and docs...$(CLEAR)"
	@rm -rf internal/apigen/* api/docs/*
	@echo "$(BIN)/buf generate --template ./api/buf.gen.yaml ./api/proto"
	@PATH="$(BIN):$(PATH)" $(BIN)/buf generate --template ./api/buf.gen.yaml ./api/proto

# TOOLS
$(BIN)/buf: Makefile
	@mkdir -p $(@D)
	go install github.com/bufbuild/buf/cmd/buf@$(BUFF_VERSION)

$(BIN)/protoc-gen-go: Makefile
	@mkdir -p $(@D)
	go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)

$(BIN)/protoc-gen-validate-go: Makefile
	@mkdir -p $(@D)
	go install github.com/envoyproxy/protoc-gen-validate/cmd/protoc-gen-validate-go@$(PROTOC_GEN_VALIDATE_GO_VERSION)

$(BIN)/protoc-gen-connect-go: Makefile
	@mkdir -p $(@D)
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go@$(PROTOC_GEN_CONNECT_GO_VERSION)

$(BIN)/golangci-lint: Makefile
	@mkdir -p $(@D)
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCILINT_VERSION)

$(BIN)/CompileDaemon: Makefile
	@mkdir -p $(@D)
	go install github.com/githubnemo/CompileDaemon@$(COMPILEDAEMON_VERSION)

$(BIN)/gofumpt: Makefile
	@mkdir -p $(@D)
	go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)

