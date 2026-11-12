# 接入通用 MCP 客户端

Halo 可以包装任意基于 **stdio** 的 MCP server：MCP 客户端连接 Halo，Halo 再启动并桥接真实的 MCP server，全程记录工具调用。

## 原理

```
MCP 客户端 ──stdio──► halo mcp ──stdio──► 真实 MCP server
                        │
                        └──► 工具调用事件 / 检测 / 审计（同时提供 HTTP 控制台）
```

## 用法

把 MCP 客户端原本的启动命令：

```
npx -y @modelcontextprotocol/server-filesystem /tmp
```

包装成：

```sh
halo mcp --config halo.yaml -- npx -y @modelcontextprotocol/server-filesystem /tmp
```

`--` 之后是真实 MCP server 的完整启动命令，原样填写即可。

## 配置示例（MCP 客户端）

在支持 `mcpServers` 配置的客户端（如 Claude Desktop、支持 MCP 的编辑器）中：

```json
{
  "mcpServers": {
    "filesystem": {
      "command": "halo",
      "args": [
        "mcp",
        "--config",
        "/path/to/halo.yaml",
        "--",
        "npx",
        "-y",
        "@modelcontextprotocol/server-filesystem",
        "/tmp"
      ]
    }
  }
}
```

## 查看数据

桥接进程会同时启动 HTTP 服务，打开 `http://127.0.0.1:8787/console/`：

- 「事件流」中协议显示为 `mcp+stdio`
- 「调用图谱」展示每次 `tools/call` 及其结果
- 所有调用同样写入防篡改审计链

## 说明

- MCP 的 stdio 通道承载 JSON-RPC，Halo 的运行日志输出到 **stderr**，不会干扰协议。
- MCP over HTTP（SSE/Streamable HTTP）的调用，通过普通 HTTP 代理方式接入即可，见 [接入总览](./README.md)。
