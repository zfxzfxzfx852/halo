BINARY := vigil
PKG := ./cmd/vigil
GOFLAGS := -trimpath

.PHONY: all build run test fmt vet lint clean ci web-install web-dev web-build tidy

all: build

## build: 编译后端二进制
build:
	go build $(GOFLAGS) -o bin/$(BINARY) $(PKG)

## run: 直接运行
run:
	go run $(PKG)

## test: 运行全部测试
test:
	go test -race -cover ./...

## fmt: 格式化代码
fmt:
	go fmt ./...

## vet: 静态检查
vet:
	go vet ./...

## tidy: 整理依赖
tidy:
	go mod tidy

## web-install: 安装前端依赖
web-install:
	cd web && pnpm install

## web-dev: 前端开发服务器
web-dev:
	cd web && pnpm dev

## web-build: 构建前端
web-build:
	cd web && pnpm build

## ci: 流水线本地校验
ci: vet test build

## clean: 清理产物
clean:
	rm -rf bin/ web/dist/
