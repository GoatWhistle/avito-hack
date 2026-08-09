import { useTranslation } from 'react-i18next'
import { cn } from '#/lib/utils'

export interface StreakCardProps {
  days: number
  freezes: number
  atRisk: boolean
}

export function StreakCard({ days, freezes, atRisk }: StreakCardProps) {
  const { t } = useTranslation('pet')
  const alive = days > 0

  return (
    <section
      aria-labelledby="pet-streak-title"
      className={cn(
        'flex flex-wrap items-center gap-x-3 gap-y-2 rounded-xl px-4 py-3 ring-1 transition-colors',
        alive
          ? 'bg-warning-subtle ring-streak-flame/30'
          : 'bg-muted ring-border',
      )}
    >
      <span
        aria-hidden="true"
        className={cn(
          'text-2xl leading-none',
          alive ? 'animate-pulse' : 'opacity-40 grayscale',
        )}
      >
        🔥
      </span>

      <div className="flex min-w-32 flex-1 flex-col gap-0.5">
        <h2
          id="pet-streak-title"
          className="text-xs font-medium tracking-wide text-muted-foreground uppercase"
        >
          {t('streak.label')}
        </h2>
        <p className="text-sm font-semibold text-balance text-foreground">
          {alive ? t('streak.days', { count: days }) : t('streak.empty')}
        </p>
        <p className="text-xs text-pretty text-muted-foreground">
          {atRisk
            ? t('streak.atRisk')
            : alive
              ? t('streak.tomorrow', { count: days + 1 })
              : t('streak.startHint')}
        </p>
      </div>

      {freezes > 0 && (
        <span className="shrink-0 rounded-full bg-info-subtle px-2 py-1 text-xs font-medium text-info-subtle-foreground">
          {t('streak.freezes', { count: freezes })}
        </span>
      )}
    </section>
  )
}
