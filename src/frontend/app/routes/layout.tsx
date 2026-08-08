import { Activity } from '#/features/tamagotchi/components/activity/Activity'
import { Dashboard } from '#/features/tamagotchi/components/dashboard/Dashboard'
import { RightSidebar } from '#/features/tamagotchi/components/right-sidebar/RightSidebar'
import { SidebarContent, SidebarProvider } from '#/shared/components/ui/sidebar'
import { useLayoutEffect, useRef } from 'react'
import { Outlet, useLocation } from 'react-router'

export default function Layout() {
  const { pathname } = useLocation()
  const containerRef = useRef<HTMLDivElement>(null)

  useLayoutEffect(() => {
    containerRef.current?.scrollTo(0, 0)
  }, [pathname])

  return (
    <SidebarProvider>
      <Dashboard />
      <Activity />
      <RightSidebar>
        <SidebarContent>
          <div ref={containerRef} className="h-full overflow-y-auto">
            <Outlet />
          </div>
        </SidebarContent>
      </RightSidebar>
    </SidebarProvider>
  )
}
