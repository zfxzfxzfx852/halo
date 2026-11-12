import { useEffect, useState } from 'react'
import type { HaloEvent } from '../types'

// 连接 WebSocket 实时事件流，断线自动重连。
export function useEventStream(onEvent: (e: HaloEvent) => void): boolean {
  const [connected, setConnected] = useState(false)

  useEffect(() => {
    let ws: WebSocket
    let retry: ReturnType<typeof setTimeout>
    let closed = false

    const connect = () => {
      const proto = location.protocol === 'https:' ? 'wss' : 'ws'
      ws = new WebSocket(`${proto}://${location.host}/api/v1/stream`)
      ws.onopen = () => setConnected(true)
      ws.onmessage = (ev) => {
        try {
          onEvent(JSON.parse(ev.data))
        } catch {
          // 忽略无法解析的消息
        }
      }
      ws.onerror = () => ws.close()
      ws.onclose = () => {
        setConnected(false)
        if (!closed) retry = setTimeout(connect, 1500)
      }
    }

    connect()
    return () => {
      closed = true
      clearTimeout(retry)
      ws?.close()
    }
  }, [onEvent])

  return connected
}
