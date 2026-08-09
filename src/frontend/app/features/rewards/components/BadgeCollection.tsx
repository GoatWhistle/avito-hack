import { useTranslation } from 'react-i18next'
import { useBadges } from '#/features/rewards/hooks'
import { loadErrorKey } from '#/features/rewards/lib'
import { AchievementsGroup } from './AchievementsGroup'
import { AchievementsSummary } from './AchievementsSummary'
import { EmptyState, ErrorState, LoadingState } from './StateViews'

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

  const earned = data.filter((badge) => Boolean(badge.earned_at))
  const locked = data.filter((badge) => !badge.earned_at)

  return (
    <section aria-labelledby="badges-title" className="flex flex-col gap-4">
      <AchievementsSummary earned={earned.length} total={data.length} />

      <AchievementsGroup title={t('badges.groups.earned')} badges={earned} />

      <AchievementsGroup title={t('badges.groups.locked')} badges={locked} />
    </section>
  )
}
