BINARY := halo
PKG := ./cmd/halo
GOFLAGS := -trimpath
VERSION ?= $(shell git describe --tags --always)
DIST_DIR := dist
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

.PHONY: all build build-go run test fmt vet lint clean ci web-install web-dev web-build tidy release checksums

all: build

## build: 先构建前端，再编译含嵌入式 UI 的完整二进制
build: web-build
	go build $(GOFLAGS) -o bin/$(BINARY) $(PKG)

## build-go: 仅编译后端（要求 web/dist 已存在）
build-go:
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

## release: 交叉编译全部平台并打包（含 LICENSE 与示例配置）
release: web-build
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; arch=$${platform#*/}; \
		name=$(BINARY)-$(VERSION)-$$os-$$arch; \
		echo "building $$name"; \
		GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build $(GOFLAGS) \
			-ldflags "-X main.version=$(VERSION)" -o $(DIST_DIR)/$$name/$(BINARY) $(PKG); \
		cp LICENSE halo.example.yaml $(DIST_DIR)/$$name/; \
		tar -C $(DIST_DIR) -czf $(DIST_DIR)/$$name.tar.gz $$name; \
		rm -rf $(DIST_DIR)/$$name; \
	done

## checksums: 生成发布包 SHA256 校验和
checksums:
	cd $(DIST_DIR) && sha256sum *.tar.gz > checksums.txt

## clean: 清理产物
clean:
	rm -rf bin/ web/dist/ $(DIST_DIR)
