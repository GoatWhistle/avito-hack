import type { ReactNode } from 'react'
import { useTranslation } from 'react-i18next'
import { PetCharacter, type PetEmotion } from '#/components/pet-avatar'
import { buildAvatarLabels } from '#/features/pet/lib'
import type { XpToast } from '#/features/pet/hooks'
import type { Pet } from '#/features/pet/types'
import { XpToasts } from '../XpToasts'
import { StreakTracker } from './StreakTracker'

export interface DashboardActivityProps {
  pet: Pet
  emotion: PetEmotion | null
  xpToasts: XpToast[]
  checkedInToday?: boolean
  onStroke: () => void
  onEmotionEnd: () => void
  children: ReactNode
}

export function DashboardActivity({
  pet,
  emotion,
  xpToasts,
  checkedInToday = false,
  onStroke,
  onEmotionEnd,
  children,
}: DashboardActivityProps) {
  const { t } = useTranslation('pet')

  return (
    <div className="flex min-w-0 flex-1 flex-col items-center gap-4 lg:h-[calc(100dvh-9.5rem)]">
      <div className="flex w-full min-h-0 flex-1 flex-col items-center justify-center gap-4">
        <p
          data-testid="pet-speech"
          aria-live="polite"
          className="max-w-xs rounded-3xl bg-card px-4 py-3 text-center text-sm font-medium text-card-foreground ring-1 ring-foreground/10"
        >
          {t('dashboard.speech', {
            name: pet.name,
            stage: t(`stage.${pet.stage}`),
          })}
        </p>

        <div className="relative w-full max-w-[24rem] min-h-0 shrink">
          <XpToasts toasts={xpToasts} />
          <PetCharacter
            stage={pet.stage}
            emotion={emotion}
            satiety={pet.satiety}
            happiness={pet.happiness}
            energy={pet.energy}
            size="xl"
            labels={buildAvatarLabels(t, pet.name)}
            onStroke={onStroke}
            onEmotionEnd={onEmotionEnd}
            className="w-full"
          />
        </div>

        <div className="w-full max-w-sm">{children}</div>
      </div>

      <StreakTracker
        streakDays={pet.streak_days}
        checkedInToday={checkedInToday}
      />
    </div>
  )
}
