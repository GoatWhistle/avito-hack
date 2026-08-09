import type { PropsWithChildren } from 'react'
import { useTranslation } from 'react-i18next'
import { useLocation } from 'react-router'
import { useSession } from '#/features/auth/session'
import { PlatformTourGate } from '#/features/onboarding'
import { cn } from '#/lib/utils'
import { Header } from './Header'
import { MobileNav } from './MobileNav'

export function AppLayout({ children }: PropsWithChildren) {
  const { t } = useTranslation('common')
  const { isAuthenticated } = useSession()
  const { pathname } = useLocation()
  const isPetSection = pathname === '/pet' || pathname.startsWith('/pet/')

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
        className={cn(
          'w-full flex-1 pb-20',
          isPetSection
            ? 'px-0 py-0 lg:pb-0'
            : 'mx-auto max-w-content px-3 py-4 sm:px-4 sm:py-6 lg:pb-6',
        )}
      >
        {children}
      </div>

      {isAuthenticated && <MobileNav />}
      {isAuthenticated && <PlatformTourGate />}
    </div>
  )
}
