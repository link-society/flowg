import path from 'path'

import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from 'tailwindcss'

// https://vitejs.dev/config/
export default defineConfig({
  base: './',
  plugins: [
    react(),
  ],
  css: {
    postcss: {
      plugins: [
        tailwindcss,
      ],
    },
  },
  resolve: {
    alias: {
      '@materializecss/materialize/style': 'node_modules/@materializecss/materialize/dist/css/materialize.css',
      '@': path.resolve(__dirname, './src'),
    },
  },
})
