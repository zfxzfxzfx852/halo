package model

// Protocol 调用协议类型。
type Protocol string

const (
	ProtocolHTTP      Protocol = "http"
	ProtocolMCPStdio  Protocol = "mcp+stdio"
	ProtocolMCPSSE    Protocol = "mcp+sse"
	ProtocolMCPHTTP   Protocol = "mcp+http"
)

// Decision 策略决策结果。
type Decision string

const (
	DecisionAllow Decision = "allow"
	DecisionDeny  Decision = "deny"
	DecisionAsk   Decision = "ask"
)

// RiskLevel 风险等级。
type RiskLevel string

const (
	RiskNone     RiskLevel = "none"
	RiskLow      RiskLevel = "low"
	RiskMedium   RiskLevel = "medium"
	RiskHigh     RiskLevel = "high"
	RiskCritical RiskLevel = "critical"
)

// EventType 事件类型。
type EventType string

const (
	EventRequest    EventType = "request"
	EventResponse   EventType = "response"
	EventToolCall   EventType = "tool_call"
	EventToolResult EventType = "tool_result"
	EventDecision   EventType = "decision"
	EventAlert      EventType = "alert"
)

// Severity 返回风险等级的数值权重，便于比较与排序。
func (r RiskLevel) Severity() int {
	switch r {
	case RiskCritical:
		return 4
	case RiskHigh:
		return 3
	case RiskMedium:
		return 2
	case RiskLow:
		return 1
	default:
		return 0
	}
}
