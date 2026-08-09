import type { ReactNode } from 'react'
import { usePetScreen, type UsePetScreenOptions } from '#/features/pet/hooks'
import { CelebrationBanner } from '../CelebrationBanner'
import { PetActions } from '../PetActions'
import { PetScreenError, PetScreenSkeleton } from '../PetScreenStates'
import { DashboardActivity } from './DashboardActivity'
import { DashboardPanel } from './DashboardPanel'
import { DashboardSidebar } from './DashboardSidebar'
import { DashboardTopBar } from './DashboardTopBar'

export interface PetDashboardProps extends UsePetScreenOptions {
  children: ReactNode
}

export function PetDashboard({ children, ...options }: PetDashboardProps) {
  const screen = usePetScreen(options)
  const view = screen.view

  return (
    <div className="mx-auto flex w-full max-w-[100rem] flex-col gap-4 px-3 py-4 sm:px-4">
      <DashboardTopBar />

      {screen.celebration.banner !== null && (
        <CelebrationBanner
          banner={screen.celebration.banner}
          onDismiss={screen.celebration.dismissBanner}
        />
      )}

      {screen.isLoading && <PetScreenSkeleton />}

      {!screen.isLoading && screen.loadError !== null && (
        <PetScreenError
          message={screen.loadError}
          onRetry={screen.refetch}
          isRetrying={screen.isRetrying}
        />
      )}

      {view !== null && (
        <div className="isolate flex flex-col gap-4 lg:flex-row lg:items-start lg:gap-5">
          <DashboardSidebar
            pet={view.pet}
            progress={view.progress}
            feed={screen.feed}
          />

          <DashboardActivity
            pet={view.pet}
            emotion={screen.celebration.emotion}
            xpToasts={screen.celebration.xpToasts}
            checkedInToday={!screen.canCheckIn}
            onStroke={screen.stroke.run}
            onEmotionEnd={screen.celebration.clearEmotion}
          >
            <PetActions strokeError={screen.stroke.error} />
          </DashboardActivity>

          <DashboardPanel>{children}</DashboardPanel>
        </div>
      )}
    </div>
  )
}
