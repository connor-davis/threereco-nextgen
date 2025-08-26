import { type UserConfig } from '@hey-api/openapi-ts';

export default {
  input: {
    path: 'http://localhost:6173/api/api-spec',
  },
  output: {
    clean: true,
    format: 'prettier',
    tsConfigPath: 'tsconfig.json',
    path: 'src/api-client',
  },
  plugins: [
    '@tanstack/react-query',
    '@hey-api/client-fetch',
    '@hey-api/typescript',
    'zod',
  ],
} as UserConfig;
