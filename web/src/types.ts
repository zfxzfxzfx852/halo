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
  risk: RiskLevel
  findings?: Finding[]
  parentId?: string
}
