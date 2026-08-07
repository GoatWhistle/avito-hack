import { useTranslation } from 'react-i18next'
import { usePetScreen, type UsePetScreenOptions } from '#/features/pet/hooks'
import { CelebrationBanner } from './CelebrationBanner'
import { DailySummaryCard } from './DailySummaryCard'
import { LevelProgress } from './LevelProgress'
import { NextStepHint } from './NextStepHint'
import { PetActions } from './PetActions'
import { PetHero } from './PetHero'
import {
  PetScreenEmpty,
  PetScreenError,
  PetScreenSkeleton,
} from './PetScreenStates'
import { StatMeter } from './StatMeter'
import { StreakCard } from './StreakCard'

export type PetScreenProps = UsePetScreenOptions

export function PetScreen(options: PetScreenProps) {
  const { t } = useTranslation('pet')
  const screen = usePetScreen(options)

  return (
    <main className="mx-auto flex w-full max-w-md flex-col gap-5 px-4 py-6 sm:max-w-lg sm:py-10">
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

      {screen.view !== null && !screen.view.pet.is_hatched && (
        <PetScreenEmpty
          onCheckIn={screen.checkIn.run}
          isCheckingIn={screen.checkIn.isPending}
        />
      )}

      {screen.view !== null && screen.view.pet.is_hatched && (
        <>
          <PetHero
            pet={screen.view.pet}
            emotion={screen.celebration.emotion}
            xpToasts={screen.celebration.xpToasts}
            onStroke={screen.stroke.run}
            onEmotionEnd={screen.celebration.clearEmotion}
          />

          <LevelProgress progress={screen.view.progress} />

          <StreakCard
            days={screen.view.pet.streak_days}
            freezes={screen.view.pet.freezes}
            atRisk={screen.view.streakAtRisk}
          />

          <PetActions
            canCheckIn={screen.canCheckIn}
            isStroking={screen.stroke.isPending}
            isCheckingIn={screen.checkIn.isPending}
            strokeError={screen.stroke.error}
            checkInError={screen.checkIn.error}
            onStroke={screen.stroke.run}
            onCheckIn={screen.checkIn.run}
          />

          <section
            aria-label={t('stats.groupLabel')}
            className="flex flex-col gap-4 rounded-xl bg-card px-4 py-4 ring-1 ring-foreground/10"
          >
            {screen.view.stats.map((stat) => (
              <StatMeter
                key={stat.key}
                statKey={stat.key}
                value={stat.value}
                tone={stat.tone}
              />
            ))}
          </section>

          <NextStepHint step={screen.view.step} />

          {screen.summary !== null && (
            <DailySummaryCard
              summary={screen.summary}
              petName={screen.view.pet.name}
              onDismiss={screen.dismissSummary}
            />
          )}
        </>
      )}
    </main>
  )
}
