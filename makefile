
GO ?= go
# GOFMT ?= gofumpt "-s"
GOFMT ?= gofumpt
GOIMPORTS := goimports
# gofumpt 安装路径
GOFUMPT := gofumpt
GOIMPORTS-REVISER := goimports-reviser
# 工具命令（建议已全局安装）
GCI := gci

# GOFILES := $(shell find . -name "*.go")
GOFILES := .
# LDFLAGS := -s -w


.PHONY: install-tools
install-tools: # Install the necessary tools | 安装必要的工具
	@echo "==> 安装 golangci-lint,goimports,gofumpt,swagger,modd 等工具"
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest;
	$(GO) install golang.org/x/tools/cmd/goimports@latest
	$(GO) install github.com/incu6us/goimports-reviser/v3@latest
	$(GO) install mvdan.cc/gofumpt@latest;
	$(GO) install github.com/go-swagger/go-swagger/cmd/swagger@latest;
	$(GO) install github.com/cortesi/modd/cmd/modd@latest


.PHONY: fmt gofumpt check-fmt

fmt: # Format the codes | 格式化代码
	@echo " \033[34m ==> 自动整理 imports... <== \033[0m "
	$(GOIMPORTS) -w $(GOFILES)
	$(GOIMPORTS-REVISER) -rm-unused -format $(GOFILES)
	@echo " \033[34m ==> 格式化代码...  <==\033[0m"
	$(GOFMT) -w $(GOFILES)

# 使用 gci 进行 import 分组格式化（需要配置 prefix）
gci:
	@echo "==> 使用 gci 排序 imports..."
	@$(GCI) write --skip-generated -s standard -s default -s "prefix($(shell go list -m))" .


# 使用 gofumpt 检查代码是否已格式化（CI 场景使用）
check-fmt:
	@echo " \033[34m ==> 检查是否通过 gofumpt 格式... <== \033[0m "
	$(GOFUMPT) -l $(GOFILES) | tee /dev/stderr | test -z

# 检查 goimports 格式是否通过（CI场景可用）
check-imports:
	@echo "==> 检查 import 是否正确..."
	@! $(GOIMPORTS) -l $(GOFILES) | grep .


.PHONY: lint
lint: # Run go linter | 运行代码错误分析
	@echo " \033[34m ==> 运行代码错误分析 <== \033[0m "
	# golangci-lint run -D staticcheck
	golangci-lint run --disable=unused -D staticcheck

# 执行静态检查（例如 go vet）
lint-vet:
	@echo " \033[34m ==> 静态检查代码 (go vet)... <== \033[0m "
	@go vet ./...


.PHONY: Docker_Mac_Env
Docker_Mac_Env:
	docker-compose -f docker-compose-env.yml -p tony-go-zero-looklook  up -d

.PHONY: Docker_Mac_Start
Docker_Mac_Start:
	docker-compose up -d

.PHONY: Modd_Local
Modd_Local:
	modd -f ./modd/modd.local.conf

.PHONY: Modd_Dev
Modd_Dev:
	modd -f ./modd/modd.dev.conf





