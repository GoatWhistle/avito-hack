import { useTranslation } from 'react-i18next'
import { useMyRewards } from '#/features/rewards/hooks'
import { loadErrorKey } from '#/features/rewards/lib'
import { PromoCodeCard } from './PromoCodeCard'
import { EmptyState, ErrorState, LoadingState } from './StateViews'

export function MyRewards() {
  const { t } = useTranslation('rewards')
  const { data, isPending, isError, error, refetch } = useMyRewards()

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
    return <EmptyState title={t('empty.mine')} hint={t('empty.mineHint')} />
  }

  return (
    <div className="grid gap-3 sm:grid-cols-2">
      {data.map((item) => (
        <PromoCodeCard key={item.reward_id} item={item} />
      ))}
    </div>
  )
}
