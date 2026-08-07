import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router'
import { Button, CardContent, CardFooter, FieldGroup } from '#/components/ui'
import { useSignInForm } from '#/features/auth/forms'
import { toServerMessage } from '#/features/auth/lib'
import { useSession } from '#/features/auth/session'
import { AuthField } from './AuthField'
import { AuthFormError } from './AuthFormError'

interface SignInFormProps {
  redirectTo?: string
}

export function SignInForm({ redirectTo = '/pet' }: SignInFormProps) {
  const { t } = useTranslation(['auth', 'validation', 'errors', 'common'])
  const navigate = useNavigate()
  const { setUser } = useSession()

  const { form, isPending, error } = useSignInForm({
    onSuccess: (session) => {
      setUser(session.user)
      void navigate(redirectTo, { replace: true })
    },
  })

  const serverError = toServerMessage(error, (key) => t(key as never))

  return (
    <form
      noValidate
      data-testid="sign-in-form"
      onSubmit={(event) => {
        event.preventDefault()
        event.stopPropagation()
        void form.handleSubmit()
      }}
    >
      <CardContent>
        <FieldGroup>
          <form.Field name="email">
            {(field) => (
              <AuthField
                name="email"
                type="email"
                autoComplete="email"
                labelKey="auth:fields.email"
                placeholderKey="auth:placeholders.email"
                value={field.state.value}
                disabled={isPending}
                touched={field.state.meta.isTouched}
                issues={field.state.meta.errors}
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
                autoComplete="current-password"
                labelKey="auth:fields.password"
                placeholderKey="auth:placeholders.password"
                value={field.state.value}
                disabled={isPending}
                touched={field.state.meta.isTouched}
                issues={field.state.meta.errors}
                onChange={field.handleChange}
                onBlur={field.handleBlur}
              />
            )}
          </form.Field>

          <AuthFormError message={serverError} />
        </FieldGroup>
      </CardContent>

      <CardFooter>
        <Button type="submit" className="w-full" disabled={isPending}>
          {isPending ? t('common:status.loading') : t('auth:signIn.submit')}
        </Button>
      </CardFooter>
    </form>
  )
}
