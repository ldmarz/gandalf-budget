import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    outDir: 'dist'
  },
  server: { // Add this server configuration
    proxy: {
      '/api': { // Proxy requests that start with /api
        target: 'http://localhost:8080', // Your Go backend address
        changeOrigin: true, // Necessary for virtual hosted sites
        // secure: false, // Uncomment if your backend is not HTTPS
        // No rewrite rule for now, assuming backend expects /api/...
      }
    }
  }
})
