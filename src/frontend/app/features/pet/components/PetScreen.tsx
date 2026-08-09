import { usePetScreen, type UsePetScreenOptions } from '#/features/pet/hooks'
import { CelebrationBanner } from './CelebrationBanner'
import { DailySummaryCard } from './DailySummaryCard'
import { NextStepHint } from './NextStepHint'
import { PetActions } from './PetActions'
import { PetHud } from './PetHud'
import { PetRewardsPanel } from './PetRewardsPanel'
import { PetScreenError, PetScreenSkeleton } from './PetScreenStates'
import { PetStage } from './PetStage'
import { PetStatsPanel } from './PetStatsPanel'

export type PetScreenProps = UsePetScreenOptions

export function PetScreen(options: PetScreenProps) {
  const screen = usePetScreen(options)
  const view = screen.view

  return (
    <main className="mx-auto flex w-full max-w-content flex-col gap-4 px-1 py-2 sm:gap-5">
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
        <>
          <PetHud
            pet={view.pet}
            progress={view.progress}
            streakAtRisk={view.streakAtRisk}
          />

          <div className="grid grid-cols-1 items-start gap-4 lg:grid-cols-[minmax(0,17rem)_minmax(0,1fr)_minmax(0,19rem)] lg:gap-5">
            <div className="order-2 flex flex-col gap-4 lg:order-1">
              <PetStatsPanel
                stats={view.stats}
                streakDays={view.pet.streak_days}
                freezes={view.pet.freezes}
                streakAtRisk={view.streakAtRisk}
                feed={screen.feed}
              />

              <NextStepHint step={view.step} />
            </div>

            <div className="order-1 lg:order-2">
              <PetStage
                pet={view.pet}
                speech={view.speech}
                emotion={screen.celebration.emotion}
                xpToasts={screen.celebration.xpToasts}
                onStroke={screen.stroke.run}
                onEmotionEnd={screen.celebration.clearEmotion}
              >
                <PetActions strokeError={screen.stroke.error} />
              </PetStage>
            </div>

            <div className="order-3 flex flex-col gap-4">
              <PetRewardsPanel />

              {screen.summary !== null && (
                <DailySummaryCard
                  summary={screen.summary}
                  petName={view.pet.name}
                  onDismiss={screen.dismissSummary}
                />
              )}
            </div>
          </div>
        </>
      )}
    </main>
  )
}
