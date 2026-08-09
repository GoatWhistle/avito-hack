import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useActivateReward, useRewardTrack } from '#/features/rewards/hooks'
import {
  activationErrorKey,
  loadErrorKey,
  trackSummary,
} from '#/features/rewards/lib'
import { RewardLevelCard } from './RewardLevelCard'
import { RewardTrackSummary } from './RewardTrackSummary'
import { EmptyState, ErrorState, LoadingState } from './StateViews'

export function RewardTrack() {
  const { t } = useTranslation('rewards')
  const { data, isPending, isError, error, refetch } = useRewardTrack()
  const activate = useActivateReward()
  const [activeId, setActiveId] = useState<string | null>(null)
  const [codes, setCodes] = useState<Record<string, string>>({})
  const [errors, setErrors] = useState<Record<string, string>>({})

  if (isPending) return <LoadingState label={t('loading')} />

  if (isError) {
    return (
      <ErrorState
        message={t(loadErrorKey(error) as 'errors.loadFailed')}
        onRetry={() => void refetch()}
      />
    )
  }

  if (data.length === 0) {
    return (
      <EmptyState title={t('empty.catalog')} hint={t('empty.catalogHint')} />
    )
  }

  const summary = trackSummary(data)
  const percent = Math.round((summary.claimed / summary.total) * 100)

  const handleActivate = (rewardId: string) => {
    setActiveId(rewardId)
    setErrors((prev) => ({ ...prev, [rewardId]: '' }))

    activate.mutate(rewardId, {
      onSuccess: (result) => {
        setCodes((prev) => ({ ...prev, [rewardId]: result.code }))
      },
      onError: (mutationError) => {
        setErrors((prev) => ({
          ...prev,
          [rewardId]: t(activationErrorKey(mutationError) as 'errors.unknown'),
        }))
      },
      onSettled: () => setActiveId(null),
    })
  }

  return (
    <section
      aria-labelledby="reward-track-title"
      className="flex flex-col gap-4"
    >
      <RewardTrackSummary
        claimed={summary.claimed}
        total={summary.total}
        percent={percent}
      />

      <ul className="flex flex-col gap-2.5">
        {data.map((entry) => (
          <RewardLevelCard
            key={entry.reward.id}
            entry={entry}
            onActivate={handleActivate}
            isActivating={activeId === entry.reward.id && activate.isPending}
            code={codes[entry.reward.id]}
            error={errors[entry.reward.id]}
          />
        ))}
      </ul>
    </section>
  )
}
