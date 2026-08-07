import { isApiError } from '#/api'

export type Translate = (key: string) => string

const conflictHints = ['exists', 'занят', 'taken', 'duplicate']

const isConflictMessage = (message: string) =>
  conflictHints.some((hint) => message.toLowerCase().includes(hint))

export const toIssueMessage = (
  issue: { message?: string } | string | undefined,
  t: Translate,
): string | undefined => {
  const raw = typeof issue === 'string' ? issue : issue?.message
  if (!raw) return undefined

  return raw.includes(':') ? t(raw) : raw
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
    return error.message || t('auth:errors.invalidData')
  }

  return error.message || t(error.translationKey)
}
