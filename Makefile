# 引入环境变量
-include .env

# 项目基础变量
PROJECTNAME := $(shell basename "$(PWD)")
GOBASE      := $(shell pwd)
GOBIN       := $(GOBASE)/bin
GOFILES     := $(wildcard cmd/*/*.go)

# 获取 cmd 目录下所有的子目录名，作为构建目标
COMMANDS    := $(notdir $(patsubst %/,%,$(wildcard cmd/*/)))

# 编译输出重定向
STDERR      := /tmp/.$(PROJECTNAME)-stderr.txt
PID         := /tmp/.$(PROJECTNAME).pid

# 让 Make 输出更简洁
MAKEFLAGS += --silent

.PHONY: all build compile clean install test watch start stop help

default: help

## install: 安装并同步依赖 (go mod tidy)
install:
	@echo "  >  Syncing dependencies..."
	go mod tidy
	go mod download

## build: 编译所有 cmd 下的二进制文件
build: clean $(COMMANDS)

# 动态匹配 cmd/ 下的目录进行编译
$(COMMANDS):
	@echo "  >  Building binary: $@"
	GOBIN=$(GOBIN) go build -o $(GOBIN)/$@ ./cmd/$@

## compile: 快捷编译（默认编译主程序）
compile: build

## start: 开发模式：自动监控并重启
start:
	@bash -c "trap 'make stop' EXIT; $(MAKE) watch run='make compile start-server'"

start-server: stop-server
	@echo "  >  $(PROJECTNAME) is starting..."
	@$(GOBIN)/$(PROJECTNAME) 2>$(STDERR) & echo $$! > $(PID)
	@echo "  >  PID: $$(cat $(PID))"

stop-server:
	@-touch $(PID)
	@-kill `cat $(PID)` 2> /dev/null || true
	@-rm $(PID)

## watch: 使用 yolo 或内建逻辑监控文件变化
watch:
	@echo "  >  Watching files in $(GOBASE)..."
	@yolo -i . -e vendor -e bin -c "$(run)"

## test: 运行单元测试
test:
	@echo "  >  Running tests..."
	go test -v ./...

## clean: 清理二进制文件和缓存
clean:
	@echo "  >  Cleaning build cache and binaries..."
	go clean
	rm -rf $(GOBIN)
	rm -f $(STDERR)

## help: 显示所有命令
help: Makefile
	@echo
	@echo " Choose a command to run in $(PROJECTNAME):"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
