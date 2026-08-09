import { useTranslation } from 'react-i18next'
import { Link, NavLink } from 'react-router'
import { cn } from '#/lib/utils'
import { dashboardNavItems } from './dashboard-nav'

export function DashboardTopBar() {
  const { t } = useTranslation(['pet', 'common'])

  return (
    <div className="flex flex-col gap-3">
      <div className="flex items-center justify-between gap-3">
        <Link
          to="/items"
          className="flex items-center gap-2 rounded-lg px-2 py-1.5 text-sm font-medium text-muted-foreground no-underline transition-colors hover:bg-muted hover:text-foreground focus-visible:ring-3 focus-visible:ring-ring/50 focus-visible:outline-none"
        >
          <span aria-hidden="true">←</span>
          {t('pet:dashboard.backToMarket')}
        </Link>
      </div>

      <nav
        aria-label={t('pet:dashboard.navLabel')}
        className="scrollbar-none -mx-3 overflow-x-auto px-3 lg:hidden"
      >
        <ul className="flex w-max gap-1.5">
          {dashboardNavItems.map((item) => (
            <li key={item.key}>
              <NavLink
                to={item.path}
                end={item.end}
                className={({ isActive }) =>
                  cn(
                    'flex items-center gap-1.5 rounded-full px-3 py-1.5 text-sm font-medium whitespace-nowrap no-underline transition-colors outline-none',
                    'focus-visible:ring-3 focus-visible:ring-ring/50',
                    isActive
                      ? 'bg-primary text-primary-foreground'
                      : 'bg-muted text-muted-foreground hover:text-foreground',
                  )
                }
              >
                <item.icon aria-hidden="true" className="size-4 shrink-0" />
                {t(`pet:dashboard.nav.${item.key}`)}
              </NavLink>
            </li>
          ))}
        </ul>
      </nav>
    </div>
  )
}
