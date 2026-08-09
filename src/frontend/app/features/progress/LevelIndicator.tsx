import { useTranslation } from 'react-i18next'
import { useLocation } from 'react-router'
import { Flame } from 'lucide-react'
import { cn } from '#/lib/utils'
import { useProgress } from './useProgress'
import { xpProgressRatio } from './progress.types'

interface LevelIndicatorProps {
  className?: string
  compact?: boolean
}

export function LevelIndicator({ className, compact }: LevelIndicatorProps) {
  const { t } = useTranslation('common')
  const { pathname } = useLocation()
  const { data, isPending, isError } = useProgress()
  const isPetDashboard = pathname === '/pet' || pathname.startsWith('/pet/')

  if (isPetDashboard) return null

  if (isPending) {
    return (
      <div
        className={cn('flex items-center gap-2', className)}
        data-testid="level-indicator-skeleton"
        aria-hidden="true"
      >
        <div className="h-5 w-14 animate-pulse rounded-md bg-muted" />
        <div className="h-2 w-20 animate-pulse rounded-full bg-muted" />
      </div>
    )
  }

  if (isError || !data) return null

  const ratio = xpProgressRatio(data)
  const percent = Math.round(ratio * 100)
  const atMax = data.xpToNextLevel <= 0

  const label = t('progress.indicatorLabel', {
    level: data.level,
    xp: data.xp,
    streak: data.currentStreak,
  })

  return (
    <div
      className={cn('flex items-center gap-2', className)}
      data-testid="level-indicator"
      aria-label={label}
    >
      <span className="shrink-0 rounded-md bg-primary-subtle px-1.5 py-0.5 text-xs font-semibold text-primary-subtle-foreground tabular-nums">
        {t('progress.levelShort', { level: data.level })}
      </span>

      <div className="flex min-w-0 flex-col gap-1">
        <div
          className="h-1.5 w-16 overflow-hidden rounded-full bg-xp-track sm:w-24"
          role="progressbar"
          aria-valuenow={percent}
          aria-valuemin={0}
          aria-valuemax={100}
          aria-label={t('progress.xpProgress')}
        >
          <div
            className="h-full rounded-full bg-xp-fill transition-[width] duration-normal ease-out motion-reduce:transition-none"
            style={{ width: `${percent}%` }}
          />
        </div>
        {!compact && (
          <span className="text-[0.6875rem] leading-none text-muted-foreground tabular-nums">
            {atMax
              ? t('progress.maxLevel')
              : t('progress.xpOfNext', {
                  xp: data.xp,
                  next: data.xpToNextLevel,
                })}
          </span>
        )}
      </div>

      {data.currentStreak > 0 && (
        <span
          className="flex shrink-0 items-center gap-0.5 text-xs font-medium text-streak-flame tabular-nums"
          aria-label={t('progress.streakLabel', { count: data.currentStreak })}
        >
          <Flame className="size-3.5" aria-hidden="true" />
          {data.currentStreak}
        </span>
      )}
    </div>
  )
}
