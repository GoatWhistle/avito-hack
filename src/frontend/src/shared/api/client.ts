import axios from 'axios';

import { API_BASE_URL } from '@/shared/config/env';

import { attachAuth } from './interceptors/auth';
import { attachErrorMapper } from './interceptors/error';
import { attachLocale } from './interceptors/locale';

const REQUEST_TIMEOUT_MS = 15_000;

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: REQUEST_TIMEOUT_MS,
  headers: { 'Content-Type': 'application/json' },
});

attachAuth(apiClient);
attachLocale(apiClient);
attachErrorMapper(apiClient);
