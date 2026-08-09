import { useLayoutEffect, useRef, type ReactNode } from 'react'
import { useTranslation } from 'react-i18next'
import { useLocation } from 'react-router'
import { Sidebar, SidebarContent } from '#/components/ui'

export interface DashboardPanelProps {
  children: ReactNode
}

export function DashboardPanel({ children }: DashboardPanelProps) {
  const { t } = useTranslation('pet')
  const { pathname } = useLocation()
  const containerRef = useRef<HTMLDivElement>(null)

  useLayoutEffect(() => {
    const container = containerRef.current
    if (container === null) return
    if (typeof container.scrollTo === 'function') {
      container.scrollTo(0, 0)

      return
    }
    container.scrollTop = 0
  }, [pathname])

  return (
    <Sidebar
      side="right"
      aria-label={t('dashboard.panelLabel')}
      className="w-full min-w-0 self-start [--sidebar-width:22rem] lg:sticky lg:top-[4.5rem] lg:max-h-[calc(100dvh-9.5rem)] lg:w-(--sidebar-width) lg:shrink-0 xl:[--sidebar-width:24rem]"
    >
      <SidebarContent className="scrollbar-thin min-h-0 px-0" ref={containerRef}>
        {children}
      </SidebarContent>
    </Sidebar>
  )
}
