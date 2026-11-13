export type RiskLevel = '' | 'none' | 'low' | 'medium' | 'high' | 'critical'

export type EventType =
  | 'request'
  | 'response'
  | 'tool_call'
  | 'tool_result'
  | 'decision'
  | 'alert'

export interface ToolCall {
  id: string
  name: string
  arguments: unknown
  server?: string
}

export interface Payload {
  headers?: Record<string, string>
  body?: unknown
  status?: number
}

export interface Finding {
  rule: string
  level: RiskLevel
  message: string
  field?: string
}

export interface DecisionRecord {
  decision: 'allow' | 'deny' | 'ask'
  reason: string
  policy?: string
  reviewedBy?: string
  reviewedAt?: string
}

// 后端 approval.Request 的 JSON 字段为大写开头。
export interface ApprovalRequest {
  ID: string
  EventID: string
  Summary: string
  Risk: RiskLevel
  CreatedAt: string
}

export interface HaloEvent {
  id: string
  type: EventType
  protocol: string
  agent: string
  sessionId: string
  timestamp: string
  method?: string
  request?: Payload
  response?: Payload
  toolCall?: ToolCall
  toolResult?: { content?: unknown; isError?: boolean; durationMs?: number }
  decision?: DecisionRecord
  risk: RiskLevel
  findings?: Finding[]
  parentId?: string
}

export interface Stats {
  Total: number
  ByType: Record<string, number>
  ByRisk: Record<string, number>
}

export interface SessionStat {
  sessionId: string
  agent: string
  events: number
  maxRisk: RiskLevel
  lastSeen: number
  toolCalls: number
}

export interface FindingStat {
  rule: string
  count: number
  maxRisk: RiskLevel
}

export interface PolicyRule {
  name: string
  tool?: string
  method?: string
  command?: string
  paths?: string[]
  args?: Record<string, string>
  decision: 'allow' | 'deny' | 'ask'
  maxRisk?: RiskLevel
}

export interface RuntimeInfo {
  version: string
  commit: string
  buildTime: string
  mode: string
  addr: string
  dataDir: string
  retentionDays: number
  policyFile: string
  upstream: string
  wsClients: number
}

export interface AuditReport {
  valid: boolean
  checked: number
  missingEvents: number
  brokenSeq?: number
  reason?: string
}
