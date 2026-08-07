import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { Button } from '#/components/ui'
import { useSession } from '#/features/auth/session'
import { useLeaderboard } from '#/features/leaderboard/hooks'
import { leaderboardErrorKey, mergePages } from '#/features/leaderboard/lib'
import { LeaderboardTable } from './LeaderboardTable'
import { NeighborsPanel } from './NeighborsPanel'
import {
  LeaderboardEmpty,
  LeaderboardError,
  LeaderboardGuest,
  LeaderboardLoading,
} from './StateViews'

export function LeaderboardScreen() {
  const { t } = useTranslation('leaderboard')
  const { isAuthenticated, isLoading: isSessionLoading } = useSession()
  const {
    data,
    isPending,
    isError,
    error,
    refetch,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = useLeaderboard()

  const entries = useMemo(() => mergePages(data?.pages ?? []), [data])
  const myRank = data?.pages[0]?.myRank ?? null
  const isGuest = !isSessionLoading && !isAuthenticated

  return (
    <main className="mx-auto flex w-full max-w-3xl flex-col gap-5 px-4 py-6">
      <header className="flex flex-col gap-1">
        <h1 className="font-heading text-xl font-semibold text-foreground">
          {t('title')}
        </h1>
        <p className="text-sm text-muted-foreground">{t('subtitle')}</p>
      </header>

      {isAuthenticated && (
        <>
          <p
            className="rounded-lg bg-primary-subtle px-3 py-2 text-sm font-medium tabular-nums text-primary-subtle-foreground"
            data-testid="my-rank"
          >
            {myRank === null ? t('notRanked') : t('myRank', { rank: myRank })}
          </p>

          <NeighborsPanel />
        </>
      )}

      <section
        aria-labelledby="leaderboard-top"
        className="flex flex-col gap-3"
      >
        <h2
          id="leaderboard-top"
          className="text-sm font-semibold text-foreground"
        >
          {t('top.title')}
        </h2>

        {isGuest && <LeaderboardGuest />}

        {!isGuest && isPending && <LeaderboardLoading label={t('loading')} />}

        {!isGuest && !isPending && isError && (
          <LeaderboardError
            message={t(leaderboardErrorKey(error) as 'errors.loadFailed')}
            onRetry={() => void refetch()}
          />
        )}

        {!isGuest && !isPending && !isError && entries.length === 0 && (
          <LeaderboardEmpty title={t('empty')} hint={t('emptyHint')} />
        )}

        {!isGuest && !isPending && !isError && entries.length > 0 && (
          <>
            <LeaderboardTable
              entries={entries}
              myRank={myRank}
              caption={t('caption')}
            />
            {hasNextPage && (
              <Button
                variant="outline"
                size="sm"
                className="self-center"
                onClick={() => void fetchNextPage()}
                disabled={isFetchingNextPage}
              >
                {isFetchingNextPage
                  ? t('actions.loading')
                  : t('actions.loadMore')}
              </Button>
            )}
          </>
        )}
      </section>
    </main>
  )
}
