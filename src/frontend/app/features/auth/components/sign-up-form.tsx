import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router'
import { Button, CardContent, CardFooter, FieldGroup } from '#/components/ui'
import { useFocusOnServerField, useSignUpForm } from '#/features/auth/forms'
import { toServerFieldName, toServerMessage } from '#/features/auth/lib'
import { useSession } from '#/features/auth/session'
import { requestPlatformTour } from '#/features/onboarding'
import { AuthField } from './AuthField'
import { AuthFormError } from './AuthFormError'

interface SignUpFormProps {
  redirectTo?: string
}

export function SignUpForm({ redirectTo = '/onboarding' }: SignUpFormProps) {
  const { t } = useTranslation(['auth', 'validation', 'errors', 'common'])
  const navigate = useNavigate()
  const { setUser } = useSession()

  const { form, isPending, error, register, focusFirstInvalid } = useSignUpForm(
    {
      onSuccess: (session) => {
        setUser(session.user)
        requestPlatformTour()
        void navigate(redirectTo, { replace: true })
      },
    },
  )

  const translate = (key: string, options?: Record<string, unknown>): string =>
    String(t(key as never, options as never))

  const serverError = toServerMessage(error, translate)
  const serverField = toServerFieldName(error)

  useFocusOnServerField(serverField, focusFirstInvalid)

  const messageFor = (name: string) =>
    serverField === name && serverError ? serverError : undefined

  return (
    <form
      noValidate
      data-testid="sign-up-form"
      onSubmit={(event) => {
        event.preventDefault()
        event.stopPropagation()
        void form.handleSubmit()
      }}
    >
      <CardContent>
        <FieldGroup>
          <form.Field name="fullName">
            {(field) => (
              <AuthField
                name="fullName"
                autoComplete="name"
                labelKey="auth:fields.fullName"
                placeholderKey="auth:placeholders.fullName"
                inputRef={register('fullName')}
                value={field.state.value}
                disabled={isPending}
                touched={field.state.meta.isTouched}
                issues={field.state.meta.errors}
                serverMessage={messageFor('fullName')}
                onChange={field.handleChange}
                onBlur={field.handleBlur}
              />
            )}
          </form.Field>

          <form.Field name="email">
            {(field) => (
              <AuthField
                name="email"
                type="email"
                autoComplete="email"
                labelKey="auth:fields.email"
                placeholderKey="auth:placeholders.email"
                inputRef={register('email')}
                value={field.state.value}
                disabled={isPending}
                touched={field.state.meta.isTouched}
                issues={field.state.meta.errors}
                serverMessage={messageFor('email')}
                onChange={field.handleChange}
                onBlur={field.handleBlur}
              />
            )}
          </form.Field>

          <form.Field name="password">
            {(field) => (
              <AuthField
                name="password"
                type="password"
                autoComplete="new-password"
                labelKey="auth:fields.password"
                placeholderKey="auth:placeholders.password"
                hintKey="auth:hints.passwordRule"
                inputRef={register('password')}
                value={field.state.value}
                disabled={isPending}
                touched={field.state.meta.isTouched}
                issues={field.state.meta.errors}
                serverMessage={messageFor('password')}
                onChange={field.handleChange}
                onBlur={field.handleBlur}
              />
            )}
          </form.Field>

          <AuthFormError message={serverField ? null : serverError} />
        </FieldGroup>
      </CardContent>

      <CardFooter>
        <Button type="submit" className="w-full" disabled={isPending}>
          {isPending ? t('common:status.loading') : t('auth:signUp.submit')}
        </Button>
      </CardFooter>
    </form>
  )
}
