import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { ApprovalCenter } from './components/ApprovalCenter'
import { AuditPanel } from './components/AuditPanel'
import { CallGraph } from './components/CallGraph'
import { EventDetail } from './components/EventDetail'
import { EventList } from './components/EventList'
import { Icon } from './components/Icon'
import { IntegrationGuide } from './components/IntegrationGuide'
import { Overview } from './components/Overview'
import { ParticleField } from './components/ParticleField'
import { PoliciesPanel } from './components/PoliciesPanel'
import { RuntimePanel } from './components/RuntimePanel'
import { SessionsPanel } from './components/SessionsPanel'
import { useEventStream } from './hooks/useEventStream'
import {
  fetchApprovals,
  fetchEvent,
  fetchEvents,
  fetchFindings,
  fetchPolicies,
  fetchRuntime,
  fetchSessions,
  fetchStats,
  resolveApproval,
} from './lib/api'
import type {
  ApprovalRequest,
  FindingStat,
  HaloEvent,
  PolicyRule,
  RuntimeInfo,
  SessionStat,
  Stats,
} from './types'

const MAX_EVENTS = 500
const MAIN_TABS = 4

const TABS = [
  { key: 'overview', label: '总览', icon: 'sparkles' },
  { key: 'events', label: '事件流', icon: 'pulse' },
  { key: 'graph', label: '图谱', icon: 'graph' },
  { key: 'approvals', label: '审批', icon: 'seal' },
  { key: 'audit', label: '审计', icon: 'shield' },
  { key: 'sessions', label: '会话', icon: 'layers' },
  { key: 'policies', label: '策略', icon: 'check' },
  { key: 'guide', label: '接入', icon: 'book' },
  { key: 'runtime', label: '运行', icon: 'settings' },
] as const

type ViewKey = (typeof TABS)[number]['key']

export default function App() {
  const [theme, setTheme] = useState<'light' | 'dark'>(() => {
    const saved = localStorage.getItem('halo-theme')
    if (saved === 'light' || saved === 'dark') return saved
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
  })
  const [view, setView] = useState<ViewKey>('overview')
  const [events, setEvents] = useState<HaloEvent[]>([])
  const [stats, setStats] = useState<Stats>()
  const [approvals, setApprovals] = useState<ApprovalRequest[]>([])
  const [sessions, setSessions] = useState<SessionStat[]>([])
  const [findings, setFindings] = useState<FindingStat[]>([])
  const [runtime, setRuntime] = useState<RuntimeInfo>()
  const [policy, setPolicy] = useState<{ default: string; rules: PolicyRule[] }>({
    default: '',
    rules: [],
  })
  const [selected, setSelected] = useState<HaloEvent>()
  const [selectedId, setSelectedId] = useState<string>()
  const [sessionFilter, setSessionFilter] = useState<string>()
  const [query, setQuery] = useState('')
  const [live, setLive] = useState(true)
  const [moreOpen, setMoreOpen] = useState(false)
  const [toast, setToast] = useState<string>()
  const liveRef = useRef(live)
  const searchRef = useRef<HTMLInputElement | null>(null)

  liveRef.current = live

  useEffect(() => {
    document.documentElement.dataset.theme = theme
    localStorage.setItem('halo-theme', theme)
  }, [theme])

  const showToast = useCallback((msg: string) => {
    setToast(msg)
  }, [])

  useEffect(() => {
    if (!toast) return
    const t = setTimeout(() => setToast(undefined), 2600)
    return () => clearTimeout(t)
  }, [toast])

  const refreshMeta = useCallback(() => {
    fetchStats().then(setStats).catch(() => undefined)
    fetchSessions().then((d) => setSessions(d.items || [])).catch(() => undefined)
    fetchFindings().then((d) => setFindings(d.items || [])).catch(() => undefined)
    fetchRuntime().then(setRuntime).catch(() => undefined)
    fetchPolicies().then((d) => setPolicy({ default: d.default || '', rules: d.rules || [] })).catch(() => undefined)
  }, [])

  const refreshApprovals = useCallback(() => {
    fetchApprovals()
      .then((d) => setApprovals(d.items || []))
      .catch(() => undefined)
  }, [])

  useEffect(() => {
    fetchEvents({ limit: 300 })
      .then((d) => setEvents(d.items || []))
      .catch(() => showToast('初始事件加载失败'))
    refreshMeta()
    refreshApprovals()
  }, [refreshMeta, refreshApprovals, showToast])

  useEffect(() => {
    const t = setInterval(refreshApprovals, 5000)
    const m = setInterval(refreshMeta, 12000)
    return () => {
      clearInterval(t)
      clearInterval(m)
    }
  }, [refreshApprovals, refreshMeta])

  const onEvent = useCallback(
    (e: HaloEvent) => {
      if (!liveRef.current) return
      setEvents((prev) => [e, ...prev].slice(0, MAX_EVENTS))
      if (e.type === 'decision') refreshApprovals()
    },
    [refreshApprovals],
  )
  const connected = useEventStream(onEvent)

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        setSelected(undefined)
        setSelectedId(undefined)
      }
      if (e.key === '/' && document.activeElement?.tagName !== 'INPUT') {
        e.preventDefault()
        searchRef.current?.focus()
      }
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [])

  const openEvent = useCallback(
    (id: string) => {
      setSelectedId(id)
      const hit = events.find((e) => e.id === id)
      if (hit) {
        setSelected(hit)
        return
      }
      fetchEvent(id)
        .then(setSelected)
        .catch(() => showToast('事件详情加载失败'))
    },
    [events, showToast],
  )

  const onResolve = useCallback(
    async (id: string, decision: 'allow' | 'deny') => {
      try {
        await resolveApproval(id, decision)
        showToast(decision === 'allow' ? '已允许该调用' : '已拒绝该调用')
        refreshApprovals()
      } catch {
        showToast('审批操作失败，可能已超时')
      }
    },
    [refreshApprovals, showToast],
  )

  const onCopy = useCallback(
    (text: string, label: string) => {
      navigator.clipboard
        ?.writeText(text)
        .then(() => showToast(label))
        .catch(() => showToast('复制失败，请手动选择'))
    },
    [showToast],
  )

  const filteredEvents = useMemo(() => {
    const key = query.trim().toLowerCase()
    if (!key) return events
    return events.filter((e) =>
      [e.id, e.sessionId, e.agent, e.method, e.toolCall?.name, e.type]
        .filter(Boolean)
        .some((v) => String(v).toLowerCase().includes(key)),
    )
  }, [events, query])

  const openSession = (sessionId: string) => {
    setSessionFilter(sessionId)
    setView('events')
  }

  return (
    <div className="app">
      <ParticleField theme={theme} />

      <header className="topbar glass">
        <div className="brand">
          <span className="brand-orb" />
          <div className="brand-text">
            <h1>Halo</h1>
            <span>Agent 可观测与安全控制平面</span>
          </div>
        </div>

        <label className="search-field">
          <Icon name="search" size={17} />
          <input
            ref={searchRef}
            placeholder="搜索事件 / 工具 / 会话…"
            value={query}
            onChange={(e) => {
              setQuery(e.target.value)
              if (e.target.value.trim()) setView('events')
            }}
          />
          <kbd>/</kbd>
        </label>

        <div className="top-actions">
          <span className={`conn-pill ${connected ? 'on' : 'off'}`}>
            <i />
            {connected ? '实时' : '重连中'}
          </span>
          <span className="mode-pill">{runtime?.mode || 'observe'}</span>
          <button
            className="icon-btn"
            title="切换主题"
            onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
          >
            <Icon name={theme === 'dark' ? 'sun' : 'moon'} size={18} />
          </button>
          <a
            className="icon-btn"
            href="https://github.com/zfxzfxzfx852/halo"
            target="_blank"
            rel="noreferrer"
            title="GitHub"
          >
            <Icon name="arrow" size={18} />
          </a>
        </div>
      </header>

      <nav className="tabbar glass" aria-label="主导航">
        {TABS.map((t) => (
          <button
            key={t.key}
            className={`tab ${view === t.key ? 'active' : ''}`}
            onClick={() => setView(t.key)}
          >
            <span className="tab-icon">
              <Icon name={t.icon} size={19} />
              {t.key === 'approvals' && approvals.length > 0 && (
                <i className="tab-badge">{approvals.length}</i>
              )}
            </span>
            <span>{t.label}</span>
          </button>
        ))}
      </nav>

      <nav className="mobile-tabbar glass" aria-label="移动导航">
        {TABS.slice(0, MAIN_TABS).map((t) => (
          <button
            key={t.key}
            className={`tab ${view === t.key ? 'active' : ''}`}
            onClick={() => {
              setView(t.key)
              setMoreOpen(false)
            }}
          >
            <span className="tab-icon">
              <Icon name={t.icon} size={19} />
              {t.key === 'approvals' && approvals.length > 0 && (
                <i className="tab-badge">{approvals.length}</i>
              )}
            </span>
            <span>{t.label}</span>
          </button>
        ))}
        <button
          className={`tab tab-more ${TABS.slice(MAIN_TABS).some((t) => t.key === view) ? 'active' : ''}`}
          onClick={() => setMoreOpen((v) => !v)}
        >
          <span className="tab-icon">
            <Icon name="more" size={19} />
          </span>
          <span>更多</span>
        </button>
      </nav>

      {moreOpen && (
        <>
          <div className="more-backdrop" onClick={() => setMoreOpen(false)} />
          <nav className="more-sheet glass" aria-label="更多功能">
            {TABS.slice(MAIN_TABS).map((t) => (
              <button
                key={t.key}
                className={`more-item ${view === t.key ? 'active' : ''}`}
                onClick={() => {
                  setView(t.key)
                  setMoreOpen(false)
                }}
              >
                <span className="more-icon">
                  <Icon name={t.icon} size={22} />
                </span>
                <span>{t.label}</span>
              </button>
            ))}
          </nav>
        </>
      )}

      <main className="main">
        {view === 'overview' && (
          <Overview
            stats={stats}
            runtime={runtime}
            events={events}
            sessions={sessions}
            findings={findings}
            approvals={approvals.length}
            connected={connected}
            onNavigate={(v) => setView(v as ViewKey)}
            onCopy={onCopy}
            onSelectSession={openSession}
          />
        )}
        {view === 'events' && (
          <EventList
            events={filteredEvents}
            selectedId={selectedId}
            onSelect={openEvent}
            live={live}
            onToggleLive={() => setLive((v) => !v)}
            sessionFilter={sessionFilter}
            onClearSession={() => setSessionFilter(undefined)}
            total={stats?.Total || events.length}
          />
        )}
        {view === 'graph' && <CallGraph events={events} onSelectEvent={openEvent} />}
        {view === 'approvals' && <ApprovalCenter approvals={approvals} onResolve={onResolve} />}
        {view === 'audit' && <AuditPanel onToast={showToast} />}
        {view === 'sessions' && <SessionsPanel sessions={sessions} onOpen={openSession} />}
        {view === 'policies' && (
          <PoliciesPanel
            defaultDecision={policy.default}
            rules={policy.rules}
            mode={runtime?.mode || 'observe'}
            onCopy={onCopy}
          />
        )}
        {view === 'guide' && <IntegrationGuide runtime={runtime} onCopy={onCopy} />}
        {view === 'runtime' && <RuntimePanel runtime={runtime} connected={connected} />}
      </main>

      <footer className="foot">
        <span>Halo · 本地优先 · 数据不出机器</span>
        <span>{stats?.Total || 0} 事件已归档</span>
      </footer>

      {selected && (
        <EventDetail
          event={selected}
          onClose={() => {
            setSelected(undefined)
            setSelectedId(undefined)
          }}
          onCopy={onCopy}
        />
      )}

      {toast && <div className="toast glass">{toast}</div>}
    </div>
  )
}
