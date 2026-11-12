package mcp

import "encoding/json"

// JSONRPC 版本标识。
const JSONRPCVersion = "2.0"

// MCP 常用 JSON-RPC 方法。
const (
	MethodInitialized = "notifications/initialized"
	MethodToolsList   = "tools/list"
	MethodToolsCall   = "tools/call"
)

// Request 表示 JSON-RPC 请求或通知。
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// IsNotification 通知没有 ID，无需响应。
func (r *Request) IsNotification() bool {
	return len(r.ID) == 0
}

// RPCError 是 JSON-RPC 错误对象。
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Response 表示 JSON-RPC 响应。
type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

// ToolCallParams 是 tools/call 的参数。
type ToolCallParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments,omitempty"`
}
