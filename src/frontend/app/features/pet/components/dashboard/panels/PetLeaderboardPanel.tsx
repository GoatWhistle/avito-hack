import { Trophy } from 'lucide-react'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { useLeaderboard } from '#/features/leaderboard'
import {
  LeaderboardCardList,
  LeaderboardEmpty,
  LeaderboardError,
  LeaderboardLoading,
  MyRankCard,
} from '#/features/leaderboard/components'
import { mergePages } from '#/features/leaderboard/lib'

export function PetLeaderboardPanel() {
  const { t } = useTranslation(['pet', 'leaderboard'])
  const { data, isPending, isError, refetch } = useLeaderboard()

  const entries = useMemo(() => mergePages(data?.pages ?? []), [data])
  const myRank = data?.pages[0]?.myRank ?? null

  return (
    <section className="flex flex-col gap-3 p-3">
      <h2 className="flex items-center gap-2 text-sm font-semibold text-foreground">
        <Trophy aria-hidden="true" className="size-4" />
        {t('pet:dashboard.nav.leaderboard')}
      </h2>

      {!isPending && !isError && entries.length > 0 && (
        <MyRankCard rank={myRank} />
      )}

      {isPending && <LeaderboardLoading label={t('leaderboard:loading')} />}

      {!isPending && isError && (
        <LeaderboardError
          message={t('leaderboard:errors.loadFailed')}
          onRetry={() => void refetch()}
        />
      )}

      {!isPending && !isError && entries.length === 0 && (
        <LeaderboardEmpty
          title={t('leaderboard:empty')}
          hint={t('leaderboard:emptyHint')}
        />
      )}

      {!isPending && !isError && entries.length > 0 && (
        <LeaderboardCardList
          entries={entries}
          myRank={myRank}
          label={t('leaderboard:listLabel')}
        />
      )}
    </section>
  )
}
