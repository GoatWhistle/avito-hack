import axios, { AxiosError, type CreateAxiosDefaults } from 'axios'
import { ApiError, apiErrorFromEnvelope } from './api-error'
import { clearToken, getToken } from './token-store'

const rawBase = import.meta.env.VITE_API_URL ?? 'http://localhost:8080/'

export const apiBaseUrl = `${rawBase.replace(/\/+$/, '')}/api/v1`

const config: CreateAxiosDefaults = {
  baseURL: apiBaseUrl,
  headers: {
    'Content-Type': 'application/json',
  },
}

export const httpClient = axios.create(config)

httpClient.interceptors.request.use((request) => {
  const token = getToken()
  if (token) {
    request.headers.set('Authorization', `Bearer ${token}`)
  }
  return request
})

const toApiError = (error: unknown): ApiError => {
  if (error instanceof ApiError) return error

  if (error instanceof AxiosError) {
    if (error.code === AxiosError.ECONNABORTED || error.code === 'ETIMEDOUT') {
      return new ApiError({
        kind: 'timeout',
        message: error.message,
        cause: error,
      })
    }

    const response = error.response
    if (!response) {
      return new ApiError({
        kind: 'network',
        message: error.message,
        cause: error,
      })
    }

    const fromEnvelope = apiErrorFromEnvelope(response.data, response.status)
    if (fromEnvelope) return fromEnvelope

    return new ApiError({
      kind: response.status === 401 ? 'unauthorized' : 'internal_error',
      message: error.message,
      status: response.status,
      cause: error,
    })
  }

  return new ApiError({
    kind: 'internal_error',
    message: 'Unexpected error',
    cause: error,
  })
}

httpClient.interceptors.response.use(
  (response) => response,
  (error: unknown) => {
    const apiError = toApiError(error)
    if (apiError.isUnauthorized) {
      clearToken()
    }
    return Promise.reject(apiError)
  },
)
