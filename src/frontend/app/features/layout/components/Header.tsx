import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { Button } from '#/components/ui'
import { useSession } from '#/features/auth/session'
import { BrandMark } from './BrandMark'
import { DesktopNav } from './DesktopNav'
import { LocaleToggle } from './LocaleToggle'
import { MobileMenu } from './MobileMenu'
import { ThemeToggle } from './ThemeToggle'
import { UserMenu } from './UserMenu'

export function Header() {
  const { t } = useTranslation(['common', 'auth'])
  const { isAuthenticated } = useSession()

  return (
    <header className="sticky top-0 z-header border-b bg-background/95 backdrop-blur">
      <div className="mx-auto flex h-14 max-w-content items-center gap-2 px-3 sm:px-4">
        <MobileMenu />

        <Link
          to="/"
          className="flex min-w-0 items-center gap-2 no-underline"
          aria-label={t('common:app.name')}
        >
          <BrandMark className="h-7 w-auto shrink-0" />
        </Link>

        <div className="ms-1 hidden lg:block">
          <DesktopNav />
        </div>

        <div className="ms-auto flex items-center gap-1 sm:gap-2">
          <LocaleToggle />
          <ThemeToggle />
          {isAuthenticated ? (
            <UserMenu />
          ) : (
            <Button
              render={<Link to="/sign-in" />}
              size="sm"
              className="no-underline"
            >
              {t('auth:signIn.submit')}
            </Button>
          )}
        </div>
      </div>
    </header>
  )
}
