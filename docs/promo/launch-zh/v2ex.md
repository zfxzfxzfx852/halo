# V2EX — 分享帖

## 节点

`/go/programmer`（也可考虑 `/go/opensource`、`/go/security`），按节点规则只发一次。

## 标题（选一个）

**首选**
```
开源了 Halo：本地优先的 AI Agent 可观测与安全控制平面
```

**备选**
```
做了个给 AI Agent 用的「Wireshark」，能看调用也能拦危险操作
```

## 正文

```
最近用各种 coding agent，发现一个问题：它们自己在调工具、跑命令、读写文件，
但我对具体发生了什么几乎没什么可见性，也没法在危险操作发生前拦住。现有的
一些 agent 安全工具要么是命令行日志，要么偏 B 端网关，实时可视化做得都不太好。

于是自己写了 Halo，MIT 开源，本地优先，定位是所有 agent 都能依赖的可观测
和安全基础设施，而不是再做一个 agent 框架。

它透明地拦在 agent 和模型/工具之间，主要能力：
- 透明代理 MCP（stdio/SSE/HTTP）和 OpenAI 兼容接口，不用改 agent 代码
- 实时事件流 + 可交互调用图谱，风险用颜色区分，有点像 Wireshark
- 声明式 allow / deny / 人工审批策略，支持热加载
- 离线检测提示注入、越狱、PII 和密钥，自动脱敏，不依赖外部服务
- 高危命令在沙箱里跑，预留 Docker/gVisor/Firecracker 驱动
- 防篡改哈希链审计，能完整回放会话

理念是「先观察，后管控」：默认 observe 模式只记录不阻断，看清正常行为、
建立基线之后，再切到 guard 模式上策略，接入过程没有打断工作流的风险。

工程上是单个静态编译的 Go 二进制（内嵌 React 控制台、纯 Go 的 SQLite，
无 CGO），数据都在本机，默认零遥测，也能配合 Ollama 这类本地模型。
一次策略决策大概 40 微秒，核心包覆盖率 70%+，CI 开了 race 检测。

macOS / Linux 一行装：
curl -fsSL https://raw.githubusercontent.com/zfxzfxzfx852/halo/main/scripts/install.sh | sh

仓库：https://github.com/zfxzfxzfx852/halo
README 里有演示 GIF 和截图。

欢迎拍砖，尤其是 MCP 桥接、检测规则和策略模型上的边界情况，也欢迎 PR。
```

## 发帖注意

- 时间：工作日 9:00、12:30、20:00 左右。
- 用真实账号，认真回复每一层楼；不要在多个节点重复发同一个链接。
- 被指出问题就承认并修，不要硬辩；不要刷回复顶帖。
