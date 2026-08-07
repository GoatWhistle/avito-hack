import axios from 'axios'
import { ApiError, apiErrorFromEnvelope, isApiError } from '#/api/api-error'

export const toApiError = (error: unknown): ApiError => {
  if (isApiError(error)) return error

  if (axios.isAxiosError(error)) {
    const status = error.response?.status ?? null

    if (status !== null) {
      const fromEnvelope = apiErrorFromEnvelope(error.response?.data, status)
      if (fromEnvelope) return fromEnvelope
    }

    if (error.code === 'ECONNABORTED' || error.code === 'ETIMEDOUT') {
      return new ApiError({ kind: 'timeout', message: error.message, status })
    }

    if (status === null) {
      return new ApiError({ kind: 'network', message: error.message, status })
    }

    return new ApiError({
      kind: status === 409 ? 'conflict' : 'internal_error',
      message: error.message,
      status,
    })
  }

  return new ApiError({
    kind: 'internal_error',
    message: error instanceof Error ? error.message : String(error),
  })
}

const messagePatterns: Array<[RegExp, string]> = [
  [/already been activated/i, 'errors.alreadyActivated'],
  [/has expired/i, 'errors.expired'],
  [/not been granted/i, 'errors.notGranted'],
  [/condition is not met/i, 'errors.locked'],
  [/cannot be activated/i, 'errors.notActivatable'],
]

export const activationErrorKey = (error: unknown): string => {
  const apiError = toApiError(error)

  for (const [pattern, key] of messagePatterns) {
    if (pattern.test(apiError.message)) return key
  }

  if (apiError.kind === 'conflict') return 'errors.alreadyActivated'
  if (apiError.kind === 'validation_error') return 'errors.notGranted'
  if (apiError.kind === 'unauthorized') return 'errors.unauthorized'
  if (apiError.kind === 'network' || apiError.kind === 'timeout') {
    return 'errors.network'
  }

  return 'errors.unknown'
}

export const loadErrorKey = (error: unknown): string => {
  const apiError = toApiError(error)

  if (apiError.kind === 'network' || apiError.kind === 'timeout') {
    return 'errors.network'
  }
  if (apiError.kind === 'unauthorized') return 'errors.unauthorized'

  return 'errors.loadFailed'
}
