import { useState } from 'react'
import type { HaloEvent } from '../types'
import { Icon } from './Icon'
import {
  decisionLabel,
  formatDateTime,
  riskClass,
  riskLabel,
  typeLabel,
} from '../lib/ui'

interface Props {
  event: HaloEvent
  onClose: () => void
  onCopy: (text: string, label: string) => void
}

export function EventDetail({ event, onClose, onCopy }: Props) {
  const [tab, setTab] = useState<'overview' | 'raw'>('overview')
  const json = JSON.stringify(event, null, 2)

  return (
    <div className="drawer-wrap" onClick={onClose}>
      <aside className="drawer glass" onClick={(e) => e.stopPropagation()}>
        <header className="drawer-head">
          <div>
            <span className={`type-chip type-${event.type}`}>{typeLabel(event.type)}</span>
            <h3>{event.toolCall?.name || event.method || event.id}</h3>
            <small>{event.id}</small>
          </div>
          <div className="drawer-head-actions">
            <button className="icon-btn" onClick={() => onCopy(json, '事件 JSON 已复制')} title="复制 JSON">
              <Icon name="copy" size={17} />
            </button>
            <button className="icon-btn" onClick={onClose} title="关闭">
              <Icon name="close" size={17} />
            </button>
          </div>
        </header>

        <div className="drawer-tabs">
          <button className={tab === 'overview' ? 'active' : ''} onClick={() => setTab('overview')}>
            概览
          </button>
          <button className={tab === 'raw' ? 'active' : ''} onClick={() => setTab('raw')}>
            原始 JSON
          </button>
        </div>

        <div className="drawer-body">
          {tab === 'overview' ? (
            <>
              <div className="detail-grid">
                <Detail label="风险等级">
                  <span className={`badge ${riskClass(event.risk)}`}>{riskLabel(event.risk)}</span>
                </Detail>
                <Detail label="协议">{event.protocol}</Detail>
                <Detail label="会话">{event.sessionId}</Detail>
                <Detail label="Agent">{event.agent || '—'}</Detail>
                <Detail label="时间">{formatDateTime(event.timestamp)}</Detail>
                <Detail label="父事件">{event.parentId || '—'}</Detail>
              </div>

              {event.decision && (
                <section className="detail-block">
                  <h4>策略决策</h4>
                  <div className={`decision-box decision-${event.decision.decision}`}>
                    <strong>{decisionLabel(event.decision.decision)}</strong>
                    <span>{event.decision.reason}</span>
                    {event.decision.policy && <small>规则：{event.decision.policy}</small>}
                    {event.decision.reviewedBy && (
                      <small>审批人：{event.decision.reviewedBy}</small>
                    )}
                  </div>
                </section>
              )}

              {!!event.findings?.length && (
                <section className="detail-block">
                  <h4>检测命中</h4>
                  <ul className="detail-findings">
                    {event.findings.map((f, i) => (
                      <li key={`${f.rule}-${i}`}>
                        <span className={`badge ${riskClass(f.level)}`}>{riskLabel(f.level)}</span>
                        <div>
                          <code>{f.rule}</code>
                          <p>{f.message}</p>
                        </div>
                      </li>
                    ))}
                  </ul>
                </section>
              )}

              {event.toolCall && (
                <section className="detail-block">
                  <h4>工具调用</h4>
                  <p className="detail-kv">
                    <span>名称</span>
                    <code>{event.toolCall.name}</code>
                  </p>
                  {event.toolCall.server && (
                    <p className="detail-kv">
                      <span>服务</span>
                      <code>{event.toolCall.server}</code>
                    </p>
                  )}
                  <pre className="code-block">
                    {JSON.stringify(event.toolCall.arguments ?? {}, null, 2)}
                  </pre>
                </section>
              )}

              {event.request && (
                <section className="detail-block">
                  <h4>请求</h4>
                  <p className="detail-kv">
                    <span>方法</span>
                    <code>{event.method || '—'}</code>
                  </p>
                  <Headers headers={event.request.headers} />
                  <pre className="code-block">{pretty(event.request.body)}</pre>
                </section>
              )}

              {event.response && (
                <section className="detail-block">
                  <h4>响应</h4>
                  <p className="detail-kv">
                    <span>状态</span>
                    <code>{event.response.status ?? '—'}</code>
                  </p>
                  <Headers headers={event.response.headers} />
                  <pre className="code-block">{pretty(event.response.body)}</pre>
                </section>
              )}

              {event.toolResult && (
                <section className="detail-block">
                  <h4>工具结果</h4>
                  <pre className="code-block">{pretty(event.toolResult.content)}</pre>
                </section>
              )}
            </>
          ) : (
            <pre className="code-block raw-json">{json}</pre>
          )}
        </div>
      </aside>
    </div>
  )
}

function Detail({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="detail-item">
      <span>{label}</span>
      <strong>{children}</strong>
    </div>
  )
}

function Headers({ headers }: { headers?: Record<string, string> }) {
  if (!headers || Object.keys(headers).length === 0) return null
  return (
    <p className="detail-kv">
      <span>Headers</span>
      <code>
        {Object.entries(headers)
          .slice(0, 6)
          .map(([k, v]) => `${k}: ${v}`)
          .join('\n')}
      </code>
    </p>
  )
}

function pretty(v: unknown): string {
  if (v === undefined || v === null) return '（空）'
  if (typeof v === 'string') {
    try {
      return JSON.stringify(JSON.parse(v), null, 2)
    } catch {
      return v
    }
  }
  return JSON.stringify(v, null, 2)
}
