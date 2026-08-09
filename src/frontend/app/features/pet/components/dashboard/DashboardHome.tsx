import { useState } from 'react'
import { levelProgress, nextStep } from '#/features/pet/lib'
import { usePetQuery, useSummaryTodayQuery } from '#/features/pet/hooks'
import { canCheckInToday } from '#/features/pet/lib'
import { DailyQuestsPanel } from '#/features/quests'
import { DailySummaryCard } from '../DailySummaryCard'
import { NextStepHint } from '../NextStepHint'
import { PetRewardsPanel } from '../PetRewardsPanel'
import { NextLevelCard } from './NextLevelCard'

export function DashboardHome() {
  const { data: pet } = usePetQuery()
  const { data: summary } = useSummaryTodayQuery(pet !== undefined)
  const [dismissed, setDismissed] = useState(false)

  const step =
    pet === undefined
      ? null
      : nextStep({ pet, canCheckIn: canCheckInToday(pet) })
  const hint = step === 'maxLevel' ? null : step

  return (
    <div className="flex flex-col gap-4 p-3">
      {pet !== undefined && <NextLevelCard progress={levelProgress(pet)} />}

      {hint !== null && <NextStepHint step={hint} />}

      {pet !== undefined &&
        summary !== undefined &&
        summary !== null &&
        !dismissed && (
          <DailySummaryCard
            summary={summary}
            petName={pet.name}
            onDismiss={() => setDismissed(true)}
          />
        )}

      <DailyQuestsPanel />

      <PetRewardsPanel />
    </div>
  )
}
