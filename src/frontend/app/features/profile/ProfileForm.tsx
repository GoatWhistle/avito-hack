import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Button, Field, FieldError, FieldLabel, Input } from '#/components/ui'
import { translateApiError } from '#/api'
import type { User } from '#/types'
import { useUpdateProfile } from './useUpdateProfile'

interface ProfileFormProps {
  user: User
}

export function ProfileForm({ user }: ProfileFormProps) {
  const { t } = useTranslation(['common', 'auth', 'validation', 'errors'])
  const [editing, setEditing] = useState(false)
  const [value, setValue] = useState(user.fullName)
  const { mutateAsync, isPending, error, reset } = useUpdateProfile()

  const trimmed = value.trim()
  const invalid = trimmed.length === 0

  const serverError = translateApiError(error, t) ?? null

  const cancel = () => {
    setValue(user.fullName)
    setEditing(false)
    reset()
  }

  const submit = async (event: React.FormEvent) => {
    event.preventDefault()
    if (invalid) return

    try {
      await mutateAsync(trimmed)
      setEditing(false)
    } catch {
      return
    }
  }

  if (!editing) {
    return (
      <div className="flex items-center justify-between gap-3">
        <div className="min-w-0">
          <p className="text-xs text-muted-foreground">
            {t('auth:fields.fullName')}
          </p>
          <p className="truncate text-sm font-medium">{user.fullName}</p>
        </div>
        <Button variant="outline" size="sm" onClick={() => setEditing(true)}>
          {t('common:actions.edit')}
        </Button>
      </div>
    )
  }

  return (
    <form onSubmit={submit} noValidate>
      <Field data-invalid={invalid || undefined}>
        <FieldLabel htmlFor="profile-full-name">
          {t('auth:fields.fullName')}
        </FieldLabel>
        <Input
          id="profile-full-name"
          name="fullName"
          value={value}
          autoFocus
          disabled={isPending}
          aria-invalid={invalid || Boolean(serverError) || undefined}
          onChange={(event) => setValue(event.target.value)}
        />
        {invalid && <FieldError>{t('validation:required')}</FieldError>}
        {serverError && !invalid && <FieldError>{serverError}</FieldError>}
      </Field>

      <div className="mt-3 flex gap-2">
        <Button type="submit" size="sm" disabled={invalid || isPending}>
          {isPending ? t('common:status.loading') : t('common:actions.save')}
        </Button>
        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={cancel}
          disabled={isPending}
        >
          {t('common:actions.cancel')}
        </Button>
      </div>
    </form>
  )
}
