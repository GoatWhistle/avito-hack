import { useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import {
  DESCRIPTION_MAX,
  MAX_PRICE_RUBLES,
  TITLE_MAX,
  TITLE_MIN,
} from '#/features/items/lib'

const params: Record<string, Record<string, number>> = {
  'validation:titleTooShort': { count: TITLE_MIN },
  'validation:maxLength': { count: TITLE_MAX },
  'validation:max': { value: MAX_PRICE_RUBLES },
}

const descriptionKey = 'validation:maxLength'

type LooseTranslate = (
  key: string,
  options: Record<string, number | string>,
) => string

export const useFieldErrors = () => {
  const { t } = useTranslation()
  const translateKey = t as unknown as LooseTranslate

  const translate = useCallback(
    (message: unknown, field?: string): string | undefined => {
      if (typeof message !== 'string' || message.length === 0) return undefined
      if (!message.startsWith('validation:')) return message

      const values =
        field === 'description' && message === descriptionKey
          ? { count: DESCRIPTION_MAX }
          : (params[message] ?? {})

      return translateKey(message, { ...values, defaultValue: message })
    },
    [translateKey],
  )

  return useCallback(
    (errors: readonly unknown[] | undefined, field?: string) => {
      if (!errors?.length) return []

      const messages = errors
        .map((error) =>
          translate(
            typeof error === 'object' && error !== null && 'message' in error
              ? (error as { message?: unknown }).message
              : error,
            field,
          ),
        )
        .filter((value): value is string => Boolean(value))

      return [...new Set(messages)].map((message) => ({ message }))
    },
    [translate],
  )
}
