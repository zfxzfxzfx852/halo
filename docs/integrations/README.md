# 接入指南

Halo 是一个**本地透明代理**，默认监听 `http://127.0.0.1:8787`。接入的核心思路只有一步：

> 把工具原本指向模型厂商的 **Base URL** 改成 Halo 的地址，Halo 会把请求转发到你配置的真实上游（`upstream.baseUrl`），同时完成记录、检测与管控。

## 工作原理

```
AI 工具/Agent ──► Halo (127.0.0.1:8787) ──► 真实模型 API (upstream.baseUrl)
                        │
                        └──► 事件记录 / 风险检测 / 策略管控 / 审计
```

## 两种模式

| 模式 | 行为 | 适用 |
|---|---|---|
| `observe` | 只观察、记录、检测，不阻断任何请求 | 首次接入，零风险先看数据 |
| `guard` | 按策略允许 / 拒绝 / 转人工审批 | 确认规则后启用管控 |

## 通用步骤

1. 安装并启动 Halo：
   ```sh
   halo            # 或 halo init 生成配置后再启动
   ```
2. 在配置文件 `halo.yaml` 中设置真实上游：
   ```yaml
   upstream:
     baseUrl: "https://api.openai.com"
   ```
3. 把 AI 工具的 Base URL 改为 `http://127.0.0.1:8787`（具体方式见各工具文档）。
4. 打开控制台 `http://127.0.0.1:8787/console/` 查看实时调用。

## 按工具查看

- [Claude Code](./claude-code.md)
- [Cursor](./cursor.md)
- [Cline](./cline.md)
- [通用 MCP 客户端](./mcp.md)
