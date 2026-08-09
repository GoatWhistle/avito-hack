import { useTranslation } from 'react-i18next'
import { Card, CardContent, CardHeader } from '#/components/ui'
import type { LevelProgress } from '#/features/pet/lib'

export interface NextLevelCardProps {
  progress: LevelProgress
}

export function NextLevelCard({ progress }: NextLevelCardProps) {
  const { level, xpToNext, isMaxLevel } = progress
  const { t } = useTranslation('pet')

  return (
    <Card
      size="sm"
      data-testid="next-level-card"
      className="relative bg-level-up-surface ring-level-up-border/25"
    >
      <span
        aria-hidden="true"
        className="pointer-events-none absolute -top-4 -right-4 size-20 rounded-full bg-level-up-glow/20 blur-xl"
      />

      <CardHeader className="flex flex-row items-center justify-between gap-2 space-y-0">
        <span className="rounded-lg bg-level-up-tile/70 px-2.5 py-0.5 text-xs font-semibold text-level-up-border ring-1 ring-level-up-border/30">
          {isMaxLevel ? t('nextLevel.maxBadge') : t('nextLevel.title')}
        </span>
        <span aria-hidden="true" className="text-xl">
          {isMaxLevel ? '👑' : '🚀'}
        </span>
      </CardHeader>

      <CardContent>
        <p className="text-sm font-bold text-foreground">
          {isMaxLevel
            ? t('nextLevel.maxTitle')
            : t('nextLevel.level', { level: level + 1 })}
        </p>
        <p className="mt-0.5 text-xs font-medium text-muted-foreground">
          {isMaxLevel
            ? t('level.maxHint')
            : t('nextLevel.remaining', { count: xpToNext })}
        </p>
      </CardContent>
    </Card>
  )
}
