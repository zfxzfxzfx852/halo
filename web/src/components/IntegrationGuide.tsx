import { useState } from 'react'
import type { RuntimeInfo } from '../types'
import { Icon } from './Icon'

interface Props {
  runtime?: RuntimeInfo
  onCopy: (text: string, label: string) => void
}

const TABS = [
  { key: 'http', label: 'HTTP 上游' },
  { key: 'mcp', label: 'MCP 桥接' },
  { key: 'client', label: '客户端配置' },
] as const

type TabKey = (typeof TABS)[number]['key']

const INSTALL = `curl -fsSL https://raw.githubusercontent.com/zfxzfxzfx852/halo/main/scripts/install.sh | sh`

const SNIPPETS: Record<TabKey, { title: string; note: string; code: string }[]> = {
  http: [
    {
      title: '1. 写一份 halo.yaml',
      note: '把上游指向你真正要调用的模型服务。',
      code: `host: 127.0.0.1
port: 8787
mode: observe      # 先观察，稳定后改 guard
upstream:
  baseUrl: https://api.openai.com
storage:
  retentionDays: 30`,
    },
    {
      title: '2. 启动 Halo',
      note: '控制台与代理同端口，默认 8787。',
      code: `halo --config halo.yaml
# 控制台: http://127.0.0.1:8787/console/`,
    },
    {
      title: '3. 让客户端指向 Halo',
      note: '只改 baseURL，业务代码零改动。',
      code: `export OPENAI_BASE_URL=http://127.0.0.1:8787/v1

curl http://127.0.0.1:8787/v1/chat/completions \\
  -H "Authorization: Bearer $OPENAI_API_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{"model":"gpt-4o","messages":[{"role":"user","content":"hi"}]}'`,
    },
  ],
  mcp: [
    {
      title: '包装任意 MCP Server',
      note: 'Halo 在 stdio 两端做双向桥接，记录每一次 tools/call。',
      code: `halo mcp --config halo.yaml -- \\
  npx -y @modelcontextprotocol/server-filesystem /tmp`,
    },
    {
      title: '写入 MCP 客户端',
      note: '把原来的启动命令前面加上 halo mcp 即可。',
      code: `{
  "mcpServers": {
    "filesystem": {
      "command": "halo",
      "args": ["mcp", "--", "npx", "-y",
        "@modelcontextprotocol/server-filesystem", "/tmp"]
    }
  }
}`,
    },
    {
      title: '验证事件',
      note: '调用一次工具，事件与审计记录会同时出现。',
      code: `curl -s http://127.0.0.1:8787/api/v1/events?type=tool_call
curl -s http://127.0.0.1:8787/api/v1/audit/verify`,
    },
  ],
  client: [
    {
      title: 'Claude Code',
      note: '在 MCP 配置里用 halo mcp 包住原 server。',
      code: `claude mcp add filesystem \\
  -- halo mcp -- npx -y @modelcontextprotocol/server-filesystem /tmp`,
    },
    {
      title: 'Cursor / Cline',
      note: '在 MCP 设置中新增 command 类型的 server。',
      code: `{
  "mcpServers": {
    "filesystem": {
      "command": "halo",
      "args": ["mcp", "--", "npx", "-y",
        "@modelcontextprotocol/server-filesystem", "/tmp"]
    }
  }
}`,
    },
    {
      title: '自检',
      note: '确认控制台与接口都活着。',
      code: `curl -s http://127.0.0.1:8787/healthz
curl -s http://127.0.0.1:8787/api/v1/runtime`,
    },
  ],
}

export function IntegrationGuide({ runtime, onCopy }: Props) {
  const [tab, setTab] = useState<TabKey>('http')
  const addr = runtime?.addr || '127.0.0.1:8787'

  return (
    <section className="guide-view">
      <article className="card glass guide-hero">
        <div>
          <span className="eyebrow">
            <Icon name="sparkles" size={16} />
            三分钟接入
          </span>
          <h2>不改代码，把 Halo 放到调用链中间</h2>
          <p>
            Halo 是透明代理：Agent 以为在调用模型，模型以为在回应 Agent，
            中间多了一层可观测、可管控、可审计的安全层。
          </p>
        </div>
        <div className="flow">
          <span className="flow-node">Agent / SDK</span>
          <i />
          <span className="flow-node halo">Halo · {addr}</span>
          <i />
          <span className="flow-node">Model / Tools</span>
        </div>
      </article>

      <article className="card glass">
        <header className="card-head">
          <h2>
            <Icon name="terminal" size={18} />
            安装
          </h2>
        </header>
        <button className="code-line" onClick={() => onCopy(INSTALL, '安装命令已复制')}>
          <code>{INSTALL}</code>
          <Icon name="copy" size={16} />
        </button>
        <p className="card-desc">
          也支持 Homebrew、Docker 与源码构建：<code>brew install halo</code> ·
          <code>docker compose up</code> · <code>make build</code>。
        </p>
      </article>

      <div className="guide-tabs">
        {TABS.map((t) => (
          <button
            key={t.key}
            className={`chip ${tab === t.key ? 'active' : ''}`}
            onClick={() => setTab(t.key)}
          >
            {t.label}
          </button>
        ))}
      </div>

      <div className="snippet-grid">
        {SNIPPETS[tab].map((s) => (
          <article className="card glass snippet-card" key={s.title}>
            <header>
              <h3>{s.title}</h3>
              <button className="icon-btn" onClick={() => onCopy(s.code, '命令已复制')}>
                <Icon name="copy" size={16} />
              </button>
            </header>
            <p>{s.note}</p>
            <pre className="code-block">{s.code}</pre>
          </article>
        ))}
      </div>

      <article className="card glass">
        <header className="card-head">
          <h2>
            <Icon name="layers" size={18} />
            它到底能做什么
          </h2>
        </header>
        <div className="capability-grid">
          <Cap icon="pulse" title="实时可视化" text="请求、响应、工具调用、决策在同一条时间线里，图谱随流量自动生长。" />
          <Cap icon="shield" title="风险检测" text="离线规则识别密钥、PII、提示注入与命令注入，风险自动升级。" />
          <Cap icon="seal" title="策略管控" text="按工具、方法、命令、参数做匹配，支持 allow / deny / ask。" />
          <Cap icon="check" title="人工审批" text="高风险调用挂起，允许或拒绝都写进审计链，超时默认拒绝。" />
          <Cap icon="graph" title="可验证审计" text="SHA-256 哈希链，重启续链，篡改可定位，JSONL 可导出复核。" />
          <Cap icon="settings" title="沙箱执行" text="进程级隔离、独立进程组、超时杀组，预留 Docker/gVisor 驱动。" />
        </div>
      </article>
    </section>
  )
}

function Cap({
  icon,
  title,
  text,
}: {
  icon: 'pulse' | 'shield' | 'seal' | 'check' | 'graph' | 'settings'
  title: string
  text: string
}) {
  return (
    <div className="cap">
      <span className="cap-icon">
        <Icon name={icon} size={19} />
      </span>
      <strong>{title}</strong>
      <p>{text}</p>
    </div>
  )
}
