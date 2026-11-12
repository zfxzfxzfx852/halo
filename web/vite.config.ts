import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// 开发时把后端 API 与 WebSocket 代理到本地 halo 服务。
export default defineConfig({
  plugins: [react()],
  base: '/console/',
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:8787',
        changeOrigin: true,
        ws: true,
      },
    },
  },
  build: {
    outDir: 'dist',
    emptyOutDir: true,
  },
})
