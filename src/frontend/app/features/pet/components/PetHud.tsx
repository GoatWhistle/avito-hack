import { useTranslation } from 'react-i18next'
import { cn } from '#/lib/utils'
import type { LevelProgress as LevelProgressData } from '#/features/pet/lib'
import type { Pet } from '#/features/pet/types'

export interface PetHudProps {
  pet: Pet
  progress: LevelProgressData
  streakAtRisk: boolean
}

export function PetHud({ pet, progress, streakAtRisk }: PetHudProps) {
  const { t } = useTranslation('pet')
  const { level, xp, nextLevelXp, xpToNext, percent, isMaxLevel } = progress

  return (
    <section
      aria-labelledby="pet-hud-title"
      className="flex flex-wrap items-center gap-x-4 gap-y-3 rounded-xl bg-card px-4 py-3 ring-1 ring-foreground/10"
    >
      <span
        aria-hidden="true"
        className="grid size-12 shrink-0 place-items-center rounded-xl bg-primary font-mono text-lg font-bold tabular-nums text-primary-foreground"
      >
        {level}
      </span>

      <div className="flex min-w-48 flex-1 flex-col gap-1.5">
        <div className="flex flex-wrap items-baseline justify-between gap-x-3">
          <h2
            id="pet-hud-title"
            className="text-sm font-semibold text-foreground"
          >
            {t('hud.title', { name: pet.name, level })}
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
          className="relative h-2.5 w-full overflow-hidden rounded-full bg-xp-track"
        >
          <div
            className="h-full rounded-full bg-xp-fill transition-[width] duration-700 ease-out"
            style={{ width: `${percent}%` }}
          />
          {isMaxLevel && (
            <div className="pointer-events-none absolute inset-0 animate-pulse rounded-full bg-linear-to-r from-transparent via-primary-foreground/30 to-transparent" />
          )}
        </div>

        <p className="text-xs text-muted-foreground">
          {isMaxLevel
            ? t('level.maxHint')
            : t('level.toNext', { level: level + 1, count: xpToNext })}
        </p>
      </div>

      <ul className="flex shrink-0 flex-wrap items-center gap-2">
        <li>
          <HudChip
            icon="🔥"
            tone={
              streakAtRisk
                ? 'warning'
                : pet.streak_days > 0
                  ? 'streak'
                  : 'muted'
            }
            label={t('hud.streak', { count: pet.streak_days })}
          />
        </li>
        <li>
          <HudChip icon="🌱" tone="muted" label={t(`stage.${pet.stage}`)} />
        </li>
        <li>
          <HudChip icon="💚" tone="info" label={t(`state.${pet.state}`)} />
        </li>
      </ul>
    </section>
  )
}

const CHIP_TONE = {
  streak: 'bg-warning-subtle text-warning-subtle-foreground',
  warning: 'bg-destructive-subtle text-destructive-subtle-foreground',
  info: 'bg-info-subtle text-info-subtle-foreground',
  muted: 'bg-muted text-muted-foreground',
} as const

interface HudChipProps {
  icon: string
  label: string
  tone: keyof typeof CHIP_TONE
}

function HudChip({ icon, label, tone }: HudChipProps) {
  return (
    <span
      className={cn(
        'flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-medium whitespace-nowrap',
        CHIP_TONE[tone],
      )}
    >
      <span aria-hidden="true">{icon}</span>
      {label}
    </span>
  )
}
