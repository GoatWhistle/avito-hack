import { Award, Medal, Trophy } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { cn } from '#/lib/utils'
import { initialsOf, medalOf } from '#/features/leaderboard/lib'
import type { LeaderboardEntry } from '#/features/leaderboard/types'

const MEDAL_ICON = {
  gold: Trophy,
  silver: Medal,
  bronze: Award,
} as const

const MEDAL_TONE = {
  gold: 'text-warning',
  silver: 'text-muted-foreground',
  bronze: 'text-stat-streak-fill',
} as const

function MedalIcon({ medal }: { medal: 'gold' | 'silver' | 'bronze' }) {
  const Icon = MEDAL_ICON[medal]

  return (
    <Icon
      aria-hidden="true"
      className={cn('size-[1.125rem]', MEDAL_TONE[medal])}
    />
  )
}

export function LeaderboardCard({
  entry,
  isCurrentUser = false,
}: {
  entry: LeaderboardEntry
  isCurrentUser?: boolean
}) {
  const { t } = useTranslation('leaderboard')
  const medal = medalOf(entry.rank)

  return (
    <li
      data-testid="leaderboard-card"
      data-rank={entry.rank}
      data-me={isCurrentUser || undefined}
      aria-current={isCurrentUser ? 'true' : undefined}
      className={cn(
        'flex items-center gap-3 rounded-xl bg-card p-3 ring-1 transition-colors',
        isCurrentUser
          ? 'bg-primary-subtle ring-primary/30'
          : medal !== null
            ? 'ring-warning/30'
            : 'ring-foreground/10',
      )}
    >
      <span className="flex size-7 shrink-0 items-center justify-center">
        {medal !== null ? (
          <MedalIcon medal={medal} />
        ) : (
          <span className="font-mono text-sm font-bold tabular-nums text-muted-foreground">
            {entry.rank}
          </span>
        )}
      </span>

      <span
        aria-hidden="true"
        className="flex size-9 shrink-0 items-center justify-center rounded-full bg-primary text-xs font-bold text-primary-foreground"
      >
        {initialsOf(entry.name)}
      </span>

      <div className="flex min-w-0 flex-1 flex-col gap-0.5">
        <div className="flex flex-wrap items-center gap-1.5">
          <span
            className={cn(
              'min-w-0 text-sm font-semibold break-words',
              isCurrentUser
                ? 'text-primary-subtle-foreground'
                : 'text-foreground',
            )}
          >
            {entry.name}
          </span>
          {isCurrentUser && (
            <span className="shrink-0 rounded-full bg-primary px-1.5 py-0.5 text-[10px] font-semibold text-primary-foreground">
              {t('you')}
            </span>
          )}
        </div>

        <p className="text-[11px] tabular-nums text-muted-foreground">
          {t('cardStats', {
            level: entry.level,
            xp: entry.xp,
            streak: entry.streak_days,
          })}
        </p>
      </div>
    </li>
  )
}
