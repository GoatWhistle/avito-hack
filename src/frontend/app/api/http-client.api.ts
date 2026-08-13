import axios, {
  AxiosError,
  type CreateAxiosDefaults,
  type InternalAxiosRequestConfig,
} from 'axios'
import { ApiError, apiErrorFromEnvelope } from './api-error'
import { isAuthPath, SessionRefresher } from './session-refresh'
import { clearToken, getToken } from './token-store'

const rawBase = import.meta.env.VITE_API_URL || ''

export const apiBaseUrl = `${rawBase.replace(/\/+$/, '')}/api/v1`

const config: CreateAxiosDefaults = {
  baseURL: apiBaseUrl,
  headers: {
    'Content-Type': 'application/json',
  },
}

export const httpClient = axios.create(config)

const refreshClient = axios.create(config)

const refresher = new SessionRefresher(refreshClient)

type RetriableConfig = InternalAxiosRequestConfig & {
  retriedAfterRefresh?: boolean
}

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

const shouldAttemptRefresh = (
  apiError: ApiError,
  request: RetriableConfig | undefined,
): request is RetriableConfig =>
  apiError.isUnauthorized &&
  request !== undefined &&
  request.retriedAfterRefresh !== true &&
  !isAuthPath(request.url) &&
  getToken() !== null

httpClient.interceptors.response.use(
  (response) => response,
  async (error: unknown) => {
    const apiError = toApiError(error)
    const request =
      error instanceof AxiosError
        ? (error.config as RetriableConfig | undefined)
        : undefined

    if (!shouldAttemptRefresh(apiError, request)) {
      return Promise.reject(apiError)
    }

    const token = await refresher.refresh()
    if (!token) {
      clearToken('expired')

      return Promise.reject(apiError)
    }

    request.retriedAfterRefresh = true
    request.headers.set('Authorization', `Bearer ${token}`)

    return httpClient.request(request)
  },
)
