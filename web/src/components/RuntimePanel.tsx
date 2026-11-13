import type { RuntimeInfo } from '../types'
import { Icon } from './Icon'

interface Props {
  runtime?: RuntimeInfo
  connected: boolean
}

export function RuntimePanel({ runtime, connected }: Props) {
  if (!runtime) {
    return (
      <div className="card glass empty-state tall">
        <div className="empty-orb">
          <Icon name="settings" size={26} />
        </div>
        <h3>正在读取运行信息</h3>
      </div>
    )
  }

  return (
    <section className="runtime-view">
      <article className="card glass runtime-hero">
        <div className="runtime-brand">
          <span className="brand-orb" />
          <div>
            <h2>Halo</h2>
            <p>本地优先的 AI Agent 可观测与安全控制平面</p>
          </div>
        </div>
        <div className="runtime-badges">
          <span className={`badge ${connected ? 'risk-none' : 'risk-high'}`}>
            {connected ? '实时连接正常' : '实时连接断开'}
          </span>
          <span className="badge risk-none">{runtime.mode}</span>
          <span className="badge risk-none">v{runtime.version}</span>
        </div>
      </article>

      <div className="runtime-grid">
        <Info label="监听地址" value={runtime.addr} mono />
        <Info label="数据目录" value={runtime.dataDir} mono />
        <Info label="保留天数" value={runtime.retentionDays > 0 ? `${runtime.retentionDays} 天` : '不自动清理'} />
        <Info label="上游服务" value={runtime.upstream || '未配置（纯观测模式）'} mono />
        <Info label="策略文件" value={runtime.policyFile || '未配置'} mono />
        <Info label="实时订阅" value={`${runtime.wsClients} 个 WebSocket 客户端`} />
        <Info label="构建时间" value={runtime.buildTime} />
        <Info label="提交" value={runtime.commit} mono />
      </div>

      <div className="two-col">
        <article className="card glass">
          <header className="card-head">
            <h2>
              <Icon name="sparkles" size={18} />
              为什么轻
            </h2>
          </header>
          <ul className="perf-list">
            <li><b>单二进制</b><span>前端嵌入 Go 二进制，一个文件即可运行</span></li>
            <li><b>零外部依赖</b><span>纯 Go SQLite，无需 Docker、数据库或中间件</span></li>
            <li><b>内存优先</b><span>事件总线在内存中分发，落库异步批处理</span></li>
            <li><b>增量渲染</b><span>前端只保留最近 500 条实时事件，避免无限增长</span></li>
          </ul>
        </article>

        <article className="card glass">
          <header className="card-head">
            <h2>
              <Icon name="terminal" size={18} />
            常用接口
            </h2>
          </header>
          <ul className="api-list">
            <li><code>GET /api/v1/events</code><span>事件查询</span></li>
            <li><code>GET /api/v1/stats</code><span>聚合统计</span></li>
            <li><code>GET /api/v1/sessions</code><span>会话画像</span></li>
            <li><code>GET /api/v1/findings</code><span>风险洞察</span></li>
            <li><code>GET /api/v1/audit/verify</code><span>审计校验</span></li>
            <li><code>GET /api/v1/stream</code><span>WebSocket 实时流</span></li>
          </ul>
        </article>
      </div>
    </section>
  )
}

function Info({ label, value, mono }: { label: string; value: string; mono?: boolean }) {
  return (
    <div className="info-card glass">
      <span>{label}</span>
      <strong className={mono ? 'mono' : ''}>{value}</strong>
    </div>
  )
}
