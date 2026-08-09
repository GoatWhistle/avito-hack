import { useTranslation } from 'react-i18next'
import { Card, CardContent } from '#/components/ui'
import { cn } from '#/lib/utils'
import {
  MAX_FREEZES,
  STREAK_BONUS_DAYS,
  STREAK_BONUS_FACTOR,
  STREAK_MILESTONES,
  STREAK_RESET_HOURS,
  isBonusActive,
  streakRules,
  type StreakRuleTone,
  type StreakRuleView,
} from '#/features/pet/lib'
import type { Pet } from '#/features/pet/types'

export interface StreakMechanicsCardProps {
  pet: Pet
}

const toneClass: Record<StreakRuleTone, string> = {
  active: 'bg-stat-streak-subtle text-stat-streak-subtle-foreground',
  pending: 'bg-primary-subtle text-primary-subtle-foreground',
  muted: 'bg-muted text-muted-foreground',
}

function RuleBadge({ rule, label }: { rule: StreakRuleView; label: string }) {
  return (
    <span
      data-testid={`streak-rule-badge-${rule.key}`}
      className={cn(
        'shrink-0 rounded-full px-2 py-0.5 text-[10px] font-semibold whitespace-nowrap tabular-nums',
        toneClass[rule.tone],
      )}
    >
      {label}
    </span>
  )
}

export function StreakMechanicsCard({ pet }: StreakMechanicsCardProps) {
  const { t } = useTranslation('pet')
  const rules = streakRules(pet)
  const streakDays = Math.max(0, Math.trunc(pet.streak_days))

  const badgeLabel = (rule: StreakRuleView): string => {
    if (rule.key === 'bonus') {
      return isBonusActive(streakDays)
        ? t('mechanics.badge.active')
        : t('mechanics.badge.daysLeft', { count: rule.badgeCount })
    }
    if (rule.key === 'reset') {
      return t('mechanics.badge.hours', { count: STREAK_RESET_HOURS })
    }

    return rule.badgeCount > 0
      ? t('mechanics.badge.freezesLeft', { count: rule.badgeCount })
      : t('mechanics.badge.unavailable')
  }

  const ruleText = (rule: StreakRuleView): string => {
    if (rule.key === 'bonus') {
      return t('mechanics.bonus.text', {
        days: STREAK_BONUS_DAYS,
        factor: STREAK_BONUS_FACTOR,
      })
    }
    if (rule.key === 'reset') return t('mechanics.reset.text')

    return t('mechanics.freeze.text', {
      max: MAX_FREEZES,
      milestones: STREAK_MILESTONES.join(', '),
    })
  }

  const ruleTitle = (rule: StreakRuleView): string =>
    rule.key === 'bonus'
      ? t('streak.days', { count: streakDays })
      : t(`mechanics.${rule.key}.title`)

  return (
    <Card size="sm" data-testid="streak-mechanics">
      <CardContent className="flex flex-col gap-3">
        <h3 className="text-[11px] font-bold tracking-wide text-muted-foreground uppercase">
          {t('mechanics.title')}
        </h3>

        <ul className="flex flex-col gap-2">
          {rules.map((rule) => (
            <li
              key={rule.key}
              data-testid={`streak-rule-${rule.key}`}
              className="flex items-start gap-2.5 rounded-xl bg-muted/40 p-2.5 ring-1 ring-foreground/5"
            >
              <span
                aria-hidden="true"
                className="flex size-7 shrink-0 items-center justify-center rounded-full bg-card text-sm ring-1 ring-foreground/10"
              >
                {rule.icon}
              </span>

              <div className="flex min-w-0 flex-1 flex-col gap-1">
                <div className="flex flex-wrap items-center justify-between gap-x-2 gap-y-1">
                  <p className="text-xs font-semibold text-foreground">
                    {ruleTitle(rule)}
                  </p>
                  <RuleBadge rule={rule} label={badgeLabel(rule)} />
                </div>
                <p className="text-[11px] leading-snug text-muted-foreground">
                  {ruleText(rule)}
                </p>
              </div>
            </li>
          ))}
        </ul>
      </CardContent>
    </Card>
  )
}
