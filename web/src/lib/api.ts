import type { ApprovalRequest, HaloEvent } from '../types'

export interface EventQuery {
  type?: string
  protocol?: string
  session?: string
  riskMin?: string
  limit?: number
  offset?: number
}

export async function fetchEvents(query: EventQuery): Promise<{
  items: HaloEvent[]
  total: number
}> {
  const params = new URLSearchParams()
  Object.entries(query).forEach(([k, v]) => {
    if (v) params.set(k, String(v))
  })
  const res = await fetch(`/api/v1/events?${params.toString()}`)
  if (!res.ok) throw new Error(`fetch events failed: ${res.status}`)
  return res.json()
}

export async function fetchApprovals(): Promise<{
  items: ApprovalRequest[]
  total: number
}> {
  const res = await fetch('/api/v1/approvals')
  if (!res.ok) throw new Error(`fetch approvals failed: ${res.status}`)
  return res.json()
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
