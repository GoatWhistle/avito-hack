import { useTranslation } from 'react-i18next'
import { loadErrorKey } from '#/features/rewards/lib'
import { useRewardCatalog } from '#/features/rewards/hooks'
import type { RewardGroup } from '#/features/rewards/types'
import { RewardCard } from './RewardCard'
import { EmptyState, ErrorState, LoadingState } from './StateViews'

function GroupSection({ group }: { group: RewardGroup }) {
  const { t } = useTranslation('rewards')
  const isClose = group.key === 'close'

  return (
    <section
      aria-labelledby={`reward-group-${group.key}`}
      className="flex flex-col gap-3"
    >
      <div className="flex items-baseline justify-between gap-2">
        <h2
          id={`reward-group-${group.key}`}
          className="text-sm font-semibold text-foreground"
        >
          {t(`groups.${group.key}` as 'groups.close')}
        </h2>
        <span className="tabular-nums text-xs text-muted-foreground">
          {group.entries.length}
        </span>
      </div>
      <p className="-mt-2 text-xs text-muted-foreground">
        {t(`groupHints.${group.key}` as 'groupHints.close')}
      </p>
      <div className="grid gap-3 sm:grid-cols-2">
        {group.entries.map((entry) => (
          <RewardCard
            key={entry.reward.id}
            entry={entry}
            highlighted={isClose}
          />
        ))}
      </div>
    </section>
  )
}

export function RewardCatalog() {
  const { t } = useTranslation('rewards')
  const { data, isPending, isError, error, refetch } = useRewardCatalog()

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

  return (
    <div className="flex flex-col gap-6">
      {data.map((group) => (
        <GroupSection key={group.key} group={group} />
      ))}
    </div>
  )
}
