import path from 'path'
import fs from 'fs'

import { defineConfig } from 'vite'
import { compression } from 'vite-plugin-compression2'
import react from '@vitejs/plugin-react'
import tailwindcss from 'tailwindcss'

// https://vitejs.dev/config/
export default defineConfig({
  base: './',
  plugins: [
    react(),
    compression({
      filename: '[path][base]',
      deleteOriginalAssets: true
    }),
  ],
  css: {
    postcss: {
      plugins: [
        tailwindcss,
      ],
    },
  },
  build: {
    chunkSizeWarningLimit: 1024,
    rollupOptions: {
      output: {
        manualChunks: () => 'bundle.js',
      },
    },
  },
  resolve: {
    alias: {
      '@materializecss/materialize/style': 'node_modules/@materializecss/materialize/dist/css/materialize.css',
      '@': path.resolve(__dirname, './src'),
    },
  },
  define: {
    'import.meta.env.FLOWG_VERSION': JSON.stringify(fs.readFileSync('../../VERSION.txt', 'utf8').trim()),
  },
})
