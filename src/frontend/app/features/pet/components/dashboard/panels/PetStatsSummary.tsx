import { useTranslation } from 'react-i18next'
import { Card, CardContent } from '#/components/ui'
import { useBadges } from '#/features/rewards'
import { levelProgress } from '#/features/pet/lib'
import type { Pet } from '#/features/pet/types'

export interface PetStatsSummaryProps {
  pet: Pet
}

interface SummaryTileProps {
  label: string
  value: string
  testId: string
}

function SummaryTile({ label, value, testId }: SummaryTileProps) {
  return (
    <div
      data-testid={testId}
      className="min-w-0 rounded-xl bg-muted/50 px-3 py-2 ring-1 ring-foreground/5"
    >
      <p className="truncate text-[10px] font-medium tracking-wide text-muted-foreground uppercase">
        {label}
      </p>
      <p className="mt-1 truncate font-mono text-base font-bold tabular-nums text-foreground">
        {value}
      </p>
    </div>
  )
}

export function PetStatsSummary({ pet }: PetStatsSummaryProps) {
  const { t } = useTranslation('pet')
  const { data: badges } = useBadges()
  const { level, xp, isMaxLevel } = levelProgress(pet)
  const badgeCount = badges?.length ?? 0
  const earned = badges?.map((badge) => badge.name).filter(Boolean) ?? []

  return (
    <Card size="sm" data-testid="pet-stats-summary">
      <CardContent className="flex flex-col gap-3">
        <h3 className="text-[11px] font-bold tracking-wide text-muted-foreground uppercase">
          {t('summaryCard.title')}
        </h3>

        <div className="grid grid-cols-2 gap-2">
          <SummaryTile
            testId="summary-level"
            label={t('stats.level')}
            value={isMaxLevel ? t('level.maxBadge') : String(level)}
          />
          <SummaryTile
            testId="summary-xp"
            label={t('stats.xp')}
            value={String(xp)}
          />
          <SummaryTile
            testId="summary-streak"
            label={t('streak.label')}
            value={t('dashboard.hud.streak', { count: pet.streak_days })}
          />
          <SummaryTile
            testId="summary-badges"
            label={t('summaryCard.badges')}
            value={String(badgeCount)}
          />
        </div>

        <p className="text-xs text-muted-foreground">
          {earned.length > 0
            ? t('summaryCard.badgeList', { list: earned.join(', ') })
            : t('summaryCard.badgesEmpty')}
        </p>
      </CardContent>
    </Card>
  )
}
