import { useTranslation } from 'react-i18next'
import { cn } from '#/lib/utils'
import { useBadges } from '#/features/rewards/hooks'
import { formatDate, loadErrorKey } from '#/features/rewards/lib'
import type { BadgeItem } from '#/features/rewards/types'
import { EmptyState, ErrorState, LoadingState } from './StateViews'

function BadgeTile({ badge }: { badge: BadgeItem }) {
  const { t, i18n } = useTranslation('rewards')
  const earned = Boolean(badge.earned_at)

  return (
    <li
      data-testid="badge-tile"
      data-earned={earned}
      className={cn(
        'flex flex-col gap-1 rounded-xl p-3 ring-1 ring-foreground/10',
        earned ? 'bg-card' : 'bg-muted/40 opacity-60',
      )}
    >
      <span className="line-clamp-2 text-sm font-medium text-balance text-foreground">
        {badge.name}
      </span>
      <span className="line-clamp-2 text-xs text-muted-foreground">
        {badge.description}
      </span>
      <span className="mt-1 text-xs text-muted-foreground">
        {earned
          ? t('badges.earnedAt', {
              date: formatDate(badge.earned_at, i18n.language),
            })
          : t('badges.locked')}
      </span>
    </li>
  )
}

export function BadgeCollection() {
  const { t } = useTranslation('rewards')
  const { data, isPending, isError, error, refetch } = useBadges()

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
    return <EmptyState title={t('badges.empty')} hint={t('badges.emptyHint')} />
  }

  const earnedCount = data.filter((badge) => Boolean(badge.earned_at)).length

  return (
    <section aria-labelledby="badges-title" className="flex flex-col gap-3">
      <div className="flex items-baseline justify-between gap-2">
        <h2 id="badges-title" className="text-sm font-semibold text-foreground">
          {t('badges.title')}
        </h2>
        <span className="tabular-nums text-xs text-muted-foreground">
          {t('badges.counter', { earned: earnedCount, total: data.length })}
        </span>
      </div>
      <ul className="grid grid-cols-2 gap-3 sm:grid-cols-3">
        {data.map((badge) => (
          <BadgeTile key={badge.id} badge={badge} />
        ))}
      </ul>
    </section>
  )
}
