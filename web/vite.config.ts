import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    outDir: 'dist'
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        configure: (proxy, _options) => {
          proxy.on('error', (err, _req, _res) => {
            console.log('vite proxy error:', err.message); // Log error message
          });
          proxy.on('proxyReq', (proxyReq, req, _res) => {
            // Shorten the log by not logging the full proxyReq object itself
            console.log('vite proxy: Sending Request to Target:', req.method, req.url, '->', proxyReq.protocol + '//' + proxyReq.host + proxyReq.path);
          });
          proxy.on('proxyRes', (proxyRes, req, _res) => {
            console.log('vite proxy: Received Response from Target:', proxyRes.statusCode, req.url);
          });
        }
      }
    }
  }
})
