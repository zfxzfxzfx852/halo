import { useMemo } from 'react'
import type {
  FindingStat,
  RiskLevel,
  HaloEvent,
  RuntimeInfo,
  SessionStat,
  Stats,
} from '../types'
import { Icon } from './Icon'
import { formatNumber, fromNow, riskClass, riskLabel } from '../lib/ui'

interface Props {
  stats?: Stats
  runtime?: RuntimeInfo
  events: HaloEvent[]
  sessions: SessionStat[]
  findings: FindingStat[]
  approvals: number
  connected: boolean
  onNavigate: (view: string) => void
  onCopy: (text: string, label: string) => void
  onSelectSession: (id: string) => void
}

const TYPE_ORDER = ['request', 'response', 'tool_call', 'tool_result', 'decision', 'alert']
const TYPE_NAMES: Record<string, string> = {
  request: '请求',
  response: '响应',
  tool_call: '工具调用',
  tool_result: '工具结果',
  decision: '决策',
  alert: '告警',
}

const RISK_ORDER = ['critical', 'high', 'medium', 'low', 'none']
const RISK_COLORS: Record<string, string> = {
  critical: '#b65c55',
  high: '#a97b3f',
  medium: '#a98a45',
  low: '#6a92bd',
  none: '#4f8a70',
}

export function Overview({
  stats,
  runtime,
  events,
  sessions,
  findings,
  approvals,
  connected,
  onNavigate,
  onCopy,
  onSelectSession,
}: Props) {
  const risk = useMemo(() => {
    const byRisk = stats?.ByRisk || {}
    const total = RISK_ORDER.reduce((s, k) => s + (byRisk[k] || 0), 0)
    let acc = 0
    const segments = RISK_ORDER.map((k) => {
      const v = byRisk[k] || 0
      const from = total ? (acc / total) * 360 : 0
      acc += v
      const to = total ? (acc / total) * 360 : 0
      return `${RISK_COLORS[k]} ${from}deg ${to}deg`
    }).join(', ')
    return {
      total,
      segments: total ? segments : 'rgba(128,128,128,0.15) 0deg 360deg',
      counts: byRisk,
    }
  }, [stats])

  const typeRows = useMemo(() => {
    const byType = stats?.ByType || {}
    const max = Math.max(1, ...TYPE_ORDER.map((t) => byType[t] || 0))
    return TYPE_ORDER.map((t) => ({ type: t, count: byType[t] || 0, pct: ((byType[t] || 0) / max) * 100 }))
  }, [stats])

  const highRisk = (risk.counts.high || 0) + (risk.counts.critical || 0)
  const toolCalls = stats?.ByType?.tool_call || 0
  const recent = events.slice(-6).reverse()
  const topFindings = findings.slice(0, 5)
  const recentSessions = sessions.slice(0, 4)

  const quickCmd = 'curl -fsSL https://raw.githubusercontent.com/zfxzfxzfx852/halo/main/scripts/install.sh | sh'
  const isEmpty = stats ? (stats.Total || 0) === 0 : false

  return (
    <div className="overview">
      <section className="hero-card glass">
        <div className="hero-copy">
          <span className="eyebrow">
            <span className={`live-dot ${connected ? 'on' : ''}`} />
            {connected ? '实时流已连接' : '正在重连实时流'}
          </span>
          <h1>
            看清每一次 <em>Agent</em> 调用
          </h1>
          <p>把 Agent 的每次模型与工具调用，变成可观测、可管控、可审计的本地数据。</p>
          <div className="hero-actions">
            <button className="btn primary" onClick={() => onNavigate('events')}>
              <Icon name="pulse" size={17} />
              进入实时事件流
            </button>
            <button className="btn ghost" onClick={() => onNavigate('guide')}>
              <Icon name="book" size={17} />
              3 分钟接入
            </button>
          </div>
        </div>
        <div className="hero-ring">
          <div className="risk-donut" style={{ background: `conic-gradient(${risk.segments})` }}>
            <div className="donut-core">
              <strong>{formatNumber(risk.total)}</strong>
              <span>已捕获事件</span>
            </div>
          </div>
          <div className="donut-legend">
            {RISK_ORDER.map((k) => (
              <span key={k}>
                <i style={{ background: RISK_COLORS[k] }} />
                {riskLabel(k as RiskLevel)} {risk.counts[k] || 0}
              </span>
            ))}
          </div>
        </div>
      </section>

      {isEmpty && (
        <section className="card glass onboarding">
          <span className="onboarding-orb">
            <Icon name="sparkles" size={24} />
          </span>
          <div className="onboarding-copy">
            <h2>还没有事件，先接入第一个 Agent</h2>
            <p>Halo 已在运行。把 Agent 的模型地址指向它，下一次模型或工具调用就会被捕获，风险、审批和审计自动出现。</p>
          </div>
          <div className="onboarding-actions">
            <button
              className="btn ghost"
              onClick={() => onCopy('http://127.0.0.1:8787/v1', '接入地址已复制')}
            >
              复制接入地址
            </button>
            <button className="btn primary" onClick={() => onNavigate('guide')}>
              查看接入指南
            </button>
          </div>
          <div className="onboarding-steps">
            <div className="onboarding-step">
              <b>01</b>
              <div>
                <strong>指向 Halo</strong>
                <code>http://127.0.0.1:8787/v1</code>
              </div>
            </div>
            <div className="onboarding-step">
              <b>02</b>
              <div>
                <strong>发起一次调用</strong>
                <p>对话、工具调用或 MCP 请求都可以。</p>
              </div>
            </div>
            <div className="onboarding-step">
              <b>03</b>
              <div>
                <strong>回到事件流</strong>
                <p>链路、风险、审批与审计会实时出现。</p>
              </div>
            </div>
          </div>
        </section>
      )}

      <section className="kpi-grid">
        <Kpi
          icon="layers"
          label="事件总量"
          value={formatNumber(stats?.Total)}
          hint={`${formatNumber(toolCalls)} 次工具调用`}
          accent="blue"
          onClick={() => onNavigate('events')}
        />
        <Kpi
          icon="shield"
          label="高风险事件"
          value={formatNumber(highRisk)}
          hint={highRisk ? '需要关注' : '当前干净'}
          accent={highRisk ? 'red' : 'green'}
          onClick={() => onNavigate('events')}
        />
        <Kpi
          icon="seal"
          label="待审批"
          value={formatNumber(approvals)}
          hint={approvals ? '等待人工放行' : '没有积压'}
          accent={approvals ? 'orange' : 'green'}
          onClick={() => onNavigate('approvals')}
        />
        <Kpi
          icon="check"
          label="运行模式"
          value={runtime?.mode === 'guard' ? 'Guard' : 'Observe'}
          hint={runtime ? `v${runtime.version} · ${runtime.addr}` : '读取中'}
          accent="purple"
          onClick={() => onNavigate('runtime')}
        />
      </section>

      <div className="two-col">
        <section className="card glass">
          <header className="card-head">
            <h2>
              <Icon name="graph" size={18} />
              调用构成
            </h2>
            <button className="link" onClick={() => onNavigate('events')}>
              查看全部 <Icon name="chevron" size={14} />
            </button>
          </header>
          <div className="type-bars">
            {typeRows.map((r) => (
              <div className="type-row" key={r.type}>
                <span className="type-name">{TYPE_NAMES[r.type]}</span>
                <div className="track">
                  <i className={`fill fill-${r.type}`} style={{ width: `${Math.max(2, r.pct)}%` }} />
                </div>
                <span className="type-count">{r.count}</span>
              </div>
            ))}
          </div>
        </section>

        <section className="card glass">
          <header className="card-head">
            <h2>
              <Icon name="risk" size={18} />
              风险洞察
            </h2>
            <span className="subtle">{findings.length} 类规则命中</span>
          </header>
          {topFindings.length === 0 ? (
            <div className="mini-empty">
              <Icon name="check" size={22} />
              <span>还没有命中检测规则，流量保持干净</span>
            </div>
          ) : (
            <ul className="finding-list">
              {topFindings.map((f) => (
                <li key={f.rule}>
                  <span className={`risk-dot ${riskClass(f.maxRisk)}`} />
                  <code>{f.rule}</code>
                  <b>{f.count}</b>
                </li>
              ))}
            </ul>
          )}
        </section>
      </div>

      <div className="two-col">
        <section className="card glass">
          <header className="card-head">
            <h2>
              <Icon name="clock" size={18} />
              最近活动
            </h2>
            <span className="subtle">自动刷新</span>
          </header>
          {recent.length === 0 ? (
            <div className="mini-empty">
              <Icon name="pulse" size={22} />
              <span>等待第一条事件</span>
            </div>
          ) : (
            <ul className="recent-list">
              {recent.map((e) => (
                <li key={e.id} onClick={() => onNavigate('events')}>
                  <span className={`type-chip type-${e.type}`}>{e.type.replace('_', ' ')}</span>
                  <div>
                    <strong>{e.toolCall?.name || e.method || e.type}</strong>
                    <small>{e.sessionId} · {e.protocol}</small>
                  </div>
                  <span className={`badge ${riskClass(e.risk)}`}>{riskLabel(e.risk)}</span>
                </li>
              ))}
            </ul>
          )}
        </section>

        <section className="card glass">
          <header className="card-head">
            <h2>
              <Icon name="layers" size={18} />
              活跃会话
            </h2>
            <button className="link" onClick={() => onNavigate('sessions')}>
              全部 <Icon name="chevron" size={14} />
            </button>
          </header>
          {recentSessions.length === 0 ? (
            <div className="mini-empty">
              <Icon name="clock" size={22} />
              <span>暂无会话</span>
            </div>
          ) : (
            <ul className="session-mini">
              {recentSessions.map((s) => (
                <li key={s.sessionId} onClick={() => onSelectSession(s.sessionId)}>
                  <span className="avatar">{(s.sessionId || '未命名').slice(0, 2).toUpperCase()}</span>
                  <div>
                    <strong>{s.sessionId || '未命名会话'}</strong>
                    <small>{s.agent || '未知 Agent'} · {s.events} 事件</small>
                  </div>
                  <span className={`badge ${riskClass(s.maxRisk)}`}>{riskLabel(s.maxRisk)}</span>
                  <small className="when">{fromNow(s.lastSeen)}</small>
                </li>
              ))}
            </ul>
          )}
        </section>
      </div>

      <section className="card glass install-card">
        <div>
          <h2>
            <Icon name="terminal" size={18} />
            一条命令接入
          </h2>
          <p>安装后设置上游地址，客户端无需改代码即可被 Halo 观测与保护。</p>
        </div>
        <button className="code-line" onClick={() => onCopy(quickCmd, '安装命令已复制')}>
          <code>{quickCmd}</code>
          <Icon name="copy" size={16} />
        </button>
        <div className="install-hints">
          <span>HTTP：<code>HALO_UPSTREAM=http://127.0.0.1:8787</code></span>
          <span>MCP：<code>halo mcp -- npx -y your-mcp-server</code></span>
        </div>
      </section>
    </div>
  )
}

function Kpi({
  icon,
  label,
  value,
  hint,
  accent,
  onClick,
}: {
  icon: 'layers' | 'shield' | 'seal' | 'check'
  label: string
  value: string
  hint: string
  accent: string
  onClick: () => void
}) {
  return (
    <button className={`kpi glass accent-${accent}`} onClick={onClick}>
      <span className="kpi-icon">
        <Icon name={icon} size={20} />
      </span>
      <span className="kpi-label">{label}</span>
      <strong>{value}</strong>
      <small>{hint}</small>
    </button>
  )
}
