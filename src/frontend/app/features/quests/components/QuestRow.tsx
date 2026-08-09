import { useTranslation } from 'react-i18next'
import { cn } from '#/lib/utils'
import type { QuestProgressView } from '#/features/quests/types'

export interface QuestRowProps {
  entry: QuestProgressView
}

export function QuestRow({ entry }: QuestRowProps) {
  const { t } = useTranslation('quests')
  const { quest, current, target, percent, completed, claimed } = entry
  const title = t(`items.${quest.id}` as 'items.favorite_three', {
    count: target,
    defaultValue: t('items.fallback', { count: target }),
  })

  return (
    <li
      data-testid="quest-row"
      data-quest-id={quest.id}
      data-completed={completed}
      className="flex flex-col gap-1.5"
    >
      <div className="flex items-start justify-between gap-2">
        <p
          className={cn(
            'text-xs leading-snug font-medium wrap-anywhere',
            completed ? 'text-muted-foreground' : 'text-foreground',
          )}
        >
          {completed && <span aria-hidden="true">✓ </span>}
          {title}
        </p>
        <span
          className={cn(
            'shrink-0 rounded-full px-2 py-0.5 text-[11px] font-semibold',
            claimed
              ? 'bg-success-subtle text-success-subtle-foreground'
              : 'bg-primary-subtle text-primary-subtle-foreground',
          )}
        >
          {claimed
            ? t('rewardClaimed', { xp: quest.reward_xp })
            : t('reward', { xp: quest.reward_xp })}
        </span>
      </div>

      <div
        role="progressbar"
        aria-valuemin={0}
        aria-valuemax={100}
        aria-valuenow={percent}
        aria-label={t('meterLabel', { title, current, target })}
        className="h-1.5 w-full overflow-hidden rounded-full bg-[var(--xp-track)]"
      >
        <div
          className={cn(
            'h-full rounded-full transition-all',
            completed ? 'bg-success' : 'bg-[var(--xp-fill)]',
          )}
          style={{ width: `${percent}%` }}
        />
      </div>

      <p
        data-testid="quest-progress"
        className="text-[11px] leading-snug text-muted-foreground"
      >
        {t('progress', { current, target })}
      </p>
    </li>
  )
}
