import { useMemo } from 'react'
import {
  Background,
  Controls,
  Handle,
  MarkerType,
  Position,
  ReactFlow,
  type Edge,
  type Node,
  type NodeProps,
} from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import type { HaloEvent, RiskLevel } from '../types'
import { Icon } from './Icon'
import { riskLabel, severity } from '../lib/ui'

interface Props {
  events: HaloEvent[]
  onSelectEvent: (id: string) => void
}

interface ToolAgg {
  name: string
  count: number
  maxRisk: RiskLevel
  agents: Set<string>
}

const nodeTypes = {
  tool: ToolNode,
}

function ToolNode({ data }: NodeProps) {
  const d = data as unknown as { name: string; count: number; risk: RiskLevel }
  return (
    <div className="graph-tool-node">
      <Handle type="target" position={Position.Left} />
      <strong>{d.name}</strong>
      <span>
        {d.count} 次 · {riskLabel(d.risk)}
      </span>
    </div>
  )
}

function buildGraph(events: HaloEvent[]): { nodes: Node[]; edges: Edge[]; tools: number } {
  const tools = new Map<string, ToolAgg>()
  const agents = new Set<string>()
  const hasTools = events.some((e) => e.toolCall?.name)
  const source = hasTools ? events.filter((e) => e.toolCall?.name) : events

  for (const e of source) {
    const name = e.toolCall?.name || e.method
    if (!name) continue
    const agent = e.agent || 'local-agent'
    agents.add(agent)
    const cur = tools.get(name)
    if (!cur) {
      tools.set(name, {
        name,
        count: 1,
        maxRisk: e.risk,
        agents: new Set([agent]),
      })
    } else {
      cur.count += 1
      cur.agents.add(agent)
      if (severity(e.risk) > severity(cur.maxRisk)) cur.maxRisk = e.risk
    }
  }

  const nodes: Node[] = []
  const edges: Edge[] = []
  const agentList = Array.from(agents)
  const agentHeight = 96
  const agentOffset = -((agentList.length - 1) * agentHeight) / 2

  agentList.forEach((a, i) => {
    nodes.push({
      id: `agent-${a}`,
      type: 'input',
      position: { x: 0, y: 240 + agentOffset + i * agentHeight },
      data: { label: `Agent · ${a}` },
      className: 'node-agent',
    })
  })

  const list = Array.from(tools.values()).sort((a, b) => b.count - a.count)
  const height = 104
  const offset = -((list.length - 1) * height) / 2

  list.forEach((t, i) => {
    const risky = severity(t.maxRisk) >= 2
    nodes.push({
      id: `tool-${t.name}`,
      type: 'tool',
      position: { x: 520, y: 240 + offset + i * height },
      data: { name: t.name, count: t.count, risk: t.maxRisk },
      className: risky ? 'node-risk' : 'node-tool',
    })
    t.agents.forEach((a) => {
      edges.push({
        id: `edge-${a}-${t.name}`,
        source: `agent-${a}`,
        target: `tool-${t.name}`,
        animated: risky,
        style: risky ? { stroke: '#b65c55', strokeWidth: 1.6 } : undefined,
        markerEnd: { type: MarkerType.ArrowClosed },
      })
    })
  })

  return { nodes, edges, tools: list.length }
}

export function CallGraph({ events, onSelectEvent }: Props) {
  const { nodes, edges, tools } = useMemo(() => buildGraph(events), [events])

  if (tools === 0) {
    return (
      <div className="card glass empty-state tall">
        <div className="empty-orb">
          <Icon name="graph" size={26} />
        </div>
        <h3>调用图谱等待数据</h3>
        <p>接入 MCP 或 HTTP 上游后，Agent 与工具的调用关系会自动绘制在这里。</p>
      </div>
    )
  }

  return (
    <section className="card glass graph-card">
      <header className="card-head graph-head">
        <div>
          <h2>
            <Icon name="graph" size={18} />
            实时调用图谱
          </h2>
          <span className="subtle">点击工具节点查看最近一次调用详情</span>
        </div>
        <div className="graph-legend">
          <span><i className="dot safe" />安全</span>
          <span><i className="dot risk" />风险</span>
          <span className="subtle">{tools} 个工具</span>
        </div>
      </header>
      <div className="graph-canvas">
        <ReactFlow
          nodes={nodes}
          edges={edges}
          nodeTypes={nodeTypes}
          fitView
          fitViewOptions={{ padding: 0.22, maxZoom: 1.15 }}
          minZoom={0.4}
          maxZoom={2}
          proOptions={{ hideAttribution: true }}
          onNodeClick={(_, node) => {
            if (node.id.startsWith('tool-')) {
              const latest = events.find((e) => (e.toolCall?.name || e.method) === node.id.slice(5))
              if (latest) onSelectEvent(latest.id)
            }
          }}
        >
          <Background gap={22} size={1.2} color="rgba(120,140,170,0.18)" />
          <Controls showInteractive={false} />
        </ReactFlow>
      </div>
    </section>
  )
}
