import type { HaloEvent } from '../types'

interface Props {
  event: HaloEvent
  onClose: () => void
}

export function EventDetail({ event, onClose }: Props) {
  return (
    <div className="drawer">
      <div className="drawer-head">
        <strong>{event.id}</strong>
        <button onClick={onClose} aria-label="close">
          ×
        </button>
      </div>
      <pre className="json">{JSON.stringify(event, null, 2)}</pre>
    </div>
  )
}
