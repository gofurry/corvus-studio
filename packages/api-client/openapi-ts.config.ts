import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
  input: '../../apps/core/openapi/openapi.yaml',
  output: 'src/generated',
  plugins: ['@hey-api/client-fetch', '@hey-api/typescript', '@hey-api/sdk'],
});
