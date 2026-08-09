import { useTranslation } from 'react-i18next'
import { cn } from '#/lib/utils'

const WEEK_LENGTH = 7

export interface StreakTrackerProps {
  streakDays: number
  checkedInToday?: boolean
}

export interface StreakWeek {
  week: number
  filled: number
}

export const streakWeek = (streakDays: number): StreakWeek => {
  const safe = Math.max(0, Math.trunc(streakDays))
  const remainder = safe % WEEK_LENGTH
  const filled = remainder === 0 && safe > 0 ? WEEK_LENGTH : remainder

  return { week: Math.floor((safe - 1) / WEEK_LENGTH) + 1, filled }
}

export function StreakTracker({
  streakDays,
  checkedInToday = false,
}: StreakTrackerProps) {
  const { t } = useTranslation('pet')
  const { week, filled } = streakWeek(streakDays)
  const days = Array.from({ length: WEEK_LENGTH }, (_, index) => index + 1)

  return (
    <section
      aria-label={t('dashboard.streakTracker')}
      className="flex w-full max-w-sm items-center gap-3 rounded-3xl bg-card px-4 py-2.5 ring-1 ring-foreground/10 lg:max-w-none"
    >
      <span aria-hidden="true" className="text-lg leading-none">
        🔥
      </span>

      <div className="min-w-0">
        <p className="font-mono text-lg leading-none font-bold tabular-nums text-stat-streak-fill">
          {streakDays}
        </p>
        <p className="mt-0.5 text-[11px] leading-tight whitespace-nowrap text-muted-foreground">
          {streakDays > 0
            ? t('dashboard.streakWeek', { count: week })
            : t('streak.startHint')}
        </p>
      </div>

      <ol className="ml-auto flex items-center gap-1.5">
        {days.map((day) => {
          const isFilled = day <= filled
          const isTarget = !checkedInToday && day === filled + 1

          return (
            <li
              key={day}
              aria-label={
                isFilled
                  ? t('dashboard.streakDayDone', { day })
                  : t('dashboard.streakDay', { day })
              }
              aria-current={isTarget ? 'step' : undefined}
              className={cn(
                'size-3.5 rounded-full',
                isFilled ? 'bg-stat-streak-fill' : 'bg-muted',
                isTarget &&
                  'ring-2 ring-stat-streak-fill/50 ring-offset-1 ring-offset-card',
              )}
            />
          )
        })}
      </ol>
    </section>
  )
}
