import { Info } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { useLocation, useNavigate } from 'react-router'
import { useSession } from '#/features/auth/session'
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
  Tabs,
  TabsContent,
  TabsIndicator,
  TabsList,
  TabsTrigger,
} from '#/components/ui'
import { SignInForm } from './sign-in.form'
import { SignUpForm } from './sign-up-form'
import { useAnimatedHeight } from './useAnimatedHeight'

import './auth-screen.css'

export type AuthTab = 'signIn' | 'signUp'

interface AuthScreenProps {
  mode: AuthTab
}

const tabPath: Record<AuthTab, string> = {
  signIn: '/sign-in',
  signUp: '/sign-up',
}

export function AuthScreen({ mode }: AuthScreenProps) {
  const { t } = useTranslation('auth')
  const navigate = useNavigate()
  const { sessionExpired } = useSession()
  const location = useLocation()
  const state = location.state as { from?: unknown } | null
  const redirectTo = typeof state?.from === 'string' ? state.from : undefined
  const { ref: panelsRef, style: panelsStyle } = useAnimatedHeight(mode)

  return (
    <main className="auth-screen flex min-h-svh items-center justify-center px-4 py-8">
      <Card className="auth-card w-full max-w-form">
        <CardHeader className="text-center">
          <CardTitle className="font-heading text-2xl font-semibold tracking-tight">
            {t(`${mode}.title`)}
          </CardTitle>
          <CardDescription className="text-base">
            {t(`${mode}.subtitle`)}
          </CardDescription>
        </CardHeader>

        {sessionExpired && (
          <div className="px-(--card-spacing) pb-1">
            <p
              role="status"
              data-testid="session-expired-notice"
              className="flex items-start gap-2 rounded-lg bg-accent-subtle px-3 py-2 text-sm text-accent-subtle-foreground"
            >
              <Info className="mt-0.5 size-4 shrink-0" aria-hidden="true" />
              <span className="min-w-0">{t('sessionExpired')}</span>
            </p>
          </div>
        )}

        <Tabs
          value={mode}
          onValueChange={(value) => {
            void navigate(tabPath[value as AuthTab], { replace: true })
          }}
        >
          <div className="px-(--card-spacing)">
            <TabsList className="grid w-full grid-cols-2">
              <TabsIndicator />
              <TabsTrigger value="signIn">{t('tabs.signIn')}</TabsTrigger>
              <TabsTrigger value="signUp">{t('tabs.signUp')}</TabsTrigger>
            </TabsList>
          </div>

          <div
            ref={panelsRef}
            style={panelsStyle}
            className="auth-card__panels"
          >
            <TabsContent value="signIn" className="auth-panel">
              <SignInForm redirectTo={redirectTo} />
            </TabsContent>

            <TabsContent value="signUp" className="auth-panel">
              <SignUpForm />
            </TabsContent>
          </div>
        </Tabs>
      </Card>
    </main>
  )
}
