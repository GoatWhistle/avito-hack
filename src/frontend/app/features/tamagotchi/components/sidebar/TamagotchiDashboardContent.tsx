import { SidebarGroup, SidebarMenu } from '#/shared/components/ui/sidebar'
import type { PropsWithChildren } from 'react'

export function TamagotchiDashboardContent({ children }: PropsWithChildren) {
  return (
    <SidebarGroup className="space-y-4 p-4">
      <SidebarMenu className="space-y-4">{children}</SidebarMenu>
    </SidebarGroup>
  )
}
