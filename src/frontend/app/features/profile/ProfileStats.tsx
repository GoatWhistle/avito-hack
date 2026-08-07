import { useTranslation } from 'react-i18next'
import { useProgress } from '#/features/progress'
import { useItemCount } from './useItemCount'

interface StatTileProps {
  label: string
  value: string
  loading: boolean
}

function StatTile({ label, value, loading }: StatTileProps) {
  return (
    <div className="rounded-lg bg-muted/50 p-3">
      <dt className="text-xs text-muted-foreground">{label}</dt>
      <dd className="mt-1 text-lg font-semibold tabular-nums">
        {loading ? (
          <span
            className="block h-6 w-10 animate-pulse rounded-md bg-muted"
            aria-hidden="true"
          />
        ) : (
          value
        )}
      </dd>
    </div>
  )
}

export function ProfileStats() {
  const { t } = useTranslation('common')
  const { data: progress, isPending: progressPending } = useProgress()
  const { data: itemCount, isPending: itemsPending } = useItemCount()

  const rewardCount = progress?.badges.filter(
    (badge) => badge.earnedAt !== null,
  ).length

  return (
    <dl
      className="grid grid-cols-2 gap-2 sm:grid-cols-3"
      data-testid="profile-stats"
    >
      <StatTile
        label={t('progress.level')}
        value={String(progress?.level ?? 0)}
        loading={progressPending}
      />
      <StatTile
        label={t('progress.xp')}
        value={String(progress?.xp ?? 0)}
        loading={progressPending}
      />
      <StatTile
        label={t('progress.streak')}
        value={String(progress?.currentStreak ?? 0)}
        loading={progressPending}
      />
      <StatTile
        label={t('nav.myItems')}
        value={String(itemCount ?? 0)}
        loading={itemsPending}
      />
      <StatTile
        label={t('progress.badges')}
        value={String(rewardCount ?? 0)}
        loading={progressPending}
      />
      <StatTile
        label={t('progress.rewards')}
        value={String(progress?.earnedBadgeCount ?? 0)}
        loading={progressPending}
      />
    </dl>
  )
}
