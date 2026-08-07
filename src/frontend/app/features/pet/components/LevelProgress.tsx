import { useTranslation } from 'react-i18next'
import type { LevelProgress as LevelProgressData } from '#/features/pet/lib'

export interface LevelProgressProps {
  progress: LevelProgressData
}

export function LevelProgress({ progress }: LevelProgressProps) {
  const { t } = useTranslation('pet')
  const { level, xp, nextLevelXp, xpToNext, percent, isMaxLevel } = progress

  return (
    <section className="flex flex-col gap-2" aria-labelledby="pet-level-title">
      <div className="flex items-baseline justify-between gap-3">
        <h2
          id="pet-level-title"
          className="text-sm font-medium text-muted-foreground"
        >
          {t('stats.levelValue', { level })}
        </h2>
        <span className="font-mono text-xs tabular-nums text-muted-foreground">
          {isMaxLevel
            ? t('level.maxBadge')
            : t('stats.xpProgress', { current: xp, next: nextLevelXp })}
        </span>
      </div>

      <div
        role="progressbar"
        aria-valuenow={percent}
        aria-valuemin={0}
        aria-valuemax={100}
        aria-label={
          isMaxLevel
            ? t('level.maxAria')
            : t('level.aria', { level: level + 1, count: xpToNext })
        }
        className="relative h-3 w-full overflow-hidden rounded-full bg-xp-track"
      >
        <div
          className="h-full rounded-full bg-xp-fill transition-[width] duration-700 ease-out"
          style={{ width: `${percent}%` }}
        />
        {isMaxLevel && (
          <div className="pointer-events-none absolute inset-0 animate-pulse rounded-full bg-linear-to-r from-transparent via-primary-foreground/30 to-transparent" />
        )}
      </div>

      <p className="text-sm text-foreground">
        {isMaxLevel
          ? t('level.maxHint')
          : t('level.toNext', { level: level + 1, count: xpToNext })}
      </p>
    </section>
  )
}
