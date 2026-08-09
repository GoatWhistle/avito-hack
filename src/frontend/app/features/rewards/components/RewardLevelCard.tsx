import { useTranslation } from 'react-i18next'
import { Button, Card, CardContent } from '#/components/ui'
import { cn } from '#/lib/utils'
import { isActivatable, rewardText } from '#/features/rewards/lib'
import type { RewardTrackEntry } from '#/features/rewards/types'
import { RewardCodeReveal } from './RewardCodeReveal'
import {
  REWARD_CARD_TONE,
  REWARD_KIND_ICON,
  REWARD_KIND_TONE,
  REWARD_MARKER_TONE,
  REWARD_UNIT_KEY,
  rewardCardState,
} from './reward-card-style'

export interface RewardLevelCardProps {
  entry: RewardTrackEntry
  onActivate: (rewardId: string) => void
  isActivating?: boolean
  error?: string
  code?: string
}

export function RewardLevelCard({
  entry,
  onActivate,
  isActivating = false,
  error,
  code,
}: RewardLevelCardProps) {
  const { t } = useTranslation('rewards')
  const { t: tCatalog } = useTranslation('catalog')
  const { reward, state, level, current, target } = entry
  const text = rewardText(tCatalog, reward.id, {
    title: reward.title,
    description: reward.description,
  })
  const revealed = code || entry.code
  const canActivate = isActivatable(state)
  const cardState = rewardCardState(state)
  const unitKey = REWARD_UNIT_KEY[reward.condition_type] ?? 'track.unit.generic'

  return (
    <li>
      <Card
        size="sm"
        data-testid="reward-level-card"
        data-reward-id={reward.id}
        data-state={state}
        className={cn('relative transition-all', REWARD_CARD_TONE[cardState])}
      >
        {cardState !== 'locked' && (
          <span
            aria-hidden="true"
            className="pointer-events-none absolute -top-4 -right-4 size-20 rounded-full bg-reward-glow/15 blur-xl"
          />
        )}

        <CardContent className="flex flex-col gap-3">
          <div className="flex items-start justify-between gap-3">
            <div className="flex min-w-0 flex-1 items-start gap-3">
              <span
                aria-hidden="true"
                className={cn(
                  'mt-0.5 grid size-10 shrink-0 place-items-center rounded-xl font-mono text-sm font-bold tabular-nums',
                  REWARD_MARKER_TONE[cardState],
                )}
              >
                {cardState === 'completed' ? '✓' : level}
              </span>

              <div className="flex min-w-0 flex-1 flex-col gap-1">
                <div className="flex items-start gap-1.5">
                  <span
                    aria-hidden="true"
                    className={cn(
                      'shrink-0 text-xs leading-5',
                      REWARD_KIND_TONE[reward.kind] ?? 'text-muted-foreground',
                    )}
                  >
                    {cardState === 'locked'
                      ? '🔒'
                      : (REWARD_KIND_ICON[reward.kind] ?? '🎁')}
                  </span>
                  <h3 className="min-w-0 text-sm leading-tight font-bold wrap-anywhere text-balance text-foreground">
                    {text.title}
                  </h3>
                </div>
                <p className="line-clamp-2 text-xs leading-snug wrap-anywhere text-muted-foreground">
                  {text.description}
                </p>
              </div>
            </div>

            <div className="flex shrink-0 items-center self-center">
              {cardState === 'completed' && (
                <span className="rounded-full bg-success px-2.5 py-1 text-xs font-bold whitespace-nowrap text-success-foreground">
                  {t('status.claimed')}
                </span>
              )}

              {cardState === 'locked' && (
                <span className="rounded-full px-2.5 py-1 font-mono text-xs font-bold tabular-nums whitespace-nowrap text-muted-foreground ring-1 ring-border">
                  {t('track.progressShort', {
                    current,
                    target,
                    unit: t(unitKey as 'track.unit.level'),
                  })}
                </span>
              )}

              {cardState === 'current' && canActivate && (
                <Button
                  size="sm"
                  disabled={isActivating}
                  onClick={() => onActivate(reward.id)}
                >
                  {isActivating ? t('actions.activating') : t('actions.claim')}
                </Button>
              )}
            </div>
          </div>

          {revealed && <RewardCodeReveal code={revealed} />}

          {error && (
            <p role="alert" className="text-xs text-destructive">
              {error}
            </p>
          )}
        </CardContent>
      </Card>
    </li>
  )
}
