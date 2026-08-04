import axios, { type AxiosInstance } from 'axios';

import { ApiError, isErrorCode, type ApiErrorBody } from '../types';

const DEFAULT_STATUS = 0;

export function attachErrorMapper(client: AxiosInstance): void {
  client.interceptors.response.use(
    (response) => response,
    (error: unknown) => Promise.reject(toApiError(error)),
  );
}

export function toApiError(error: unknown): ApiError {
  if (error instanceof ApiError) {
    return error;
  }

  if (!axios.isAxiosError<ApiErrorBody>(error)) {
    return new ApiError({
      code: 'internal_error',
      message: 'unexpected error',
      status: DEFAULT_STATUS,
    });
  }

  const status = error.response?.status ?? DEFAULT_STATUS;
  const body = error.response?.data.error;

  if (body === undefined) {
    return new ApiError({
      code: 'network_error',
      message: error.message,
      status,
    });
  }

  return new ApiError({
    code: isErrorCode(body.code) ? body.code : 'internal_error',
    message: body.message,
    status,
    ...(body.field === undefined ? {} : { field: body.field }),
    ...(body.request_id === undefined ? {} : { requestId: body.request_id }),
  });
}
