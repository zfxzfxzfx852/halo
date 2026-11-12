import type { HaloEvent } from '../types'
import { formatTime, riskClass } from '../lib/ui'

interface Props {
  events: HaloEvent[]
  selectedId?: string
  onSelect: (id: string) => void
}

export function EventList({ events, selectedId, onSelect }: Props) {
  if (events.length === 0) {
    return (
      <div className="empty">
        等待事件…将客户端指向 Halo 并发起请求，即可看到实时数据
      </div>
    )
  }
  return (
    <div className="event-list">
      {events.map((e) => (
        <button
          key={e.id}
          className={`event-row ${e.id === selectedId ? 'active' : ''}`}
          onClick={() => onSelect(e.id)}
        >
          <span className={`badge type-${e.type}`}>{e.type}</span>
          <span className="event-name">
            {e.toolCall?.name || e.method || e.protocol}
          </span>
          <span className={`risk ${riskClass(e.risk)}`}>
            {e.risk || 'none'}
          </span>
          <span className="event-time">{formatTime(e.timestamp)}</span>
        </button>
      ))}
    </div>
  )
}
