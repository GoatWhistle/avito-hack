import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { Button } from '#/components/ui'
import { useRewardCatalog } from '#/features/rewards'
import type { RewardGroup, RewardProgress } from '#/features/rewards'

const VISIBLE_LIMIT = 3

export const nextRewards = (
  groups: RewardGroup[],
  limit = VISIBLE_LIMIT,
): RewardProgress[] =>
  groups
    .flatMap((group) => group.entries)
    .filter((entry) => !entry.reward.unlocked && !entry.reward.claimed)
    .sort((a, b) => b.ratio - a.ratio || a.remaining - b.remaining)
    .slice(0, limit)

export function PetRewardsPanel() {
  const { t } = useTranslation('pet')
  const { data, isPending, isError } = useRewardCatalog()

  const upcoming = useMemo(
    () => (data === undefined ? [] : nextRewards(data)),
    [data],
  )

  return (
    <section
      aria-labelledby="pet-rewards-title"
      data-testid="pet-rewards-panel"
      className="flex flex-col gap-3 rounded-xl bg-card px-4 py-4 ring-1 ring-foreground/10"
    >
      <div className="flex items-baseline justify-between gap-2">
        <h2
          id="pet-rewards-title"
          className="text-xs font-medium tracking-wide text-muted-foreground uppercase"
        >
          {t('rewards.title')}
        </h2>
        <Button
          size="sm"
          variant="ghost"
          render={<Link to="/pet/rewards">{t('rewards.all')}</Link>}
        />
      </div>

      {isPending && (
        <div aria-hidden="true" className="flex flex-col gap-3">
          {[0, 1, 2].map((row) => (
            <div key={row} className="flex flex-col gap-1.5">
              <div className="h-3 w-2/3 animate-pulse rounded-full bg-muted" />
              <div className="h-2 w-full animate-pulse rounded-full bg-muted" />
            </div>
          ))}
        </div>
      )}

      {!isPending && (isError || upcoming.length === 0) && (
        <p className="text-sm text-muted-foreground">
          {isError ? t('rewards.error') : t('rewards.empty')}
        </p>
      )}

      <ul className="flex flex-col gap-3">
        {upcoming.map((entry) => (
          <li key={entry.reward.id} className="flex flex-col gap-1.5">
            <div className="flex items-baseline justify-between gap-2">
              <span className="truncate text-sm font-medium text-foreground">
                {entry.reward.title}
              </span>
              <span className="shrink-0 font-mono text-xs tabular-nums text-muted-foreground">
                {t('rewards.progress', {
                  current: entry.current,
                  target: entry.target,
                })}
              </span>
            </div>

            <div
              role="meter"
              aria-valuenow={entry.percent}
              aria-valuemin={0}
              aria-valuemax={100}
              aria-label={t('rewards.meterLabel', {
                title: entry.reward.title,
                percent: entry.percent,
              })}
              className="h-1.5 w-full overflow-hidden rounded-full bg-muted"
            >
              <div
                className="h-full rounded-full bg-accent transition-[width] duration-500 ease-out"
                style={{ width: `${entry.percent}%` }}
              />
            </div>
          </li>
        ))}
      </ul>
    </section>
  )
}
