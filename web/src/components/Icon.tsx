type IconName =
  | 'pulse'
  | 'list'
  | 'graph'
  | 'shield'
  | 'check'
  | 'seal'
  | 'layers'
  | 'book'
  | 'copy'
  | 'close'
  | 'search'
  | 'pause'
  | 'play'
  | 'chevron'
  | 'sparkles'
  | 'alert'
  | 'clock'
  | 'terminal'
  | 'download'
  | 'refresh'
  | 'risk'
  | 'settings'
  | 'arrow'
  | 'sun'
  | 'moon'
  | 'more'

const paths: Record<IconName, string> = {
  pulse: 'M3 12h4l3-7 4 14 3-7h4',
  list: 'M4 7h16M4 12h16M4 17h10',
  graph: 'M5 6h4v4H5zM15 14h4v4h-4zM9 8l6 8M15 8h4',
  shield: 'M12 3l8 3v6c0 5-3.5 8-8 9-4.5-1-8-4-8-9V6z',
  check: 'M5 13l4 4 10-11',
  seal: 'M12 3l2.4 2.1 3.2-.3.4 3.2 2.8 1.6-1.4 2.9 1.4 2.9-2.8 1.6-.4 3.2-3.2-.3L12 21l-2.4-2.1-3.2.3-.4-3.2-2.8-1.6 1.4-2.9-1.4-2.9 2.8-1.6.4-3.2 3.2.3z',
  layers: 'M12 3l9 5-9 5-9-5 9-5zM3 13l9 5 9-5M3 17l9 5 9-5',
  book: 'M6 4h11a2 2 0 012 2v14H7a2 2 0 01-2-2V6a2 2 0 012-2zM6 18h13',
  copy: 'M9 8h9a1 1 0 011 1v10a1 1 0 01-1 1H9a1 1 0 01-1-1V9a1 1 0 011-1zM6 4h9M5 6v10',
  close: 'M6 6l12 12M18 6L6 18',
  search: 'M11 4a7 7 0 100 14 7 7 0 000-14zM16.5 16.5L20 20',
  pause: 'M9 5v14M15 5v14',
  play: 'M8 5l11 7-11 7z',
  chevron: 'M9 6l6 6-6 6',
  sparkles: 'M12 3l1.6 4.4 4.4 1.6-4.4 1.6L12 15l-1.6-4.4L6 9l4.4-1.6zM19 14l.9 2.1L22 17l-2.1.9L19 20l-.9-2.1L16 17l2.1-.9z',
  alert: 'M12 3l9 16H3zM12 9v4M12 16.5v.5',
  clock: 'M12 4a8 8 0 100 16 8 8 0 000-16zM12 8v4l3 2',
  terminal: 'M4 5h16v14H4zM8 10l2 2-2 2M12 14h4',
  download: 'M12 4v10M8 11l4 4 4-4M5 19h14',
  refresh: 'M5 12a7 7 0 0112-5M19 12a7 7 0 01-12 5M5 7v5h5M19 17v-5h-5',
  risk: 'M12 3l9 16H3zM12 9v4M12 16.5v.5',
  settings: 'M12 8a4 4 0 100 8 4 4 0 000-8zM4 12h2M18 12h2M12 4v2M12 18v2',
  arrow: 'M5 12h14M13 6l6 6-6 6',
  sun: 'M12 8a4 4 0 100 8 4 4 0 000-8zM12 3v2M12 19v2M3 12h2M19 12h2M5.6 5.6l1.4 1.4M17 17l1.4 1.4M18.4 5.6L17 7M7 17l-1.4 1.4',
  moon: 'M20 14.5A8 8 0 119.5 4a6.5 6.5 0 0010.5 10.5z',
  more: 'M5 12h.5M12 12h.5M19 12h.5',
}

export function Icon({
  name,
  size = 20,
}: {
  name: IconName
  size?: number
}) {
  return (
    <svg
      className="icon"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.8"
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
    >
      <path d={paths[name]} />
    </svg>
  )
}
