import type {
  AuditReport,
  ApprovalRequest,
  FindingStat,
  HaloEvent,
  PolicyRule,
  RuntimeInfo,
  SessionStat,
  Stats,
} from '../types'

export interface EventQuery {
  type?: string
  protocol?: string
  session?: string
  riskMin?: string
  limit?: number
  offset?: number
}

async function getJSON<T>(url: string): Promise<T> {
  const res = await fetch(url)
  if (!res.ok) throw new Error(`request failed: ${res.status}`)
  return (await res.json()) as T
}

export async function fetchEvents(query: EventQuery): Promise<{
  items: HaloEvent[]
  total: number
}> {
  const params = new URLSearchParams()
  Object.entries(query).forEach(([k, v]) => {
    if (v !== undefined && v !== '') params.set(k, String(v))
  })
  return getJSON(`/api/v1/events?${params.toString()}`)
}

export async function fetchEvent(id: string): Promise<HaloEvent> {
  return getJSON(`/api/v1/events/${encodeURIComponent(id)}`)
}

export async function fetchStats(): Promise<Stats> {
  return getJSON('/api/v1/stats')
}

export async function fetchApprovals(): Promise<{
  items: ApprovalRequest[]
  total: number
}> {
  return getJSON('/api/v1/approvals')
}

export async function resolveApproval(
  id: string,
  decision: 'allow' | 'deny',
): Promise<void> {
  const res = await fetch(`/api/v1/approvals/${id}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ decision }),
  })
  if (!res.ok) throw new Error(`resolve approval failed: ${res.status}`)
}

export async function fetchSessions(limit = 50): Promise<{
  items: SessionStat[]
  total: number
}> {
  return getJSON(`/api/v1/sessions?limit=${limit}`)
}

export async function fetchFindings(): Promise<{
  items: FindingStat[]
  total: number
}> {
  return getJSON('/api/v1/findings')
}

export async function fetchPolicies(): Promise<{
  default: string
  rules: PolicyRule[]
}> {
  return getJSON('/api/v1/policies')
}

export async function fetchRuntime(): Promise<RuntimeInfo> {
  return getJSON('/api/v1/runtime')
}

export async function verifyAudit(): Promise<AuditReport> {
  return getJSON('/api/v1/audit/verify')
}

// exportAudit 拉取 JSONL 并触发浏览器下载。
export async function downloadAudit(): Promise<number> {
  const res = await fetch('/api/v1/audit/export')
  if (!res.ok) throw new Error(`export failed: ${res.status}`)
  const text = await res.text()
  const lines = text.trim() ? text.trim().split('\n').length : 0
  const blob = new Blob([text], { type: 'application/x-ndjson' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `halo-audit-${new Date().toISOString().slice(0, 19).replace(/[:T]/g, '-')}.jsonl`
  document.body.appendChild(a)
  a.click()
  a.remove()
  URL.revokeObjectURL(url)
  return lines
}
