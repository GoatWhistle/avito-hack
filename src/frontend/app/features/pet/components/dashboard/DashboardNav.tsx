import { useTranslation } from 'react-i18next'
import { NavLink } from 'react-router'
import {
  SidebarContent,
  SidebarGroup,
  SidebarMenu,
  SidebarMenuItem,
} from '#/components/ui'
import { cn } from '#/lib/utils'
import { dashboardNavItems } from './dashboard-nav'

export function DashboardNav() {
  const { t } = useTranslation('pet')

  return (
    <SidebarContent className="px-2">
      <SidebarGroup>
        <nav aria-label={t('dashboard.navLabel')}>
          <SidebarMenu>
            {dashboardNavItems.map((item) => (
              <SidebarMenuItem key={item.key}>
                <NavLink
                  to={item.path}
                  end={item.end}
                  className={({ isActive }) =>
                    cn(
                      'flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium no-underline transition-colors outline-none',
                      'focus-visible:ring-3 focus-visible:ring-sidebar-ring/50',
                      isActive
                        ? 'bg-sidebar-primary text-sidebar-primary-foreground'
                        : 'text-sidebar-foreground hover:bg-sidebar-accent',
                    )
                  }
                >
                  <item.icon
                    aria-hidden="true"
                    className="size-[1.125rem] shrink-0"
                  />
                  <span className="truncate">
                    {t(`dashboard.nav.${item.key}`)}
                  </span>
                </NavLink>
              </SidebarMenuItem>
            ))}
          </SidebarMenu>
        </nav>
      </SidebarGroup>
    </SidebarContent>
  )
}
