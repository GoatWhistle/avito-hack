import type { PropsWithChildren } from 'react'
import { useTranslation } from 'react-i18next'
import { Navigate } from 'react-router'
import { useSession } from '#/features/auth/session'

export function RequireAuth({ children }: PropsWithChildren) {
  const { t } = useTranslation('common')
  const { isAuthenticated, isLoading } = useSession()

  if (isLoading) {
    return (
      <p
        role="status"
        className="py-10 text-center text-sm text-muted-foreground"
      >
        {t('status.loading')}
      </p>
    )
  }

  if (!isAuthenticated) return <Navigate to="/sign-in" replace />

  return <>{children}</>
}
