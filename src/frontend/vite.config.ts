import { reactRouter } from '@react-router/dev/vite'
import tailwindcss from '@tailwindcss/vite'
import { devtools } from '@tanstack/devtools-vite'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [devtools(), tailwindcss(), reactRouter()],
  resolve: {
    tsconfigPaths: true,
  },
})
