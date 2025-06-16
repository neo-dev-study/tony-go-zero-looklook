GO ?= go
GOFMT ?= gofumpt "-s"
GOFILES := $(shell find . -name "*.go")
LDFLAGS := -s -w


.PHONY: tools
tools: # Install the necessary tools | 安装必要的工具
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest;
	$(GO) install golang.org/x/tools/cmd/goimports@latest
	$(GO) install mvdan.cc/gofumpt@latest;
	$(GO) install github.com/go-swagger/go-swagger/cmd/swagger@latest;
	$(GO) install github.com/cortesi/modd/cmd/modd@latest


.PHONY: fmt
fmt: # Format the codes | 格式化代码
	goimports -w .
	$(GOFMT) -w $(GOFILES)


.PHONY: lint
lint: # Run go linter | 运行代码错误分析
	# golangci-lint run -D staticcheck
	golangci-lint run --disable=unused -D staticcheck

.PHONY: Docker_Mac_Env
Docker_Mac_Env:
	docker-compose -f docker-compose-env-mac.yml up -d

.PHONY: Docker_Mac_Start
Docker_Mac_Start:
	docker-compose up -d

.PHONY: Modd_Local
Modd_Local:
	modd -f ./modd/local.modd.conf

.PHONY: Modd_Dev
Modd_Dev:
	modd -f ./modd/dev.modd.conf