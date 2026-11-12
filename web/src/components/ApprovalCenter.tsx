import type { ApprovalRequest } from '../types'
import { formatTime, riskClass } from '../lib/ui'

interface Props {
  approvals: ApprovalRequest[]
  onResolve: (id: string, decision: 'allow' | 'deny') => void
}

export function ApprovalCenter({ approvals, onResolve }: Props) {
  if (approvals.length === 0) {
    return <div className="empty">没有待审批的调用，所有请求都已处理</div>
  }
  return (
    <div className="approval-list">
      {approvals.map((a) => (
        <div key={a.ID} className="approval-card">
          <div className="approval-info">
            <strong>{a.Summary}</strong>
            <div className="approval-meta">
              <span className={`risk ${riskClass(a.Risk)}`}>
                {a.Risk || 'none'}
              </span>
              <span className="event-time">{formatTime(a.CreatedAt)}</span>
            </div>
          </div>
          <div className="approval-actions">
            <button
              className="ap-btn deny"
              onClick={() => onResolve(a.ID, 'deny')}
            >
              拒绝
            </button>
            <button
              className="ap-btn allow"
              onClick={() => onResolve(a.ID, 'allow')}
            >
              允许
            </button>
          </div>
        </div>
      ))}
    </div>
  )
}
