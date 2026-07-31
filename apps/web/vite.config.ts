import react from '@vitejs/plugin-react';
import { defineConfig } from 'vitest/config';

const coreTarget = process.env.CORVUS_CORE_URL ?? 'http://127.0.0.1:8765';

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': coreTarget,
      '/healthz': coreTarget,
    },
  },
  test: {
    environment: 'jsdom',
    setupFiles: './src/test-setup.ts',
  },
});
