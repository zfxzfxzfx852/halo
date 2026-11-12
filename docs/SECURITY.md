# 安全白皮书与威胁模型 / Security & Threat Model

## 中文

### 1. 定位与信任边界

Halo 是一个**本地优先**的防御性控制平面，运行在用户自己的机器上，位于 AI Agent 与模型/工具之间。

```
[AI Agent] ──► [Halo（本地，信任边界）] ──► [模型 API / 工具 / MCP server]
```

- Halo 与 Agent 之间、Halo 与上游之间的所有流量都在本机或用户授权的网络内处理。
- 默认不向任何外部服务上报数据，不内置遥测，数据存储在本地 SQLite。
- 上游返回的内容（模型响应、工具结果、网页文本）一律视为**不可信输入**，只做解析与展示，不赋予其任何执行权限。

### 2. 要防御的威胁

| 威胁 | 机制 |
|---|---|
| 提示注入 / 越狱 | 离线规则检测，标记风险并可按策略阻断 |
| 密钥与 PII 泄露 | 密钥/个人信息识别与脱敏，出站前可拦截 |
| 危险工具调用 | 声明式策略 allow/deny，高危转人工审批 |
| 命令执行风险 | 沙箱隔离（独立进程组、权限降级，预留容器化驱动） |
| 事后抵赖 / 篡改 | 防篡改哈希链审计，任何改动可被校验发现 |
| 敏感头泄露 | 代理对 Authorization 等敏感头脱敏展示 |

### 3. 不承诺防御的范围

- Halo 不替代操作系统级安全边界；进程级沙箱是纵深防御的一层，强隔离需启用 Docker/gVisor 等驱动（路线图内）。
- 无法阻止用户**主动、知情地**批准一个本就危险的操作；审批是把决定权交给人，而非替人判断。
- 不防御已在 Halo 之外、未经过 Halo 的流量。要获得保护，所有调用必须经 Halo 代理。

### 4. 安全设计原则

1. 默认安全：默认 `observe` 模式不阻断，先让用户看清数据再切 `guard`。
2. 最小权限：本地监听默认绑定 `127.0.0.1`，不暴露到公网。
3. 零外部依赖：检测规则纯离线，不把数据发给第三方。
4. 可验证：审计链让每一次调用都可独立校验。

### 5. 漏洞披露

欢迎通过 GitHub 的 **Private vulnerability reporting** 私密报告安全问题。请不要在公开 issue 中披露未修复漏洞细节。我们会在确认后尽快响应并发布修复。

---

## English

### Scope & Trust Boundary

Halo is a local-first defensive control plane running on the user's own machine, between the AI agent and its models/tools. Traffic is processed locally, there is no built-in telemetry, and data stays in local SQLite. Upstream content (model responses, tool output, web text) is always treated as untrusted input and never granted execution authority.

### Threats Addressed

Prompt injection/jailbreak, secret and PII leakage, dangerous tool calls, risky command execution, tampering/repudiation, and sensitive-header exposure — addressed through offline detection, declarative policy, human approval, sandboxing, and a tamper-evident hash chain.

### Out of Scope

Halo does not replace OS-level isolation (strong isolation requires container drivers on the roadmap), cannot stop a user from knowingly approving a risky action, and does not protect traffic that bypasses the proxy.

### Reporting

Please report security issues privately via GitHub **Private vulnerability reporting** rather than public issues.
