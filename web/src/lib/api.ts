import type { HaloEvent } from '../types'

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
