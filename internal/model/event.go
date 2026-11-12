package model

import (
	"encoding/json"
	"time"
)

// ToolCall 描述一次工具调用。
type ToolCall struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
	Server    string          `json:"server,omitempty"` // MCP server 标识
}

// Payload 记录 HTTP/RPC 的原始载荷。
type Payload struct {
	Headers map[string]string `json:"headers,omitempty"`
	Body    json.RawMessage `json:"body,omitempty"`
	Status  int             `json:"status,omitempty"`
}

// Finding 表示一次检测命中。
type Finding struct {
	Rule    string    `json:"rule"`
	Level   RiskLevel `json:"level"`
	Message string    `json:"message"`
	Field   string    `json:"field,omitempty"`
}

// DecisionRecord 策略决策与审批信息。
type DecisionRecord struct {
	Decision   Decision   `json:"decision"`
	Policy     string     `json:"policy,omitempty"`
	Reason     string     `json:"reason"`
	ReviewedBy string     `json:"reviewedBy,omitempty"`
	ReviewedAt *time.Time `json:"reviewedAt,omitempty"`
}

// Event 是系统内统一流转的事件。
type Event struct {
	ID        string          `json:"id"`
	Type      EventType       `json:"type"`
	Protocol  Protocol        `json:"protocol"`
	Agent     string          `json:"agent"`
	SessionID string          `json:"sessionId"`
	Timestamp time.Time       `json:"timestamp"`
	Method    string          `json:"method,omitempty"` // JSON-RPC method 或 HTTP 路径
	Request   *Payload        `json:"request,omitempty"`
	Response  *Payload        `json:"response,omitempty"`
	ToolCall  *ToolCall       `json:"toolCall,omitempty"`
	Decision  *DecisionRecord `json:"decision,omitempty"`
	Risk      RiskLevel       `json:"risk"`
	Findings  []Finding       `json:"findings,omitempty"`
	ParentID  string          `json:"parentId,omitempty"`
}

// AuditRecord 防篡改审计哈希链中的一条记录。
type AuditRecord struct {
	Seq       int64     `json:"seq"`
	EventID   string    `json:"eventId"`
	Hash      string    `json:"hash"`
	PrevHash  string    `json:"prevHash"`
	Timestamp time.Time `json:"timestamp"`
}
