.PHONY: build install test clean run fmt lint help pre-commit setup-hooks

# 变量定义
BINARY_NAME=goon
BUILD_DIR=bin
MAIN_PATH=main.go
INSTALL_PATH=/usr/local/bin

# 默认目标
.DEFAULT_GOAL := help

## build: 编译项目
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## install: 安装到系统
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/
	@echo "Installation complete"

## test: 运行测试
test:
	@echo "Running tests..."
	@go test -v ./...

## test-coverage: 运行测试并生成覆盖率报告
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## clean: 清理构建产物
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

## run: 运行程序
run: build
	@$(BUILD_DIR)/$(BINARY_NAME)

## fmt: 格式化代码
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete"

## lint: 代码检查
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

## tidy: 整理依赖
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy
	@echo "Tidy complete"

## dev: 开发模式（格式化 + 构建 + 运行）
dev: fmt build
	@$(BUILD_DIR)/$(BINARY_NAME)

## pre-commit: 运行提交前检查
pre-commit:
	@echo "Running pre-commit checks..."
	@make fmt
	@make lint
	@make test

## setup-hooks: 设置 git hooks
setup-hooks:
	@echo "Setting up git hooks..."
	@echo "#!/bin/sh" > .git/hooks/pre-commit
	@echo "make pre-commit" >> .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "Git hooks installed successfully!"

## help: 显示帮助信息
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
