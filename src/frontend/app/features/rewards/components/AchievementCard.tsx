import { useTranslation } from 'react-i18next'
import { Card, CardContent } from '#/components/ui'
import { cn } from '#/lib/utils'
import {
  badgeProgress,
  badgeRemainingLabel,
  badgeText,
  formatDate,
} from '#/features/rewards/lib'
import type { BadgeItem } from '#/features/rewards/types'
import { ProgressBar } from './ProgressBar'

export interface AchievementCardProps {
  badge: BadgeItem
}

export function AchievementCard({ badge }: AchievementCardProps) {
  const { t, i18n } = useTranslation('rewards')
  const { t: tCatalog } = useTranslation('catalog')
  const earned = Boolean(badge.earned_at)
  const progress = badgeProgress(badge)
  const remaining = badgeRemainingLabel(t, badge)
  const showProgress = !earned && progress.hasProgress
  const text = badgeText(tCatalog, badge.id, {
    title: badge.name,
    description: badge.description,
  })

  return (
    <li className="flex">
      <Card
        size="sm"
        data-testid="badge-tile"
        data-earned={earned}
        className={cn(
          'flex-1 transition-[box-shadow,opacity] duration-200',
          earned
            ? 'bg-achievement-earned ring-achievement-earned-border/40 hover:ring-achievement-earned-border/60'
            : 'opacity-75 hover:opacity-100',
        )}
      >
        <CardContent className="flex h-full flex-col gap-3">
          <span
            aria-hidden="true"
            className={cn(
              'flex size-11 shrink-0 items-center justify-center rounded-xl text-lg ring-1',
              earned
                ? 'bg-achievement-earned-tile ring-achievement-earned-border/40'
                : 'bg-achievement-locked-tile opacity-70 grayscale ring-border',
            )}
          >
            {badge.icon_url ? (
              <img
                src={badge.icon_url}
                alt=""
                className="size-7 object-contain"
              />
            ) : (
              <span>{earned ? '🎖️' : '🔒'}</span>
            )}
          </span>

          <div className="flex min-w-0 flex-col gap-1">
            <p className="line-clamp-2 text-sm leading-tight font-semibold wrap-anywhere text-balance text-foreground">
              {text.title}
            </p>
            <p className="line-clamp-2 text-xs leading-snug font-medium wrap-anywhere text-muted-foreground">
              {text.description}
            </p>
          </div>

          <div className="mt-auto flex flex-col gap-1.5">
            {showProgress && (
              <>
                <ProgressBar
                  value={progress.percent}
                  label={t('badges.progressMeter', {
                    name: text.title,
                    current: progress.current,
                    target: progress.target,
                  })}
                  tone={progress.remaining <= 1 ? 'close' : 'default'}
                />
                <p
                  data-testid="badge-progress"
                  className="text-[11px] leading-snug font-medium wrap-anywhere text-muted-foreground"
                >
                  {t('badges.progress', {
                    current: progress.current,
                    target: progress.target,
                  })}
                </p>
              </>
            )}

            <p className="text-[11px] leading-snug wrap-anywhere text-muted-foreground/80">
              {earned
                ? t('badges.earnedAt', {
                    date: formatDate(badge.earned_at, i18n.language),
                  })
                : (remaining ?? t('badges.locked'))}
            </p>
          </div>
        </CardContent>
      </Card>
    </li>
  )
}
