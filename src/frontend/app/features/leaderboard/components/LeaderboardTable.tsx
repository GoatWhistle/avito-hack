import { useTranslation } from 'react-i18next'
import { isMe } from '#/features/leaderboard/lib'
import type { LeaderboardEntry } from '#/features/leaderboard/types'
import { LeaderboardRow } from './LeaderboardRow'

export function LeaderboardTable({
  entries,
  myRank,
  caption,
}: {
  entries: LeaderboardEntry[]
  myRank: number | null
  caption: string
}) {
  const { t } = useTranslation('leaderboard')

  return (
    <div className="overflow-x-auto rounded-xl ring-1 ring-foreground/10">
      <table className="w-full border-collapse text-left">
        <caption className="sr-only">{caption}</caption>
        <thead>
          <tr className="border-b border-border bg-muted/50">
            <th
              scope="col"
              className="w-px px-2 py-2 text-xs font-medium whitespace-nowrap text-muted-foreground sm:px-3"
            >
              {t('columns.rank')}
            </th>
            <th
              scope="col"
              className="w-full px-2 py-2 text-xs font-medium text-muted-foreground sm:px-3"
            >
              {t('columns.name')}
            </th>
            <th
              scope="col"
              className="w-px px-2 py-2 text-right text-xs font-medium whitespace-nowrap text-muted-foreground sm:px-3"
            >
              {t('columns.level')}
            </th>
            <th
              scope="col"
              className="w-px px-2 py-2 text-right text-xs font-medium whitespace-nowrap text-muted-foreground sm:px-3"
            >
              {t('columns.xp')}
            </th>
            <th
              scope="col"
              className="hidden w-px px-2 py-2 text-right text-xs font-medium whitespace-nowrap text-muted-foreground sm:table-cell sm:px-3"
            >
              {t('columns.streak')}
            </th>
          </tr>
        </thead>
        <tbody>
          {entries.map((entry) => (
            <LeaderboardRow
              key={entry.user_id}
              entry={entry}
              isCurrentUser={isMe(entry, myRank)}
            />
          ))}
        </tbody>
      </table>
    </div>
  )
}
