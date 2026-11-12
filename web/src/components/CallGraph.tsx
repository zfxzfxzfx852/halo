import { useMemo } from 'react'
import {
  Background,
  Controls,
  MarkerType,
  ReactFlow,
  type Edge,
  type Node,
} from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import type { HaloEvent } from '../types'

export function CallGraph({ events }: { events: HaloEvent[] }) {
  const { nodes, edges } = useMemo(() => buildGraph(events), [events])
  return (
    <ReactFlow
      nodes={nodes}
      edges={edges}
      fitView
      proOptions={{ hideAttribution: true }}
    >
      <Background />
      <Controls />
    </ReactFlow>
  )
}

function buildGraph(events: HaloEvent[]): { nodes: Node[]; edges: Edge[] } {
  const toolNames: string[] = []
  const risky = new Set<string>()

  events.forEach((e) => {
    const name = e.toolCall?.name
    if (name && !toolNames.includes(name)) {
      toolNames.push(name)
      if (e.risk && e.risk !== 'none') risky.add(name)
    }
  })

  const nodes: Node[] = [
    {
      id: 'agent',
      type: 'input',
      position: { x: 0, y: 160 },
      data: { label: 'Agent' },
      className: 'node-agent',
    },
  ]
  const edges: Edge[] = []

  toolNames.forEach((name, i) => {
    const id = `tool-${name}`
    nodes.push({
      id,
      position: { x: 340, y: i * 100 },
      data: { label: name },
      className: risky.has(name) ? 'node-risk' : 'node-tool',
    })
    edges.push({
      id: `edge-${name}`,
      source: 'agent',
      target: id,
      animated: risky.has(name),
      markerEnd: { type: MarkerType.ArrowClosed },
    })
  })

  return { nodes, edges }
}
