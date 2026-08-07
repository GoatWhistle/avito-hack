import { useTranslation } from 'react-i18next'
import { cn } from '#/lib/utils'
import type { StatKey, StatTone } from '#/features/pet/lib'

const TONE_FILL: Record<StatTone, string> = {
  critical: 'bg-destructive',
  warning: 'bg-warning',
  good: 'bg-success',
}

const TONE_TEXT: Record<StatTone, string> = {
  critical: 'text-destructive',
  warning: 'text-warning',
  good: 'text-success',
}

const ICON: Record<StatKey, string> = {
  satiety: '🍒',
  happiness: '💚',
  energy: '⚡',
}

export interface StatMeterProps {
  statKey: StatKey
  value: number
  tone: StatTone
}

export function StatMeter({ statKey, value, tone }: StatMeterProps) {
  const { t } = useTranslation('pet')
  const label = t(`stats.${statKey}`)
  const toneLabel = t(`tone.${tone}`)

  return (
    <div className="flex flex-col gap-1.5">
      <div className="flex items-baseline justify-between gap-2">
        <span className="flex items-center gap-1.5 text-sm text-muted-foreground">
          <span aria-hidden="true">{ICON[statKey]}</span>
          {label}
        </span>
        <span className={cn('font-mono text-sm tabular-nums', TONE_TEXT[tone])}>
          {value}
        </span>
      </div>

      <div
        role="meter"
        aria-valuenow={value}
        aria-valuemin={0}
        aria-valuemax={100}
        aria-label={t('stats.meterLabel', {
          stat: label,
          value,
          tone: toneLabel,
        })}
        className="h-2 w-full overflow-hidden rounded-full bg-muted"
      >
        <div
          className={cn(
            'h-full rounded-full transition-[width] duration-500 ease-out',
            TONE_FILL[tone],
          )}
          style={{ width: `${value}%` }}
        />
      </div>

      {tone !== 'good' && (
        <p className={cn('text-xs', TONE_TEXT[tone])}>
          {t(`hints.${statKey}.${tone}`)}
        </p>
      )}
    </div>
  )
}
