# 引入环境变量
-include .env

# 项目基础变量
PROJECTNAME := $(shell basename "$(PWD)")
GOBASE      := $(shell pwd)
GOBIN       := $(GOBASE)/bin
GOPATH_BIN  := $(shell go env GOPATH)/bin
GOFILES     := $(wildcard cmd/*/*.go)

# 获取 cmd 目录下所有的子目录名，作为构建目标
COMMANDS    := $(notdir $(patsubst %/,%,$(wildcard cmd/*/)))

# 版本注入 (ldflags 只能设置 var，不能设置 const)
VERSION    ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT     ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE       ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# 编译标志 (类似 Cargo release profile)
LDFLAGS := -ldflags "-s -w -X hello/internal/version.VERSION=$(VERSION) -X hello/internal/version.COMMIT=$(COMMIT) -X hello/internal/version.BUILD_TIME=$(DATE)"
GOFLAGS := -trimpath -buildvcs=false $(LDFLAGS)

# CGO 显式启用 (go-sqlite3 需要)
CGO_ENABLED ?= 1
export CGO_ENABLED

# 编译输出重定向
STDERR      := /tmp/.$(PROJECTNAME)-stderr.txt
PID         := /tmp/.$(PROJECTNAME).pid

# 默认启动的 binary (可通过 BINARY=foo 覆盖)
BINARY ?= hello

# 让 Make 输出更简洁
MAKEFLAGS += --silent

.PHONY: all build compile clean install test watch start stop help fmt vet lint release cross doc run verify

default: help

## install: Install and sync dependencies (go mod tidy)
install:
	@echo "  >  Syncing dependencies..."
	go mod tidy
	go mod download
	go mod verify

## build: Build all binaries under cmd/
build: $(COMMANDS)

$(COMMANDS):
	@echo "  >  Building binary: $@"
	@mkdir -p $(GOBIN)
	GOBIN=$(GOBIN) go build $(GOFLAGS) -o $(GOBIN)/$@ ./cmd/$@

## compile: Alias for build
compile: build

## release: Optimized build
release:
	@echo "  >  Building release binaries (version: $(VERSION))..."
	@mkdir -p $(GOBIN)
	GOBIN=$(GOBIN) go build $(GOFLAGS) -o $(GOBIN)/hello ./cmd/hello
	GOBIN=$(GOBIN) go build $(GOFLAGS) -o $(GOBIN)/foo ./cmd/foo

## cross: Cross-compile for linux/darwin amd64+arm64
cross:
	@echo "  >  Cross-compiling..."
	@mkdir -p $(GOBIN)/cross
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o $(GOBIN)/cross/hello-linux-amd64 ./cmd/hello
	GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -o $(GOBIN)/cross/hello-darwin-arm64 ./cmd/hello
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o $(GOBIN)/cross/foo-linux-amd64 ./cmd/foo
	GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -o $(GOBIN)/cross/foo-darwin-arm64 ./cmd/foo
	@echo "  >  Cross-compiled binaries:"
	@ls -lh $(GOBIN)/cross/

## start: Dev mode with file watch and auto-restart
start:
	@bash -c "trap 'make stop' EXIT; $(MAKE) watch run='make compile start-server'"

start-server: stop-server build
	@echo "  >  Starting $(BINARY) (PID file: $(PID))..."
	@$(GOBIN)/$(BINARY) 2>$(STDERR) & echo $$! > $(PID)
	@echo "  >  PID: $$(cat $(PID))"

stop-server:
	@-touch $(PID)
	@-kill `cat $(PID)` 2> /dev/null || true
	@-rm $(PID)

## watch: Watch files and auto-restart (requires yolo)
watch:
	@echo "  >  Watching files in $(GOBASE)..."
	@yolo -i . -e vendor -e bin -c "$(run)"

## test: Run unit tests
test:
	@echo "  >  Running tests..."
	go test -v -count=1 ./...

## fmt: Format code
fmt:
	@echo "  >  Formatting code..."
	@gofmt -s -w $(GOFILES)

## vet: Static analysis
vet:
	@echo "  >  Running go vet..."
	@go vet ./...

## lint: Full linting
lint:
	@echo "  >  Running golangci-lint..."
	@$(GOPATH_BIN)/golangci-lint run ./...

## doc: Generate documentation
doc:
	@echo "  >  Generating documentation..."
	@go doc -all ./... 2>&1 | head -200

## run: Build and run
run: build
	@$(GOBIN)/$(BINARY)

## verify: Full quality gate (fmt + vet + test)
verify: fmt vet test
	@echo "  >  All checks passed ✓"

## clean: Clean build artifacts and cache
clean:
	@echo "  >  Cleaning build cache and binaries..."
	go clean
	rm -rf $(GOBIN)
	rm -f $(STDERR)

## help: Show available commands
help: Makefile
	@echo
	@echo " Choose a command to run in $(PROJECTNAME):"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
