
GO ?= go
GOFMT ?= gofumpt "-s"
GOIMPORTS := goimports

# gofumpt 安装路径
GOFUMPT := gofumpt
# 所有 Go 文件
# GOFILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

GOFILES := $(shell find . -name "*.go")

# LDFLAGS := -s -w


.PHONY: install-tools
install-tools: # Install the necessary tools | 安装必要的工具
	@echo "==> 安装 golangci-lint,goimports,gofumpt,swagger,modd 等工具"
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest;
	$(GO) install golang.org/x/tools/cmd/goimports@latest
	$(GO) install mvdan.cc/gofumpt@latest;
	$(GO) install github.com/go-swagger/go-swagger/cmd/swagger@latest;
	$(GO) install github.com/cortesi/modd/cmd/modd@latest


.PHONY: fmt gofumpt check-fmt

fmt: # Format the codes | 格式化代码
	@echo "==> 使用 goimports 自动整理 imports..."
	$(GOIMPORTS) -w .
	@echo "==> 使用 gofumpt 格式化代码..."
	$(GOFMT) -w $(GOFILES)

# 使用 gofumpt 检查代码是否已格式化（CI 场景使用）
check-fmt:
	@echo "==> 检查是否通过 gofumpt 格式..."
	$(GOFUMPT) -l $(GOFILES) | tee /dev/stderr | test -z

# .PHONY: lint
# lint: # Run go linter | 运行代码错误分析
# 	# golangci-lint run -D staticcheck
# 	golangci-lint run --disable=unused -D staticcheck

.PHONY: Docker_Mac_Env
Docker_Mac_Env:
	docker-compose -f docker-compose-env-mac.yml -p tony-go-zero-looklook  up -d

.PHONY: Docker_Mac_Start
Docker_Mac_Start:
	docker-compose up -d

.PHONY: Modd_Local
Modd_Local:
	modd -f ./modd/modd.local.conf

.PHONY: Modd_Dev
Modd_Dev:
	modd -f ./modd/modd.dev.conf



# ssssssssss

# 项目根路径（自动获取）
ROOT_DIR := $(shell pwd)

# Go 命令配置
GOFILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

# 工具命令（建议已全局安装）
GOIMPORTS := goimports
GCI := gci

.PHONY: all fmt lint format-imports check-imports install-tools

# 默认任务：格式化 + 导入 + 静态检查
all: fmt format-imports lint

# 使用 gofmt 格式化代码
# fmt:
# 	@echo "==> 格式化代码 (gofmt)..."
# 	@gofumpt -s -w .

# 使用 goimports 自动整理 imports（分组、删除未使用的）
format-imports:
	@echo "==> 使用 goimports 自动整理 imports..."
	@$(GOIMPORTS) -w .

# 使用 gci 进行 import 分组格式化（需要配置 prefix）
gci:
	@echo "==> 使用 gci 排序 imports..."
	@$(GCI) write --skip-generated -s standard -s default -s "prefix($(shell go list -m))" .

# 执行静态检查（例如 go vet）
lint:
	@echo "==> 静态检查代码 (go vet)..."
	@go vet ./...

# 检查 goimports 格式是否通过（CI场景可用）
check-imports:
	@echo "==> 检查 import 是否正确..."
	@! $(GOIMPORTS) -l $(GOFILES) | grep .
