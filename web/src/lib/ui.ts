import type { EventType, HaloEvent, RiskLevel } from '../types'

export function riskClass(risk: RiskLevel | undefined): string {
  if (!risk || risk === 'none') return 'risk-none'
  return `risk-${risk}`
}

export const RISK_LABELS: Record<string, string> = {
  none: '安全',
  low: '低危',
  medium: '中危',
  high: '高危',
  critical: '严重',
}

export function riskLabel(risk: RiskLevel | undefined): string {
  if (!risk) return '安全'
  return RISK_LABELS[risk] || risk
}

export const TYPE_LABELS: Record<string, string> = {
  request: '请求',
  response: '响应',
  tool_call: '工具调用',
  tool_result: '工具结果',
  decision: '策略决策',
  alert: '告警',
}

export function typeLabel(t: EventType | string): string {
  return TYPE_LABELS[t] || t
}

export function formatTime(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return d.toLocaleTimeString('zh-CN', { hour12: false })
}

export function formatDateTime(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return `${d.toLocaleDateString('zh-CN')} ${d.toLocaleTimeString('zh-CN', { hour12: false })}`
}

export function fromNow(ts: number | string): string {
  const ms = typeof ts === 'number' ? ts / 1e6 : new Date(ts).getTime()
  const diff = Date.now() - ms
  if (Number.isNaN(diff)) return '—'
  if (diff < 60_000) return '刚刚'
  if (diff < 3_600_000) return `${Math.floor(diff / 60_000)} 分钟前`
  if (diff < 86_400_000) return `${Math.floor(diff / 3_600_000)} 小时前`
  return `${Math.floor(diff / 86_400_000)} 天前`
}

export function eventName(e: HaloEvent): string {
  if (e.toolCall?.name) return e.toolCall.name
  if (e.method) return e.method
  if (e.decision?.reason) return e.decision.reason
  return typeLabel(e.type)
}

export function eventEndpoint(e: HaloEvent): string {
  if (e.toolCall) return e.toolCall.server || 'mcp server'
  return e.protocol
}

export function severity(risk: RiskLevel | undefined): number {
  switch (risk) {
    case 'critical':
      return 4
    case 'high':
      return 3
    case 'medium':
      return 2
    case 'low':
      return 1
    default:
      return 0
  }
}

export function formatNumber(n: number | undefined): string {
  if (!n) return '0'
  if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)}M`
  if (n >= 10_000) return `${(n / 10_000).toFixed(1)}w`
  return n.toLocaleString('zh-CN')
}

export function decisionClass(d?: string): string {
  if (d === 'deny') return 'risk-deny'
  if (d === 'ask') return 'risk-ask'
  return 'risk-allow'
}

export function decisionLabel(d?: string): string {
  switch (d) {
    case 'allow':
      return '放行'
    case 'deny':
      return '拦截'
    case 'ask':
      return '待审批'
    default:
      return '未知'
  }
}
