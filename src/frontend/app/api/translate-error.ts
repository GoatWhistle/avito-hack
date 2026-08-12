import { isApiError } from './api-error'

export type TranslateFn = (
  key: string,
  options?: Record<string, unknown>,
) => string

export const translateApiError = (
  error: unknown,
  t: (key: never, options?: Record<string, unknown>) => string,
): string | undefined => {
  if (error === null || error === undefined) return undefined

  const translate = t as unknown as TranslateFn

  if (isApiError(error)) {
    return translate(error.translationKey, {
      defaultValue: translate('errors:unknown'),
    })
  }

  return translate('errors:unknown')
}
