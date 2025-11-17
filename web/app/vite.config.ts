import path from 'node:path'
import fs from 'node:fs'

import { defineConfig } from 'vite'
import { compression } from 'vite-plugin-compression2'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vitejs.dev/config/
export default defineConfig({
  base: './',
  publicDir: 'static',
  plugins: [
    react(),
    tailwindcss(),
    compression({
      algorithms: ['gzip'],
      include: [
        /\.(html|css|js|map|ico|png)$/,
      ],
      filename: '[path][base]',
      deleteOriginalAssets: true
    }),
  ],
  build: {
    sourcemap: true,
    chunkSizeWarningLimit: 1024,
    rollupOptions: {
      output: {
        manualChunks: () => 'bundle.js',
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
})
