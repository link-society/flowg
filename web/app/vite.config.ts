import path from 'node:path'
import fs from 'node:fs'

import { defineConfig } from 'vite'
import { compression } from 'vite-plugin-compression2'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  base: './',
  publicDir: 'static',
  plugins: [
    react(),
    compression({
      algorithms: ['gzip'],
      include: [
        /\.(html|css|js|json|map|ico|png|woff2?)$/,
      ],
      skipIfLargerOrEqual: false,
      deleteOriginalAssets: true,
    }),
    {
      name: 'rewrite-assets-path',
      configureServer(serve) {
        serve.middlewares.use((req, _res, next) => {
          if (req.url?.startsWith('/web/assets/')) {
            req.url = req.url?.replace('/web/assets', '/assets')
          }

          next()
        })
      }
    },
    {
      name: 'go-template-substitution',
      transformIndexHtml: {
        order: 'pre',
        handler(html, ctx) {
          // Only substitute Go template placeholders when running the Vite dev
          // server. During a production build the placeholders must be left
          // intact so the Go backend can inject the real values at runtime.
          if (!ctx.server) {
            return html
          }

          return html
            .replace(/\{\{\s*\.MountPath\s*\}\}/g, '')
            .replace(/\{\{\s*\.FeatureFlags\.DemoMode\s*\}\}/g, 'false')
        },
      },
    },
  ],
  build: {
    sourcemap: true,
    chunkSizeWarningLimit: 4096,
    rollupOptions: {
      output: {
        manualChunks: () => 'bundle',
      },
    },
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  define: {
    'import.meta.env.FLOWG_VERSION': JSON.stringify(fs.readFileSync('../../VERSION.txt', 'utf8').trim()),
  },
  server: {
    open: '/web/',
    proxy: {
      '/api': {
        target: 'http://localhost:5080',
        changeOrigin: true
      },
    }
  },
})
