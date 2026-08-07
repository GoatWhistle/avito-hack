import { useTranslation } from 'react-i18next'
import { Field, FieldError, FieldLabel, Input } from '#/components/ui'
import { toIssueMessage } from '#/features/auth/lib'

interface AuthFieldProps {
  name: string
  labelKey: string
  placeholderKey: string
  type?: 'text' | 'email' | 'password'
  autoComplete?: string
  value: string
  disabled?: boolean
  touched: boolean
  issues: readonly unknown[]
  onChange: (value: string) => void
  onBlur: () => void
}

export function AuthField({
  name,
  labelKey,
  placeholderKey,
  type = 'text',
  autoComplete,
  value,
  disabled,
  touched,
  issues,
  onChange,
  onBlur,
}: AuthFieldProps) {
  const { t } = useTranslation(['auth', 'validation'])
  const translate = (key: string): string => t(key as never)

  const message = touched
    ? toIssueMessage(issues[0] as { message?: string } | undefined, translate)
    : undefined
  const errorId = `${name}-error`

  return (
    <Field data-invalid={message ? true : undefined}>
      <FieldLabel htmlFor={name}>{translate(labelKey)}</FieldLabel>
      <Input
        id={name}
        name={name}
        type={type}
        value={value}
        disabled={disabled}
        autoComplete={autoComplete}
        placeholder={translate(placeholderKey)}
        aria-invalid={message ? true : undefined}
        aria-describedby={message ? errorId : undefined}
        onChange={(event) => onChange(event.target.value)}
        onBlur={onBlur}
      />
      {message && <FieldError id={errorId}>{message}</FieldError>}
    </Field>
  )
}
