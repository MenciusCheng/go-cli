# 项目名称
PROJECT_NAME := go-cli

# 版本信息
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "v1.0.0")
BUILD_TIME := $(shell date -u '+%Y-%m-%d %H:%M:%S UTC')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 构建标志 - 注入到 cmd 包的变量中
LDFLAGS := -w -s \
	-X 'github.com/MenciusCheng/go-cli/cmd.Version=$(VERSION)' \
	-X 'github.com/MenciusCheng/go-cli/cmd.BuildTime=$(BUILD_TIME)' \
	-X 'github.com/MenciusCheng/go-cli/cmd.GitCommit=$(GIT_COMMIT)'

# 输出目录
OUTPUT_DIR := ./dist

.PHONY: all clean windows linux darwin help dev

# 默认目标
all: clean windows linux darwin

# 开发构建（当前平台）
dev:
	@echo "构建开发版本..."
	@go build -ldflags="$(LDFLAGS)" -o $(PROJECT_NAME) .

# 清理输出目录
clean:
	@echo "清理输出目录..."
	@rm -rf $(OUTPUT_DIR)
	@mkdir -p $(OUTPUT_DIR)

# 构建 Windows 版本（无平台后缀）
windows:
	@echo "构建 Windows 版本..."
	@GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(PROJECT_NAME).exe .

# 构建 Linux 版本
linux:
	@echo "构建 Linux 版本..."
	@GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(PROJECT_NAME)-linux-amd64 .

# 构建 macOS 版本
darwin:
	@echo "构建 macOS 版本..."
	@GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(PROJECT_NAME)-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(PROJECT_NAME) .

# 显示帮助信息
help:
	@echo "可用的构建目标："
	@echo "  all          - 构建所有平台版本"
	@echo "  dev          - 构建开发版本（当前平台）"
	@echo "  windows      - 构建 Windows 版本"
	@echo "  linux        - 构建 Linux 版本"
	@echo "  darwin       - 构建 macOS 版本"
	@echo "  clean        - 清理输出目录"
	@echo "  help         - 显示此帮助信息"
