/// <reference types="vitest" />
import path from 'node:path';

import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';

import { createPrettyLogger, prettyAccessLog } from './vite.logger';

export default defineConfig({
  plugins: [react(), prettyAccessLog()],
  customLogger: createPrettyLogger(),
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 5173,
    host: true,
    ...(process.env.VITE_USE_POLLING === 'true' ? { watch: { usePolling: true, interval: 400 } } : {}),
    ...(process.env.VITE_HMR_PORT ? { hmr: { clientPort: Number(process.env.VITE_HMR_PORT) } } : {}),
  },
  build: {
    outDir: 'dist',
    sourcemap: true,
    rollupOptions: {
      output: {
        manualChunks: {
          react: ['react', 'react-dom', 'react-router-dom'],
          antd: ['antd', '@ant-design/icons'],
          query: ['@tanstack/react-query', 'axios'],
        },
      },
    },
  },
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./vitest.setup.ts'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'lcov'],
      exclude: ['**/*.config.*', '**/index.ts', 'src/app/**'],
    },
  },
});
