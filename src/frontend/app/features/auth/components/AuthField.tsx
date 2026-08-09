import type { Ref } from 'react'
import { useTranslation } from 'react-i18next'
import {
  Field,
  FieldDescription,
  FieldError,
  FieldLabel,
  Input,
} from '#/components/ui'
import { toIssueMessage } from '#/features/auth/lib'

interface AuthFieldProps {
  name: string
  labelKey: string
  placeholderKey: string
  hintKey?: string
  type?: 'text' | 'email' | 'password'
  autoComplete?: string
  value: string
  disabled?: boolean
  touched: boolean
  issues: readonly unknown[]
  serverMessage?: string
  inputRef?: Ref<HTMLInputElement>
  onChange: (value: string) => void
  onBlur: () => void
}

export function AuthField({
  name,
  labelKey,
  placeholderKey,
  hintKey,
  type = 'text',
  autoComplete,
  value,
  disabled,
  touched,
  issues,
  serverMessage,
  inputRef,
  onChange,
  onBlur,
}: AuthFieldProps) {
  const { t } = useTranslation(['auth', 'validation'])
  const translate = (key: string, options?: Record<string, unknown>): string =>
    String(t(key as never, options as never))

  const localMessage = touched
    ? toIssueMessage(issues[0] as { message?: string } | undefined, translate)
    : undefined
  const message = localMessage ?? serverMessage
  const errorId = `${name}-error`
  const hintId = `${name}-hint`
  const describedBy = message ? errorId : hintKey ? hintId : undefined

  return (
    <Field data-invalid={message ? true : undefined}>
      <FieldLabel htmlFor={name}>{translate(labelKey)}</FieldLabel>
      <Input
        id={name}
        name={name}
        ref={inputRef}
        type={type}
        value={value}
        disabled={disabled}
        autoComplete={autoComplete}
        placeholder={translate(placeholderKey)}
        aria-invalid={message ? true : undefined}
        aria-describedby={describedBy}
        onChange={(event) => onChange(event.target.value)}
        onBlur={onBlur}
      />
      {message ? (
        <FieldError id={errorId}>{message}</FieldError>
      ) : (
        hintKey && (
          <FieldDescription id={hintId}>{translate(hintKey)}</FieldDescription>
        )
      )}
    </Field>
  )
}
