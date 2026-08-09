import { useTranslation } from 'react-i18next'
import { Card, CardContent, SidebarHeader } from '#/components/ui'
import type { LevelProgress } from '#/features/pet/lib'
import type { Pet } from '#/features/pet/types'

export interface DashboardHeaderProps {
  pet: Pet
  progress: LevelProgress
}

export function DashboardHeader({ pet, progress }: DashboardHeaderProps) {
  const { t } = useTranslation('pet')
  const { level, xp, nextLevelXp, xpToNext, percent, isMaxLevel } = progress

  return (
    <SidebarHeader className="p-3">
      <Card
        size="sm"
        data-testid="pet-card"
        aria-labelledby="pet-card-name"
        className="bg-sidebar-accent ring-sidebar-border"
      >
        <CardContent className="flex flex-col gap-3">
          <div className="flex items-center gap-3">
            <span
              aria-hidden="true"
              className="flex size-10 shrink-0 items-center justify-center rounded-full bg-primary text-lg text-primary-foreground"
            >
              🦝
            </span>

            <div className="min-w-0 flex-1">
              <h2
                id="pet-card-name"
                className="truncate text-sm font-semibold text-card-foreground"
              >
                {pet.name}
              </h2>
              <p className="truncate text-xs text-muted-foreground">
                {t('dashboard.stageLevel', {
                  stage: t(`stage.${pet.stage}`),
                  level,
                })}
              </p>
            </div>
          </div>

          <div className="flex flex-col gap-1.5">
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
              className="h-2 w-full overflow-hidden rounded-full bg-xp-track"
            >
              <div
                className="h-full rounded-full bg-xp-fill transition-[width] duration-700 ease-out"
                style={{ width: `${percent}%` }}
              />
            </div>

            <p className="truncate font-mono text-[11px] tabular-nums text-muted-foreground">
              {isMaxLevel
                ? t('nextLevel.xpTotal', { xp })
                : t('nextLevel.xpCounter', { xp, total: nextLevelXp })}
            </p>
          </div>
        </CardContent>
      </Card>
    </SidebarHeader>
  )
}
