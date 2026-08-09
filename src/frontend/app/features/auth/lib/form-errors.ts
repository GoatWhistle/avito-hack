import { isApiError } from '#/api'
import {
  PASSWORD_MAX_LENGTH,
  PASSWORD_MIN_LENGTH,
} from '#/features/auth/schemas/credentials.schema'

export type Translate = (
  key: string,
  options?: Record<string, unknown>,
) => string

const conflictHints = ['exists', 'занят', 'taken', 'duplicate']

const isConflictMessage = (message: string) =>
  conflictHints.some((hint) => message.toLowerCase().includes(hint))

const isTooShort = (message: string) =>
  message.toLowerCase().includes('shorter')

const isTooLong = (message: string) => message.toLowerCase().includes('longer')

export const toIssueMessage = (
  issue: { message?: string } | string | undefined,
  t: Translate,
): string | undefined => {
  const raw = typeof issue === 'string' ? issue : issue?.message
  if (!raw) return undefined

  if (raw === 'validation:passwordTooShort') {
    return t(raw, { count: PASSWORD_MIN_LENGTH })
  }

  if (raw === 'validation:passwordTooLong') {
    return t(raw, { count: PASSWORD_MAX_LENGTH })
  }

  return raw.includes(':') ? t(raw) : raw
}

export const toServerFieldName = (error: unknown): string | null => {
  if (!isApiError(error)) return null
  if (error.kind === 'conflict') return 'email'
  if (error.kind !== 'validation_error' && error.kind !== 'bad_request') {
    return null
  }

  return error.field === 'full_name' ? 'fullName' : (error.field ?? null)
}

const toValidationMessage = (
  field: string | undefined,
  message: string,
  t: Translate,
): string => {
  if (field === 'password') {
    if (isTooShort(message)) {
      return t('validation:passwordTooShort', { count: PASSWORD_MIN_LENGTH })
    }
    if (isTooLong(message)) {
      return t('validation:passwordTooLong', { count: PASSWORD_MAX_LENGTH })
    }
  }

  if (field === 'email') return t('validation:email')

  return t('auth:errors.invalidData')
}

export const toServerMessage = (
  error: unknown,
  t: Translate,
): string | null => {
  if (!error) return null

  if (!isApiError(error)) return t('errors:unknown')

  if (error.kind === 'conflict' || isConflictMessage(error.message)) {
    return t('auth:errors.emailTaken')
  }

  if (error.kind === 'unauthorized') {
    return t('auth:errors.invalidCredentials')
  }

  if (error.kind === 'validation_error' || error.kind === 'bad_request') {
    return toValidationMessage(error.field, error.message, t)
  }

  return t(error.translationKey)
}
