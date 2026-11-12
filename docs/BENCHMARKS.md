# 性能基准 / Benchmarks

## 中文

Halo 的代理在**请求处理路径**上引入的开销很小。以下数据来自 handler 层基准（上游为本地 `httptest`，不含真实网络往返），衡量的是 Halo 自身在转发、事件记录上的成本。

### 环境

- Go 1.23，linux/amd64
- CPU：AMD EPYC 9Y24（4 核参与基准）

### 结果

| 指标 | 数值 |
|---|---|
| 每请求耗时 | **40.04 µs/op** |
| 内存分配 | 14,771 B/op |
| 分配次数 | 124 allocs/op |

### 解读

- 一次典型 LLM API 调用的端到端延迟通常在数百毫秒到数秒之间，Halo 增加的约 40 µs 不到其万分之一，对实际使用基本无感。
- 代理为纯内存事件处理，无外部依赖，不发起额外网络请求。
- 流式（SSE）响应采用透传，不缓冲完整响应，长连接吞吐不受影响。

### 复现

```bash
go test -run=^$ -bench=. -benchmem ./internal/proxy/
```

---

## English

The proxy overhead added by Halo on the request path is small. The figures below come from a handler-level benchmark (local `httptest` upstream, no real network round-trip), measuring Halo's own cost in forwarding and event recording.

### Environment

- Go 1.23, linux/amd64
- CPU: AMD EPYC 9Y24 (4 benchmark cores)

### Results

| Metric | Value |
|---|---|
| Time per request | **40.04 µs/op** |
| Memory | 14,771 B/op |
| Allocations | 124 allocs/op |

### Notes

A typical LLM API call takes hundreds of milliseconds to seconds end-to-end, so the ~40 µs added by Halo is negligible. Streaming (SSE) responses are proxied through without buffering.

### Reproduce

```bash
go test -run=^$ -bench=. -benchmem ./internal/proxy/
```
