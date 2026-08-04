import type { AxiosInstance } from 'axios';

import { i18n } from '@/shared/i18n';

export function attachLocale(client: AxiosInstance): void {
  client.interceptors.request.use((config) => {
    config.headers['Accept-Language'] = i18n.language;

    return config;
  });
}
