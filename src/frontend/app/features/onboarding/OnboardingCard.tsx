import { useCallback, useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { X } from 'lucide-react'
import {
  Button,
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '#/components/ui'
import { PetAvatar } from '#/components/pet-avatar'
import { useProgress } from '#/features/progress'
import { HatchCelebration } from './HatchCelebration'
import { OnboardingChecklist, type ChecklistStep } from './OnboardingChecklist'
import {
  readOnboardingState,
  writeOnboardingState,
  type OnboardingState,
} from './onboarding-state'
import { useHatchEvent } from './useHatchEvent'

interface OnboardingCardProps {
  persistent?: boolean
}

export function OnboardingCard({ persistent = false }: OnboardingCardProps) {
  const { t } = useTranslation('common')
  const { data: progress } = useProgress()
  const [state, setState] = useState<OnboardingState | null>(null)

  useEffect(() => {
    setState(readOnboardingState())
  }, [])

  const alreadyHatched = (progress?.xp ?? 0) > 0 || (progress?.level ?? 1) > 1
  const liveHatched = useHatchEvent(state !== null && !state.hatched)
  const hatched = Boolean(state?.hatched) || liveHatched || alreadyHatched

  const update = useCallback((next: Partial<OnboardingState>) => {
    setState((current) => {
      const base = current ?? readOnboardingState()
      const merged = { ...base, ...next }
      writeOnboardingState(merged)

      return merged
    })
  }, [])

  useEffect(() => {
    if (liveHatched && state && !state.hatched) update({ hatched: true })
  }, [liveHatched, state, update])

  if (!state) return null
  if (state.dismissed && !persistent) return null

  const showCelebration = hatched && !state.celebrated
  const streak = progress?.currentStreak ?? 0

  const steps: ChecklistStep[] = [
    { key: 'publish', done: hatched, to: '/items/new' },
    { key: 'comeback', done: streak > 1 },
    { key: 'meet', done: hatched, to: '/pet' },
  ]

  return (
    <>
      <Card className="relative" data-testid="onboarding-card">
        <CardHeader>
          <CardTitle>{t('onboarding.title')}</CardTitle>
          {!persistent && (
            <Button
              variant="ghost"
              size="icon-xs"
              onClick={() => update({ dismissed: true })}
              aria-label={t('actions.close')}
              className="absolute end-2 top-2"
            >
              <X className="size-3.5" aria-hidden="true" />
            </Button>
          )}
        </CardHeader>

        <CardContent className="flex flex-col gap-4 sm:flex-row sm:items-start">
          <div className="flex shrink-0 justify-center sm:justify-start">
            <PetAvatar
              stage={hatched ? 'baby' : 'egg'}
              satiety={80}
              happiness={80}
              energy={80}
              size="sm"
            />
          </div>

          <div className="min-w-0 flex-1">
            <p className="text-sm text-muted-foreground">
              {t('onboarding.intro')}
            </p>
            <div className="mt-3">
              <OnboardingChecklist steps={steps} />
            </div>
          </div>
        </CardContent>
      </Card>

      {showCelebration && (
        <HatchCelebration onClose={() => update({ celebrated: true })} />
      )}
    </>
  )
}
