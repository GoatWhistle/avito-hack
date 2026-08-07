import { useTranslation } from 'react-i18next'
import { NavLink } from 'react-router'
import { cn } from '#/lib/utils'
import { navItems } from '#/features/layout/nav-items'

const primaryItems = navItems.filter((item) => item.primary)

export function MobileNav() {
  const { t } = useTranslation('common')

  return (
    <nav
      aria-label={t('nav.primary')}
      className="fixed inset-x-0 bottom-0 z-sticky border-t bg-background/95 pb-[env(safe-area-inset-bottom)] backdrop-blur lg:hidden"
    >
      <ul className="mx-auto flex max-w-content items-stretch justify-around">
        {primaryItems.map((item) => (
          <li key={item.key} className="min-w-0 flex-1">
            <NavLink
              to={item.to}
              className={({ isActive }) =>
                cn(
                  'flex flex-col items-center gap-0.5 px-1 py-2 text-[0.6875rem] font-medium no-underline transition-colors',
                  isActive
                    ? 'text-primary'
                    : 'text-muted-foreground hover:text-foreground',
                )
              }
            >
              {({ isActive }) => (
                <>
                  <item.icon
                    className="size-5"
                    aria-hidden="true"
                    aria-current={isActive ? 'page' : undefined}
                  />
                  <span className="max-w-full truncate">
                    {t(`nav.${item.key}`)}
                  </span>
                </>
              )}
            </NavLink>
          </li>
        ))}
      </ul>
    </nav>
  )
}
