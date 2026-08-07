import { useTranslation } from 'react-i18next'
import { cn } from '#/lib/utils'
import { medalOf } from '#/features/leaderboard/lib'
import type { LeaderboardEntry } from '#/features/leaderboard/types'

const medalClass = {
  gold: 'bg-warning-subtle text-warning-subtle-foreground',
  silver: 'bg-muted text-foreground',
  bronze: 'bg-destructive-subtle text-destructive-subtle-foreground',
} as const

export function LeaderboardRow({
  entry,
  isCurrentUser = false,
}: {
  entry: LeaderboardEntry
  isCurrentUser?: boolean
}) {
  const { t } = useTranslation('leaderboard')
  const medal = medalOf(entry.rank)

  return (
    <tr
      data-testid="leaderboard-row"
      data-rank={entry.rank}
      data-me={isCurrentUser || undefined}
      aria-current={isCurrentUser ? 'true' : undefined}
      className={cn(
        'border-b border-border last:border-0',
        isCurrentUser && 'bg-primary-subtle',
      )}
    >
      <td className="w-px px-2 py-2.5 sm:px-3">
        <span
          className={cn(
            'inline-flex h-6 min-w-6 items-center justify-center rounded-full px-1.5 text-xs font-semibold tabular-nums',
            medal ? medalClass[medal] : 'text-muted-foreground',
          )}
        >
          {entry.rank}
        </span>
      </td>
      <td className="w-full px-2 py-2.5 sm:px-3">
        <div className="flex items-center gap-1.5">
          <span
            className={cn(
              'min-w-0 text-sm break-words',
              isCurrentUser
                ? 'font-semibold text-primary-subtle-foreground'
                : 'text-foreground',
            )}
          >
            {entry.name}
          </span>
          {isCurrentUser && (
            <span className="shrink-0 rounded-full bg-primary px-1.5 py-0.5 text-[0.625rem] font-semibold text-primary-foreground">
              {t('you')}
            </span>
          )}
        </div>
      </td>
      <td className="w-px px-2 py-2.5 text-right text-sm whitespace-nowrap tabular-nums text-muted-foreground sm:px-3">
        {entry.level}
      </td>
      <td className="w-px px-2 py-2.5 text-right text-sm font-medium whitespace-nowrap tabular-nums text-foreground sm:px-3">
        {entry.xp}
      </td>
      <td className="hidden w-px px-2 py-2.5 text-right text-sm whitespace-nowrap tabular-nums text-muted-foreground sm:table-cell sm:px-3">
        {entry.streak_days}
      </td>
    </tr>
  )
}
