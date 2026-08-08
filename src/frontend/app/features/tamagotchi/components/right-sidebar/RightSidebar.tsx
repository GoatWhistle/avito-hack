import { Sidebar } from '#/shared/components/ui/sidebar'
import type { PropsWithChildren } from 'react'

export function RightSidebar({ children }: PropsWithChildren) {
  return (
    <Sidebar
      side="right"
      variant="floating"
      className="[--sidebar-width:24rem]"
    >
      {children}
    </Sidebar>
  )
}
