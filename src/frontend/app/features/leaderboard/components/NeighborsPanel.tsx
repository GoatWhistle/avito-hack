import { useTranslation } from 'react-i18next'
import { chaseTarget, neighborsOf } from '#/features/leaderboard/lib'
import { useMyNeighborhood } from '#/features/leaderboard/hooks'
import { LeaderboardTable } from './LeaderboardTable'
import { RowsSkeleton } from './StateViews'

export function NeighborsPanel() {
  const { t } = useTranslation('leaderboard')
  const { data, isPending, isError } = useMyNeighborhood()

  if (isPending) return <RowsSkeleton count={3} />
  if (isError || !data) return null

  const { myRank, items } = data

  if (myRank === null) {
    return (
      <section className="flex flex-col gap-2 rounded-xl bg-muted/50 p-4">
        <h2 className="text-sm font-semibold text-foreground">
          {t('neighbors.title')}
        </h2>
        <p className="text-sm text-muted-foreground">{t('notRanked')}</p>
      </section>
    )
  }

  const neighbors = neighborsOf(items, myRank)
  if (neighbors.length === 0) return null

  const chase = chaseTarget(items, myRank)

  return (
    <section aria-labelledby="neighbors-title" className="flex flex-col gap-2">
      <h2
        id="neighbors-title"
        className="text-sm font-semibold text-foreground"
      >
        {t('neighbors.title')}
      </h2>

      {chase ? (
        <p className="text-sm text-muted-foreground">
          {t('neighbors.chase', { name: chase.target.name, xp: chase.gap })}
        </p>
      ) : (
        <p className="text-sm text-muted-foreground">{t('neighbors.leader')}</p>
      )}

      <LeaderboardTable
        entries={neighbors}
        myRank={myRank}
        caption={t('neighbors.caption')}
      />
    </section>
  )
}
