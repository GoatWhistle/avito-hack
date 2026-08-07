import type { ReactNode } from 'react'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { Card, CardDescription, CardHeader, CardTitle } from '#/components/ui'

interface AuthScreenProps {
  mode: 'signIn' | 'signUp'
  children: ReactNode
}

const altLink = {
  signIn: {
    to: '/sign-up',
    promptKey: 'signIn.noAccount',
    ctaKey: 'signIn.goToSignUp',
  },
  signUp: {
    to: '/sign-in',
    promptKey: 'signUp.hasAccount',
    ctaKey: 'signUp.goToSignIn',
  },
} as const

export function AuthScreen({ mode, children }: AuthScreenProps) {
  const { t } = useTranslation('auth')
  const alt = altLink[mode]

  return (
    <main className="mx-auto flex w-full max-w-form flex-col gap-4 py-6 sm:py-10">
      <Card>
        <CardHeader>
          <CardTitle>{t(`${mode}.title`)}</CardTitle>
          <CardDescription>{t(`${mode}.subtitle`)}</CardDescription>
        </CardHeader>
        {children}
      </Card>

      <p className="text-center text-sm text-muted-foreground">
        {t(alt.promptKey)}{' '}
        <Link to={alt.to} className="font-medium text-primary">
          {t(alt.ctaKey)}
        </Link>
      </p>
    </main>
  )
}
