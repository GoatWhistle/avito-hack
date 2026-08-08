import { leftSidebarButtonItems } from '#/features/tamagotchi/components/dashboard/left-sidebar-button-item.data'
import {
  SidebarContent,
  SidebarGroup,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '#/shared/components/ui/sidebar'
import { NavLink, useLocation } from 'react-router'

export function DashboardContent() {
  const { pathname } = useLocation()

  return (
    <SidebarContent className="px-2">
      <SidebarGroup>
        <SidebarMenu className="gap-1.5">
          {leftSidebarButtonItems.map(item => (
            <SidebarMenuItem key={item.name}>
              <SidebarMenuButton
                title={item.name}
                size="lg"
                isActive={pathname === item.path}
                className="font-medium rounded-xl transition-all"
                render={
                  <NavLink to={item.path}>
                    <span className="text-xl font-light">{item.icon}</span>
                    <span>{item.name}</span>
                  </NavLink>
                }
              />
            </SidebarMenuItem>
          ))}
        </SidebarMenu>
      </SidebarGroup>
    </SidebarContent>
  )
}
