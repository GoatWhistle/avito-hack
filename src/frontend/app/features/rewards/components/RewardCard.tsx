import { useTranslation } from 'react-i18next'
import { Card, CardContent } from '#/components/ui'
import { cn } from '#/lib/utils'
import {
  conditionLabel,
  remainingLabel,
  rewardText,
} from '#/features/rewards/lib'
import type { RewardProgress } from '#/features/rewards/types'
import { ProgressBar } from './ProgressBar'

type RewardCardProps = {
  entry: RewardProgress
  highlighted?: boolean
}

export function RewardCard({ entry, highlighted = false }: RewardCardProps) {
  const { t } = useTranslation('rewards')
  const { t: tCatalog } = useTranslation('catalog')
  const { reward } = entry
  const text = rewardText(tCatalog, reward.id, {
    title: reward.title,
    description: reward.description,
  })
  const isAvailable = reward.unlocked || reward.claimed
  const condition = conditionLabel(
    t,
    reward.condition_type,
    reward.condition_value,
  )
  const remaining = remainingLabel(t, entry)

  const tone = isAvailable
    ? 'complete'
    : entry.group === 'close'
      ? 'close'
      : 'default'

  const statusKey = reward.claimed
    ? 'status.claimed'
    : reward.unlocked
      ? 'status.unlocked'
      : 'status.locked'

  return (
    <Card
      size="sm"
      data-testid="reward-card"
      data-reward-id={reward.id}
      className={cn(
        'transition-shadow',
        highlighted && 'ring-2 ring-warning',
        !isAvailable && entry.group === 'far' && 'opacity-80',
      )}
    >
      <CardContent className="flex flex-col gap-3">
        <div className="flex items-start justify-between gap-3">
          <div className="min-w-0">
            <h3 className="line-clamp-2 font-medium text-balance text-foreground">
              {text.title}
            </h3>
            <p className="mt-0.5 line-clamp-2 text-sm text-muted-foreground">
              {text.description}
            </p>
          </div>
          <span
            className={cn(
              'shrink-0 rounded-full px-2 py-0.5 text-xs font-medium',
              isAvailable
                ? 'bg-success-subtle text-success-subtle-foreground'
                : 'bg-muted text-muted-foreground',
            )}
          >
            {t(statusKey as 'status.locked')}
          </span>
        </div>

        <div className="flex flex-col gap-1.5">
          <div className="flex items-baseline justify-between gap-2 text-sm">
            <span className="font-medium text-foreground">{condition}</span>
            <span className="tabular-nums text-muted-foreground">
              {t('progress', { current: entry.current, target: entry.target })}
            </span>
          </div>

          <ProgressBar
            value={entry.percent}
            tone={tone}
            label={t('progressLabel', {
              title: text.title,
              percent: entry.percent,
            })}
          />

          <p className="text-xs text-muted-foreground">
            {isAvailable ? t('readyToClaim') : remaining}
          </p>
        </div>
      </CardContent>
    </Card>
  )
}
