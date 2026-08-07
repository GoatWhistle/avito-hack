import { z } from 'zod'

export const apiErrorCodes = [
  'not_found',
  'forbidden',
  'unauthorized',
  'conflict',
  'validation_error',
  'bad_request',
  'internal_error',
] as const

export type ApiErrorCode = (typeof apiErrorCodes)[number]

export const apiErrorCodeSchema = z.enum(apiErrorCodes)

export const apiErrorEnvelopeSchema = z.object({
  error: z.object({
    code: z.string(),
    message: z.string(),
    field: z.string().optional(),
    request_id: z.string().optional(),
  }),
})

export type ApiErrorEnvelope = z.infer<typeof apiErrorEnvelopeSchema>

export type ApiErrorKind =
  ApiErrorCode | 'network' | 'timeout' | 'invalid_response'

export class ApiError extends Error {
  readonly kind: ApiErrorKind
  readonly status: number | null
  readonly field?: string
  readonly requestId?: string

  constructor(params: {
    kind: ApiErrorKind
    message: string
    status?: number | null
    field?: string
    requestId?: string
    cause?: unknown
  }) {
    super(params.message, { cause: params.cause })
    this.name = 'ApiError'
    this.kind = params.kind
    this.status = params.status ?? null
    this.field = params.field
    this.requestId = params.requestId
  }

  get isUnauthorized() {
    return this.kind === 'unauthorized'
  }

  get isValidation() {
    return this.kind === 'validation_error'
  }

  get translationKey() {
    switch (this.kind) {
      case 'network':
        return 'errors:network'
      case 'timeout':
        return 'errors:timeout'
      case 'invalid_response':
        return 'errors:invalidResponse'
      default:
        return `errors:code.${this.kind}`
    }
  }
}

export const isApiError = (value: unknown): value is ApiError =>
  value instanceof ApiError

const knownCode = (code: string): ApiErrorCode | null => {
  const parsed = apiErrorCodeSchema.safeParse(code)
  return parsed.success ? parsed.data : null
}

export const apiErrorFromEnvelope = (
  payload: unknown,
  status: number,
): ApiError | null => {
  const parsed = apiErrorEnvelopeSchema.safeParse(payload)
  if (!parsed.success) return null

  const { code, message, field, request_id: requestId } = parsed.data.error

  return new ApiError({
    kind: knownCode(code) ?? 'internal_error',
    message,
    status,
    field,
    requestId,
  })
}
