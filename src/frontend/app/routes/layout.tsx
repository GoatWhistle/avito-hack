import { Activity } from '#/features/tamagotchi/components/Activity'
import { LeftSidebar } from '#/features/tamagotchi/components/LeftSidebar'
import { SidebarInset, SidebarProvider } from '#/shared/components/ui/sidebar'
import { Outlet } from 'react-router'

export default function Layout() {
  return (
    <SidebarProvider>
      <LeftSidebar />
      <Activity />
      <SidebarInset>
        <Outlet />
      </SidebarInset>
    </SidebarProvider>
  )
}
