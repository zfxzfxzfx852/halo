# ---- 前端构建 ----
FROM node:22-alpine AS web
WORKDIR /app/web
RUN corepack enable
COPY web/package.json web/pnpm-workspace.yaml web/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY web/ ./
RUN pnpm build

# ---- 后端构建（静态二进制）----
FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=web /app/web/dist ./web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath \
    -ldflags "-s -w -X main.version=docker" -o /out/halo ./cmd/halo

# ---- 运行镜像 ----
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -u 10001 halo
USER halo
COPY --from=build /out/halo /usr/local/bin/halo
EXPOSE 8787
ENTRYPOINT ["/usr/local/bin/halo"]
