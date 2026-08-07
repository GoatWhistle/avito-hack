import type { PropsWithChildren } from 'react'
import { useTranslation } from 'react-i18next'
import { useSession } from '#/features/auth/session'
import { Header } from './Header'
import { MobileNav } from './MobileNav'

export function AppLayout({ children }: PropsWithChildren) {
  const { t } = useTranslation('common')
  const { isAuthenticated } = useSession()

  return (
    <div className="flex min-h-dvh flex-col">
      <a
        href="#main"
        className="sr-only focus:not-sr-only focus:absolute focus:start-2 focus:top-2 focus:z-tooltip focus:rounded-lg focus:bg-primary focus:px-3 focus:py-2 focus:text-sm focus:text-primary-foreground focus:no-underline"
      >
        {t('nav.skipToContent')}
      </a>

      <Header />

      <div
        id="main"
        tabIndex={-1}
        className="mx-auto w-full max-w-content flex-1 px-3 py-4 pb-20 sm:px-4 sm:py-6 lg:pb-6"
      >
        {children}
      </div>

      {isAuthenticated && <MobileNav />}
    </div>
  )
}
