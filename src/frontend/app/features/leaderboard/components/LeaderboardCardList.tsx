import { isMe } from '#/features/leaderboard/lib'
import type { LeaderboardEntry } from '#/features/leaderboard/types'
import { LeaderboardCard } from './LeaderboardCard'

export function LeaderboardCardList({
  entries,
  myRank,
  label,
}: {
  entries: LeaderboardEntry[]
  myRank: number | null
  label: string
}) {
  return (
    <ul aria-label={label} className="flex flex-col gap-2">
      {entries.map((entry) => (
        <LeaderboardCard
          key={entry.user_id}
          entry={entry}
          isCurrentUser={isMe(entry, myRank)}
        />
      ))}
    </ul>
  )
}
