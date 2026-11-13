import { useMemo, useState } from 'react'
import type { EventType, HaloEvent, RiskLevel } from '../types'
import { Icon } from './Icon'
import {
  decisionClass,
  decisionLabel,
  eventName,
  formatTime,
  riskClass,
  riskLabel,
  severity,
  typeLabel,
} from '../lib/ui'

interface Props {
  events: HaloEvent[]
  selectedId?: string
  onSelect: (id: string) => void
  live: boolean
  onToggleLive: () => void
  sessionFilter?: string
  onClearSession: () => void
  total: number
}

const TYPES: { key: 'all' | EventType; label: string }[] = [
  { key: 'all', label: '全部' },
  { key: 'request', label: '请求' },
  { key: 'response', label: '响应' },
  { key: 'tool_call', label: '工具调用' },
  { key: 'tool_result', label: '工具结果' },
  { key: 'decision', label: '决策' },
]

const RISKS = [
  { key: 'all', label: '全部风险' },
  { key: 'low', label: '低危以上' },
  { key: 'medium', label: '中危以上' },
  { key: 'high', label: '高危以上' },
]

export function EventList({
  events,
  selectedId,
  onSelect,
  live,
  onToggleLive,
  sessionFilter,
  onClearSession,
  total,
}: Props) {
  const [type, setType] = useState<'all' | EventType>('all')
  const [risk, setRisk] = useState('all')

  const filtered = useMemo(() => {
    const min = risk === 'all' ? -1 : severity(risk as RiskLevel)
    return events
      .filter((e) => (type === 'all' ? true : e.type === type))
      .filter((e) => (min < 0 ? true : severity(e.risk) >= min))
      .filter((e) => (sessionFilter ? e.sessionId === sessionFilter : true))
      .slice()
  }, [events, type, risk, sessionFilter])

  const shown = filtered.slice(0, 200)

  return (
    <div className="events-view">
      <section className="card glass events-card">
        <header className="events-head">
          <div className="events-title">
            <h2>
              <Icon name="pulse" size={18} />
              实时事件流
            </h2>
            <span className="subtle">{total} 条已落库 · 显示最近 {shown.length} 条</span>
          </div>
          <div className="events-tools">
            {sessionFilter && (
              <button className="chip active" onClick={onClearSession}>
                会话 {sessionFilter} ×
              </button>
            )}
            <button className={`live-btn ${live ? 'on' : ''}`} onClick={onToggleLive}>
              <Icon name={live ? 'pause' : 'play'} size={15} />
              {live ? '实时' : '已暂停'}
            </button>
          </div>
        </header>

        <div className="filters">
          <div className="chips">
            {TYPES.map((t) => (
              <button
                key={t.key}
                className={`chip ${type === t.key ? 'active' : ''}`}
                onClick={() => setType(t.key)}
              >
                {t.label}
              </button>
            ))}
          </div>
          <div className="chips">
            {RISKS.map((r) => (
              <button
                key={r.key}
                className={`chip ${risk === r.key ? 'active' : ''}`}
                onClick={() => setRisk(r.key)}
              >
                {r.label}
              </button>
            ))}
          </div>
        </div>

        {shown.length === 0 ? (
          <div className="empty-state">
            <div className="empty-orb">
              <Icon name="pulse" size={26} />
            </div>
            <h3>当前筛选下没有事件</h3>
            <p>把客户端指向 {location.host} 并发起一次请求，这里会实时出现。</p>
          </div>
        ) : (
          <div className="event-table">
            <div className="event-table-head">
              <span>类型</span>
              <span>调用</span>
              <span>会话</span>
              <span>风险</span>
              <span>时间</span>
            </div>
            {shown.map((e) => {
              const isDecision = e.type === 'decision'
              const badgeClass = isDecision
                ? decisionClass(e.decision?.decision)
                : riskClass(e.risk)
              const rowClass = isDecision
                ? e.decision?.decision === 'deny'
                  ? 'risk-critical'
                  : e.decision?.decision === 'ask'
                    ? 'risk-high'
                    : 'risk-none'
                : riskClass(e.risk)
              return (
                <button
                  key={e.id}
                  className={`event-row ${e.id === selectedId ? 'active' : ''} ${rowClass}`}
                  onClick={() => onSelect(e.id)}
                >
                  <span className={`type-chip type-${e.type}`}>{typeLabel(e.type)}</span>
                  <span className="event-main">
                    <strong>{eventName(e)}</strong>
                    <small>
                      {isDecision
                        ? `${decisionLabel(e.decision?.decision)} · ${e.decision?.reason || e.protocol}`
                        : `${e.protocol}${e.findings?.length ? ` · ${e.findings.length} 项命中` : ''}`}
                    </small>
                  </span>
                  <span className="event-session">{e.sessionId || '—'}</span>
                  <span className={`badge ${badgeClass}`}>
                    {isDecision ? decisionLabel(e.decision?.decision) : riskLabel(e.risk)}
                  </span>
                  <span className="event-time">{formatTime(e.timestamp)}</span>
                </button>
              )
            })}
          </div>
        )}
      </section>
    </div>
  )
}
