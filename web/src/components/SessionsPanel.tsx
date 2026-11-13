import { useMemo, useState } from 'react'
import type { SessionStat } from '../types'
import { Icon } from './Icon'
import { formatNumber, fromNow, riskClass, riskLabel } from '../lib/ui'

function sessionLabel(id: string): string {
  return id || '未命名会话'
}

interface Props {
  sessions: SessionStat[]
  onOpen: (sessionId: string) => void
}

export function SessionsPanel({ sessions, onOpen }: Props) {
  const [q, setQ] = useState('')

  const filtered = useMemo(() => {
    const key = q.trim().toLowerCase()
    if (!key) return sessions
    return sessions.filter(
      (s) =>
        sessionLabel(s.sessionId).toLowerCase().includes(key) ||
        (s.agent || '').toLowerCase().includes(key),
    )
  }, [sessions, q])

  return (
    <section className="sessions-view">
      <header className="section-head">
        <div>
          <h2>
            <Icon name="layers" size={18} />
            会话
          </h2>
          <p>按会话聚合风险与工具调用，点进去可直接筛选事件流。</p>
        </div>
        <label className="search-field small">
          <Icon name="search" size={16} />
          <input
            placeholder="搜索会话 / Agent"
            value={q}
            onChange={(e) => setQ(e.target.value)}
          />
        </label>
      </header>

      {filtered.length === 0 ? (
        <div className="card glass empty-state tall">
          <div className="empty-orb">
            <Icon name="layers" size={26} />
          </div>
          <h3>暂无会话</h3>
          <p>接入流量后，一次 Agent 会话会自动归并为一条记录。</p>
        </div>
      ) : (
        <div className="session-grid">
          {filtered.map((s) => (
            <button className="session-card glass" key={s.sessionId} onClick={() => onOpen(s.sessionId)}>
              <div className="session-top">
                <span className="avatar">{sessionLabel(s.sessionId).slice(0, 2).toUpperCase()}</span>
                <span className={`badge ${riskClass(s.maxRisk)}`}>{riskLabel(s.maxRisk)}</span>
              </div>
              <h3>{sessionLabel(s.sessionId)}</h3>
              <p>{s.agent || '未知 Agent'}</p>
              <div className="session-metrics">
                <div>
                  <strong>{formatNumber(s.events)}</strong>
                  <span>事件</span>
                </div>
                <div>
                  <strong>{formatNumber(s.toolCalls)}</strong>
                  <span>工具调用</span>
                </div>
                <div>
                  <strong>{fromNow(s.lastSeen)}</strong>
                  <span>最后活动</span>
                </div>
              </div>
            </button>
          ))}
        </div>
      )}
    </section>
  )
}
