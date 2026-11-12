# Halo 开发路线图（16 周 · 从 0 到 1）

> Halo —— The local-first observability & security plane for AI agents.
> 给任何 AI Agent 装上「飞行记录仪 + 安全气囊 + 控制台」。

- 计划周期：**2026-11-09（W1）→ 2027-02-28（W16）**
- 北极星指标：4 个月内 GitHub **冲刺 40k star**（保底 6k–12k，现实 15k–25k，冲刺 35k–45k）
- 备份策略：每个可运行节点打 Git tag，远程仓库 + 本地镜像双备份（详见第 6 节）

---

## 1. 目标与成功定义

| 维度 | 成功标准 |
|---|---|
| 产品 | 单二进制/一行命令启动，零配置拦截本机 Agent，Observe→Guard 双模式 |
| 兼容 | 首发打透 Claude Code、DeepSeek Harness、Ollama 三个对象 |
| 技术 | 代理延迟开销可量化（目标 P99 追加延迟 < 5ms，不含上游） |
| 传播 | 至少 1 次 Hacker News 首页、1 次 GitHub Trending、国内 3 个平台热榜 |
| 工程 | 测试覆盖核心包 ≥ 70%，CI 全绿，每个里程碑可回滚 |

## 2. 技术栈与总体架构

- **数据面/后端**：Go 1.23（单二进制、并发、低延迟代理）
- **控制台前端**：React 18 + TypeScript + Vite，构建产物用 `embed` 打进二进制
- **存储**：SQLite（嵌入式，零外部依赖），事件流 + 审计哈希链
- **实时通信**：WebSocket
- **沙箱**：先做 `Sandbox` 接口 + 进程级隔离落地；Docker/gVisor/Firecracker/WASM 作为后续驱动（当前环境无 Docker，接口先行）
- **检测**：规则/AST/启发式 + 本地 LLM（Ollama）语义检测

```
Agent ──► Halo 透明代理(拦截) ──► 策略引擎 ──► 检测引擎 ──► 沙箱/放行 ──► 模型/工具
              │                        │            │
              ▼                        ▼            ▼
          事件采集 ──► SQLite 存储 ──► 审计哈希链 ──► WebSocket ──► React 控制台
```

模块职责：

| 包 | 职责 |
|---|---|
| `internal/proxy` | MCP（stdio/SSE/streamable HTTP）与 OpenAI 兼容 API 透明代理 |
| `internal/policy` | 声明式策略解析与决策（allow/deny/ask） |
| `internal/detect` | prompt injection、jailbreak、PII、密钥检测 |
| `internal/sandbox` | 高危命令/代码隔离执行（驱动可插拔） |
| `internal/audit` | 防篡改哈希链日志、会话录制 |
| `internal/store` | SQLite 持久化与查询 |
| `internal/server` | HTTP API、WebSocket、静态资源服务 |
| `internal/config` | 配置加载与热更新 |
| `internal/model` | 领域模型与事件结构 |

## 3. 逐周计划

> 每周固定节奏：周一拆任务 → 周三中期检查 → 周五提交并打周备份 tag → 写周报。

### M0 地基（W1–W2）

#### W1（11/09–11/15）立项与骨架
- 目标：可运行的最小骨架 + 完整规划
- 任务：
  - 初始化 Go module、目录结构、Makefile、`.gitignore`、LICENSE(MIT)
  - 落地 `ROADMAP.md`、`GITHUB_PLAN.md`、README 双语骨架
  - `cmd/halo` 最小入口：加载配置、启动 HTTP server、`/healthz`
  - git init，建立提交规范与周备份 tag 机制
- 交付物：`halo` 可编译运行，`GET /healthz` 返回 200
- 验收：`make build && ./halo` 正常；`make test` 通过
- 备份点：tag `v0.0.1-week1`

#### W2（11/16–11/22）配置与领域模型
- 目标：配置体系 + 事件模型定型
- 任务：
  - `internal/config`：YAML 配置 + 环境变量覆盖 + 默认值
  - `internal/model`：Event、ToolCall、Policy、Decision、AuditRecord 结构
  - 定义统一事件总线接口（内部发布/订阅）
  - 结构化日志（slog）与版本信息注入
- 交付物：配置可加载、事件模型单测通过
- 验收：`make test`，核心模型测试全绿
- 备份点：tag `v0.0.2-week2`

### M1 Observe（W3–W6）

#### W3（11/23–11/29）HTTP 透明代理
- 目标：拦截 OpenAI 兼容 HTTP 调用
- 任务：
  - 反向代理：请求/响应缓冲、Header 透传、流式(SSE) 转发
  - 解析 `chat/completions`、`tools/call` 载荷，生成原始事件
  - 代理吞吐与延迟基准（benchmark）
- 交付物：指向 Halo 的客户端请求可被完整记录并正确转发
- 验收：流式响应不中断；基准测试输出延迟数据
- 备份点：tag `v0.1.0-http-proxy`

#### W4（11/30–12/06）MCP 拦截
- 目标：支持 MCP stdio 与 streamable HTTP
- 任务：
  - stdio MCP 包装器（子进程 stdin/stdout 劫持）
  - MCP over HTTP/SSE 代理，JSON-RPC 解析
  - 工具清单（tools/list）与调用（tools/call）事件化
- 交付物：MCP 工具调用可被拦截、记录、转发
- 验收：用一个示例 MCP server 端到端打通
- 备份点：tag `v0.2.0-mcp`

#### W5（12/07–12/13）存储与实时推送
- 目标：事件持久化 + WebSocket 实时流
- 任务：
  - `internal/store`：SQLite 建表、写入、分页查询、保留期清理
  - `internal/server`：WebSocket 广播、事件订阅
  - REST API：事件列表、详情、过滤、统计
- 交付物：事件可落库并通过 WS 实时推送
- 验收：断线重连可补拉；查询接口单测通过
- 备份点：tag `v0.3.0-store-ws`

#### W6（12/14–12/20）Observe 控制台
- 目标：可视化实时调用流
- 任务：
  - Vite + React 工程搭建，接入 WS
  - 实时事件流列表、详情抽屉、过滤搜索
  - 调用图谱基础版（节点=模型/工具，边=调用，风险着色）
  - 前端构建产物 `embed` 进二进制
- 交付物：浏览器打开控制台看到实时调用流与图谱
- 验收：单二进制启动即可访问 UI；录制第一条演示 GIF
- 备份点：tag `v0.4.0-observe-ui`

### M2 Guard（W7–W10）

#### W7（12/21–12/27）策略引擎
- 目标：声明式策略可决策
- 任务：
  - 策略 YAML 结构与加载（工具/命令/路径/参数规则）
  - 决策器：allow / deny / ask，优先级与默认策略
  - 策略热更新，单元测试覆盖各分支
- 交付物：策略命中后可拦截或放行调用
- 验收：给定测试策略集，决策结果 100% 符合预期
- 备份点：tag `v0.5.0-policy`

#### W8（12/28–01/03）检测引擎
- 目标：识别注入与敏感信息
- 任务：
  - 规则/正则/AST 检测 prompt injection、jailbreak
  - PII 与密钥（API Key/Token）识别与脱敏
  - Ollama 语义检测适配（可选开关，默认离线规则）
  - 误报/漏报测试集与评分
- 交付物：风险事件带风险等级与命中规则
- 验收：测试集召回率/误报率输出报告
- 备份点：tag `v0.6.0-detect`

#### W9（01/04–01/10）人工审批与沙箱
- 目标：高危调用可审批、可隔离
- 任务：
  - 审批流：WS 推送审批请求，控制台一键允许/拒绝，调用挂起/恢复
  - `Sandbox` 接口 + 进程级隔离驱动（权限降级、超时、资源限制）
  - 预留 Docker/gVisor/Firecracker/WASM 驱动注册点
- 交付物：高危命令在审批后于隔离环境执行
- 验收：未审批不执行；越权/超时被正确终止
- 备份点：tag `v0.7.0-approval-sandbox`

#### W10（01/11–01/17）审计与回放
- 目标：防篡改审计链
- 任务：
  - 哈希链：每条记录含前一条哈希，可校验完整性
  - 会话录制与回放（请求/响应/决策时间线）
  - 审计导出（JSONL）与校验工具
- 交付物：任意篡改可被检测，会话可回放
- 验收：篡改一条记录即校验失败
- 备份点：tag `v0.8.0-audit`

### M3 打磨（W11–W12）

#### W11（01/18–01/24）安装与适配
- 任务：一键安装脚本（curl/brew）、三大对象接入指南、配置向导、Dockerfile 与 compose（补齐容器化）
- 验收：干净环境一条命令安装并跑通
- 备份点：tag `v0.9.0-install`

#### W12（01/25–01/31）文档/i18n/CI/性能
- 任务：文档站、中英文 i18n、benchmark 报告、GitHub Actions（lint/test/build/coverage）、安全白皮书与威胁模型
- 验收：CI 全绿；核心包覆盖率 ≥ 70%；内测分发
- 备份点：tag `v1.0.0-rc1`

### M4 引爆（W13–W14）

#### W13（02/01–02/07）首发准备
- 任务：落地页、演示视频、首发物料（中英）、媒体/KOL 名单、内测反馈收口、发布候选 `v1.0.0-rc2`
- 验收：所有发布素材就位，发布时间表冻结

#### W14（02/08–02/14）全球首发
- 任务：Hacker News + Reddit + Product Hunt + X 首发；国内 V2EX/掘金/知乎/公众号/即刻/DeepSeek 社区；发布 `v1.0.0`
- 验收：完成首发清单，按小时监控 star/流量/issue（详见 GITHUB_PLAN）

### M5 二波（W15–W16）

#### W15（02/15–02/21）响应与迭代
- 任务：issue 当天响应、热点 bug 修复、发布 `v1.0.1`、补充社区呼声最高的适配器
- 验收：高优 issue 清零，发布节奏稳定

#### W16（02/22–02/28）二波传播
- 任务：企业/合规叙事文章、案例与用户故事、贡献者计划、`v1.1` 路线公布、复盘数据报告
- 验收：输出 16 周复盘，确定后续季度路线

## 4. 里程碑与版本映射

| 里程碑 | 周次 | 版本 | 主题 |
|---|---|---|---|
| M0 | W1–W2 | v0.0.x | 骨架与模型 |
| M1 | W3–W6 | v0.1–v0.4 | Observe 可用 |
| M2 | W7–W10 | v0.5–v0.8 | Guard 闭环 |
| M3 | W11–W12 | v0.9–rc1 | 打磨与 CI |
| M4 | W13–W14 | rc2 / v1.0 | 首发 |
| M5 | W15–W16 | v1.0.1 / v1.1 | 二波与维护 |

## 5. 风险与回滚

| 风险 | 影响 | 缓解 | 回滚 |
|---|---|---|---|
| 单人工作量超支 | 延期 | 严控 MVP，非核心功能后置 | 砍功能不砍里程碑，回退到上一稳定 tag |
| 代理拖慢/断流 | 不可用 | 基准测试 + 透传降级开关 | 关闭拦截，仅旁路观察 |
| 语义检测误报高 | 体验差 | 默认离线规则，语义检测可选 | 一键关闭语义检测 |
| 环境无 Docker | 沙箱受限 | 接口先行 + 进程隔离 | 沙箱驱动后续补，不阻塞主链 |
| 同类竞争 | 差异化弱 | 可视化 + 体验快速占位 | 聚焦更细分场景 |

## 6. 备份策略（强制）

1. **本地 Git**：每个功能小步提交，周五打周 tag；提交信息遵循 Conventional Commits。
2. **远程备份**：GitHub 私有仓库先同步，首发前转公开；另推一份镜像（Gitee/备用远端）。
3. **不可变里程碑**：每个 `v0.x`/`v1.x` tag 对应可编译产物，任何阶段失败可 `git checkout` 回滚。
4. **发布前校验**：每个 tag 必须通过 `make ci`（lint+test+build）。
