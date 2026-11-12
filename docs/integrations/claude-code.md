# 接入 Claude Code

Claude Code 通过 `ANTHROPIC_BASE_URL` 环境变量支持自定义端点，把它指向 Halo 即可。

## 1. 配置上游

在 `halo.yaml` 中把上游设为 Anthropic 官方地址：

```yaml
mode: observe
upstream:
  baseUrl: "https://api.anthropic.com"
```

启动 Halo：

```sh
halo
```

## 2. 让 Claude Code 走 Halo

在启动 Claude Code 的同一个终端里设置环境变量：

```sh
export ANTHROPIC_BASE_URL=http://127.0.0.1:8787
claude
```

API Key 仍由 Claude Code 正常携带，Halo 会在转发时保留鉴权头，并对敏感头做脱敏展示。

## 3. 验证

在 Claude Code 中发起一次对话，然后打开 `http://127.0.0.1:8787/console/`，应能看到对应的请求与响应事件。

## 持久化（可选）

把 export 写进 shell 配置，让每次打开终端都生效：

```sh
echo 'export ANTHROPIC_BASE_URL=http://127.0.0.1:8787' >> ~/.zshrc
```

## 切回直连

```sh
unset ANTHROPIC_BASE_URL
```
