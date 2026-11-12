import { useCallback, useEffect, useMemo, useState } from 'react'
import { ApprovalCenter } from './components/ApprovalCenter'
import { CallGraph } from './components/CallGraph'
import { EventDetail } from './components/EventDetail'
import { EventList } from './components/EventList'
import { useEventStream } from './hooks/useEventStream'
import { fetchApprovals, fetchEvents, resolveApproval } from './lib/api'
import type { ApprovalRequest, HaloEvent } from './types'

const MAX_EVENTS = 500

export default function App() {
  const [events, setEvents] = useState<HaloEvent[]>([])
  const [selectedId, setSelectedId] = useState<string>()
  const [filter, setFilter] = useState('')
  const [view, setView] = useState<'list' | 'graph' | 'approvals'>('list')
  const [approvals, setApprovals] = useState<ApprovalRequest[]>([])

  const refreshApprovals = useCallback(() => {
    fetchApprovals()
      .then((d) => setApprovals(d.items))
      .catch(() => undefined)
  }, [])

  useEffect(() => {
    fetchEvents({ limit: 200 })
      .then((d) => setEvents(d.items.slice().reverse()))
      .catch(() => undefined)
    refreshApprovals()
  }, [refreshApprovals])

  // 收到 decision 事件即刷新；另用轮询兜底。
  const onEvent = useCallback(
    (e: HaloEvent) => {
      setEvents((prev) => [...prev, e].slice(-MAX_EVENTS))
      if (e.type === 'decision') {
        refreshApprovals()
      }
    },
    [refreshApprovals],
  )
  const connected = useEventStream(onEvent)

  useEffect(() => {
    const t = setInterval(refreshApprovals, 5000)
    return () => clearInterval(t)
  }, [refreshApprovals])

  const onResolve = useCallback(
    (id: string, decision: 'allow' | 'deny') => {
      resolveApproval(id, decision)
        .then(refreshApprovals)
        .catch(() => undefined)
    },
    [refreshApprovals],
  )

  const filtered = useMemo(() => {
    if (!filter) return events
    const f = filter.toLowerCase()
    return events.filter((e) =>
      (e.toolCall?.name || e.method || e.type || '')
        .toLowerCase()
        .includes(f),
    )
  }, [events])

  const selected = events.find((e) => e.id === selectedId)

  return (
    <div className="app">
      <header>
        <h1>Halo</h1>
        <span className={`conn ${connected ? 'on' : 'off'}`}>
          {connected ? '实时连接' : '未连接'}
        </span>
        <input
          placeholder="搜索方法 / 工具…"
          value={filter}
          onChange={(e) => setFilter(e.target.value)}
        />
        <div className="tabs">
          <button
            className={view === 'list' ? 'active' : ''}
            onClick={() => setView('list')}
          >
            事件流
          </button>
          <button
            className={view === 'graph' ? 'active' : ''}
            onClick={() => setView('graph')}
          >
            调用图谱
          </button>
          <button
            className={view === 'approvals' ? 'active' : ''}
            onClick={() => setView('approvals')}
          >
            审批{approvals.length > 0 ? ` (${approvals.length})` : ''}
          </button>
        </div>
      </header>
      <main>
        {view === 'list' && (
          <EventList
            events={filtered}
            selectedId={selectedId}
            onSelect={setSelectedId}
          />
        )}
        {view === 'graph' && <CallGraph events={events} />}
        {view === 'approvals' && (
          <ApprovalCenter approvals={approvals} onResolve={onResolve} />
        )}
        {selected && (
          <EventDetail
            event={selected}
            onClose={() => setSelectedId(undefined)}
          />
        )}
      </main>
    </div>
  )
}
