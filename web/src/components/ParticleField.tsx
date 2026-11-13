import { useEffect, useRef } from 'react'

interface Particle {
  x: number
  y: number
  vx: number
  vy: number
  r: number
  tone: number
}

const TONES_LIGHT = [
  [64, 92, 134],
  [86, 116, 160],
  [110, 140, 184],
  [138, 164, 200],
]
const TONES_DARK = [
  [116, 150, 200],
  [140, 175, 226],
  [168, 198, 238],
  [192, 216, 248],
]

// ParticleField 是全屏沉浸式粒子层：缓慢漂浮、鼠标牵引、近邻连线。
// 以面积为基准控制数量，页面隐藏时自动暂停，遵循系统"减少动态效果"偏好。
export function ParticleField({ theme }: { theme: 'light' | 'dark' }) {
  const ref = useRef<HTMLCanvasElement | null>(null)

  useEffect(() => {
    const canvas = ref.current
    if (!canvas) return
    const ctx = canvas.getContext('2d')
    if (!ctx) return

    const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    let width = 0
    let height = 0
    let dpr = 1
    let particles: Particle[] = []
    let sprites: HTMLCanvasElement[] = []
    let linkDist = 150
    let raf = 0
    let running = true
    const mouse = { x: -9999, y: -9999, active: false }

    const build = () => {
      const tones = theme === 'dark' ? TONES_DARK : TONES_LIGHT
      width = window.innerWidth
      height = window.innerHeight
      dpr = Math.min(window.devicePixelRatio || 1, 1.6)
      canvas.width = Math.floor(width * dpr)
      canvas.height = Math.floor(height * dpr)
      ctx.setTransform(dpr, 0, 0, dpr, 0, 0)

      const narrow = width < 760
      const count = narrow
        ? Math.min(34, Math.max(16, Math.floor((width * height) / 30000)))
        : Math.min(76, Math.max(24, Math.floor((width * height) / 26000)))
      linkDist = narrow ? 112 : 140

      // 预渲染光晕 sprite，避免每帧为每颗粒子创建径向渐变。
      sprites = tones.map(([r, g, b]) => {
        const s = document.createElement('canvas')
        s.width = 64
        s.height = 64
        const sc = s.getContext('2d')
        if (sc) {
          const grad = sc.createRadialGradient(32, 32, 0, 32, 32, 32)
          grad.addColorStop(0, `rgba(${r},${g},${b},0.36)`)
          grad.addColorStop(1, `rgba(${r},${g},${b},0)`)
          sc.fillStyle = grad
          sc.fillRect(0, 0, 64, 64)
        }
        return s
      })
      particles = Array.from({ length: count }, () => ({
        x: Math.random() * width,
        y: Math.random() * height,
        vx: (Math.random() - 0.5) * 0.22,
        vy: (Math.random() - 0.5) * 0.22,
        r: Math.random() * 2.1 + 0.9,
        tone: Math.floor(Math.random() * tones.length),
      }))
    }
    build()

    const draw = () => {
      ctx.clearRect(0, 0, width, height)
      const alphaBase = theme === 'dark' ? 1 : 0.75

      for (const p of particles) {
        if (!reduced) {
          p.x += p.vx
          p.y += p.vy
        }
        if (mouse.active) {
          const dx = mouse.x - p.x
          const dy = mouse.y - p.y
          const d2 = dx * dx + dy * dy
          const range = 190
          if (d2 < range * range && d2 > 1) {
            // 轻微牵引，靠近指针时更明显，但不吞掉粒子。
            const f = (1 - Math.sqrt(d2) / range) * 0.012
            p.vx += dx * f
            p.vy += dy * f
          }
        }
        // 速度阻尼，保持缓慢漂浮。
        p.vx *= 0.995
        p.vy *= 0.995
        if (p.x < -20) p.x = width + 20
        if (p.x > width + 20) p.x = -20
        if (p.y < -20) p.y = height + 20
        if (p.y > height + 20) p.y = -20

        const [r, g, b] = (theme === 'dark' ? TONES_DARK : TONES_LIGHT)[p.tone]
        const glow = p.r * 10
        ctx.globalAlpha = alphaBase
        const sprite = sprites[p.tone]
        if (sprite) ctx.drawImage(sprite, p.x - glow / 2, p.y - glow / 2, glow, glow)
        ctx.globalAlpha = 1

        ctx.fillStyle = `rgba(${r},${g},${b},${0.82 * alphaBase})`
        ctx.beginPath()
        ctx.arc(p.x, p.y, p.r, 0, Math.PI * 2)
        ctx.fill()
      }

      // 近邻连线，形成轻盈的网状结构。
      for (let i = 0; i < particles.length; i++) {
        for (let j = i + 1; j < particles.length; j++) {
          const a = particles[i]
          const b = particles[j]
          const dx = a.x - b.x
          const dy = a.y - b.y
          const d2 = dx * dx + dy * dy
          if (d2 < linkDist * linkDist) {
            const alpha = (1 - Math.sqrt(d2) / linkDist) * 0.13 * alphaBase
            ctx.strokeStyle =
              theme === 'dark'
                ? `rgba(168,198,238,${alpha})`
                : `rgba(66,96,142,${alpha})`
            ctx.lineWidth = 0.7
            ctx.beginPath()
            ctx.moveTo(a.x, a.y)
            ctx.lineTo(b.x, b.y)
            ctx.stroke()
          }
        }
      }

      if (running) raf = requestAnimationFrame(draw)
    }
    raf = requestAnimationFrame(draw)

    const onMove = (e: PointerEvent) => {
      mouse.x = e.clientX
      mouse.y = e.clientY
      mouse.active = true
    }
    const onLeave = () => {
      mouse.active = false
    }
    const onResize = () => build()
    const onVisibility = () => {
      running = document.visibilityState === 'visible'
      if (running) {
        cancelAnimationFrame(raf)
        raf = requestAnimationFrame(draw)
      }
    }

    window.addEventListener('pointermove', onMove, { passive: true })
    window.addEventListener('pointerleave', onLeave)
    window.addEventListener('resize', onResize)
    document.addEventListener('visibilitychange', onVisibility)

    return () => {
      running = false
      cancelAnimationFrame(raf)
      window.removeEventListener('pointermove', onMove)
      window.removeEventListener('pointerleave', onLeave)
      window.removeEventListener('resize', onResize)
      document.removeEventListener('visibilitychange', onVisibility)
    }
  }, [theme])

  return <canvas ref={ref} className="particle-field" aria-hidden="true" />
}
