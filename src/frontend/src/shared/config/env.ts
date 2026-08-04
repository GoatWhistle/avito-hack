import { z } from 'zod';

const envSchema = z.object({
  VITE_API_URL: z.string().url(),
  VITE_DEFAULT_LOCALE: z.enum(['ru', 'en']).default('ru'),
});

export type Env = z.infer<typeof envSchema>;

export const env: Env = envSchema.parse({
  VITE_API_URL: import.meta.env.VITE_API_URL ?? 'http://localhost:8080',
  VITE_DEFAULT_LOCALE: import.meta.env.VITE_DEFAULT_LOCALE ?? 'ru',
});

export const API_VERSION = 'v1';
export const API_BASE_URL = `${env.VITE_API_URL}/api/${API_VERSION}`;
