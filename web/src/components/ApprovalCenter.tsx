import { useState } from 'react'
import type { ApprovalRequest } from '../types'
import { Icon } from './Icon'
import { formatDateTime, fromNow, riskClass, riskLabel } from '../lib/ui'

interface Props {
  approvals: ApprovalRequest[]
  onResolve: (id: string, decision: 'allow' | 'deny') => Promise<void>
}

export function ApprovalCenter({ approvals, onResolve }: Props) {
  const [busy, setBusy] = useState<string>()

  const resolve = async (id: string, decision: 'allow' | 'deny') => {
    setBusy(id)
    try {
      await onResolve(id, decision)
    } finally {
      setBusy(undefined)
    }
  }

  if (approvals.length === 0) {
    return (
      <div className="card glass empty-state tall">
        <div className="empty-orb">
          <Icon name="check" size={26} />
        </div>
        <h3>没有待审批的调用</h3>
        <p>策略命中 ask 的请求会挂起在这里，等待你允许或拒绝。超时默认拒绝。</p>
      </div>
    )
  }

  return (
    <section className="approvals-view">
      <header className="section-head">
        <div>
          <h2>
            <Icon name="seal" size={18} />
            人工审批
          </h2>
          <p>命中审批规则的调用暂停在 Halo，你的决定会写入审计链。</p>
        </div>
        <span className="count-pill">{approvals.length} 待处理</span>
      </header>

      <div className="approval-grid">
        {approvals.map((a) => (
          <article className="approval-card glass" key={a.ID}>
            <div className="approval-top">
              <div className="approval-badges">
                <span className="badge risk-ask">待审批</span>
                {a.Risk !== 'none' && (
                  <span className={`badge ${riskClass(a.Risk)}`}>{riskLabel(a.Risk)}</span>
                )}
              </div>
              <span className="when">{fromNow(a.CreatedAt)}</span>
            </div>
            <h3>{a.Summary}</h3>
            <div className="approval-meta">
              <code>{a.EventID}</code>
              <span>{formatDateTime(a.CreatedAt)}</span>
            </div>
            <p className="approval-hint">
              允许后请求继续转发；拒绝将返回 403 并阻断本次调用。
            </p>
            <div className="approval-actions">
              <button
                className="btn danger"
                disabled={busy === a.ID}
                onClick={() => resolve(a.ID, 'deny')}
              >
                {busy === a.ID ? '处理中…' : '拒绝'}
              </button>
              <button
                className="btn primary"
                disabled={busy === a.ID}
                onClick={() => resolve(a.ID, 'allow')}
              >
                {busy === a.ID ? '处理中…' : '允许'}
              </button>
            </div>
          </article>
        ))}
      </div>
    </section>
  )
}
