import type { RiskLevel } from '../types'

export function riskClass(risk: RiskLevel): string {
  if (!risk) return 'risk-none'
  return `risk-${risk}`
}

export function formatTime(iso: string): string {
  const d = new Date(iso)
  return Number.isNaN(d.getTime()) ? iso : d.toLocaleTimeString()
}
