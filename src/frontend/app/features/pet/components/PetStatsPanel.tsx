import { useTranslation } from 'react-i18next'
import type { StatView } from '#/features/pet/lib'
import { FeedButton } from './FeedButton'
import { StatMeter } from './StatMeter'
import { StreakCard } from './StreakCard'

export interface PetStatsPanelFeed {
  run: () => void
  isPending: boolean
  error: string | null
  availableAt?: string | null
}

export interface PetStatsPanelProps {
  stats: StatView[]
  streakDays: number
  freezes: number
  streakAtRisk: boolean
  feed?: PetStatsPanelFeed
}

export function PetStatsPanel({
  stats,
  streakDays,
  freezes,
  streakAtRisk,
  feed,
}: PetStatsPanelProps) {
  const { t } = useTranslation('pet')

  return (
    <div className="flex flex-col gap-4">
      <section
        aria-label={t('stats.groupLabel')}
        className="flex flex-col gap-4 rounded-xl bg-card px-4 py-4 ring-1 ring-foreground/10"
      >
        <h2 className="text-xs font-medium tracking-wide text-muted-foreground uppercase">
          {t('stats.groupLabel')}
        </h2>

        {stats.map((stat) => (
          <StatMeter
            key={stat.key}
            statKey={stat.key}
            value={stat.value}
            tone={stat.tone}
            action={
              stat.key === 'satiety' && feed !== undefined ? (
                <FeedButton
                  value={stat.value}
                  isPending={feed.isPending}
                  error={feed.error}
                  onFeed={feed.run}
                  availableAt={feed.availableAt}
                />
              ) : undefined
            }
          />
        ))}
      </section>

      <StreakCard days={streakDays} freezes={freezes} atRisk={streakAtRisk} />
    </div>
  )
}
