import type { AxiosInstance } from 'axios';

import { tokenStorage } from '../token-storage';

type UnauthorizedHandler = () => void;

let onUnauthorized: UnauthorizedHandler | null = null;

export function setUnauthorizedHandler(handler: UnauthorizedHandler): void {
  onUnauthorized = handler;
}

export function attachAuth(client: AxiosInstance): void {
  client.interceptors.request.use((config) => {
    const token = tokenStorage.get();

    if (token !== null) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    config.headers['X-Request-Id'] = crypto.randomUUID();

    return config;
  });

  client.interceptors.response.use(
    (response) => response,
    (error: unknown) => {
      if (isUnauthorized(error)) {
        tokenStorage.clear();
        onUnauthorized?.();
      }

      return Promise.reject(error instanceof Error ? error : new Error(String(error)));
    },
  );
}

function isUnauthorized(error: unknown): boolean {
  return (
    typeof error === 'object' &&
    error !== null &&
    'response' in error &&
    typeof error.response === 'object' &&
    error.response !== null &&
    'status' in error.response &&
    error.response.status === 401
  );
}
