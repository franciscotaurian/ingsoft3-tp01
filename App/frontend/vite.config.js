import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  test: {
    coverage: {
      provider: 'v8',
      // Los cuatro reporters: texto en consola, HTML navegable, lcov y el JSON
      // que el pipeline lee para armar el resumen en el Step Summary de GitHub.
      reporter: ['text', 'html', 'lcov', 'json-summary'],
      // Mide solo los módulos que tienen tests: la lógica pura (pedido.js)
      // y el cliente de API (client.js). 
      include: ['src/lib/pedido.js', 'src/api/client.js'],
      thresholds: { lines: 70, branches: 70 },
    },
  },
})
