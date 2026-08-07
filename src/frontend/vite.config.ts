import { reactRouter } from '@react-router/dev/vite'
import tailwindcss from '@tailwindcss/vite'
import { defineConfig } from 'vite'
import { devtools } from '@tanstack/devtools-vite'

export default defineConfig({
  plugins: [devtools(), tailwindcss(), reactRouter()],
  resolve: {
    tsconfigPaths: true,
  },
  preview: {
    host: '127.0.0.1',
  },
})
