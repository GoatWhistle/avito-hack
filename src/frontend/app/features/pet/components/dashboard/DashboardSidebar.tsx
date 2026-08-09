import { Sidebar } from '#/components/ui'
import type { LevelProgress } from '#/features/pet/lib'
import type { Pet } from '#/features/pet/types'
import { DashboardHeader } from './DashboardHeader'
import { DashboardNav } from './DashboardNav'
import { DashboardStats } from './DashboardStats'

export interface DashboardSidebarFeed {
  run: () => void
  isPending: boolean
  error: string | null
  availableAt: string | null
}

export interface DashboardSidebarProps {
  pet: Pet
  progress: LevelProgress
  feed: DashboardSidebarFeed
}

export function DashboardSidebar({
  pet,
  progress,
  feed,
}: DashboardSidebarProps) {
  return (
    <Sidebar
      side="left"
      className="hidden self-start [--sidebar-width:14rem] lg:sticky lg:top-[4.5rem] lg:flex lg:max-h-[calc(100dvh-9.5rem)] lg:overflow-y-auto xl:[--sidebar-width:16rem]"
    >
      <DashboardHeader pet={pet} progress={progress} />
      <DashboardNav />
      <DashboardStats pet={pet} feed={feed} />
    </Sidebar>
  )
}
