import { useTranslation } from 'react-i18next'
import { NavLink } from 'react-router'
import { cn } from '#/lib/utils'
import { navItems } from '#/features/layout/nav-items'

export function DesktopNav() {
  const { t } = useTranslation('common')

  return (
    <nav aria-label={t('nav.primary')} className="hidden lg:block">
      <ul className="flex items-center gap-0.5">
        {navItems.map((item) => (
          <li key={item.key}>
            <NavLink
              to={item.to}
              className={({ isActive }) =>
                cn(
                  'flex items-center gap-1.5 rounded-lg px-2.5 py-1.5 text-sm font-medium no-underline transition-colors',
                  isActive
                    ? 'bg-primary-subtle text-primary-subtle-foreground'
                    : 'text-muted-foreground hover:bg-muted hover:text-foreground',
                )
              }
            >
              {({ isActive }) => (
                <>
                  <item.icon
                    className="size-4"
                    aria-hidden="true"
                    aria-current={isActive ? 'page' : undefined}
                  />
                  {t(`nav.${item.key}`)}
                </>
              )}
            </NavLink>
          </li>
        ))}
      </ul>
    </nav>
  )
}
