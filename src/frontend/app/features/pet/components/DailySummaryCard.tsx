import { X } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { Button } from '#/components/ui'
import { adviceText, summaryMessage } from '#/features/pet/lib'
import type { DailySummary } from '#/features/pet/types'

export interface DailySummaryCardProps {
  summary: DailySummary
  petName: string
  onDismiss: () => void
}

export function DailySummaryCard({
  summary,
  petName,
  onDismiss,
}: DailySummaryCardProps) {
  const { t } = useTranslation('pet')
  const { t: tCatalog } = useTranslation('catalog')
  const advice = summary.advice
  const itemId = advice?.item_id
  const message =
    summary.generated_by === 'template'
      ? summaryMessage(tCatalog, summary.facts)
      : summary.message

  return (
    <section
      aria-labelledby="pet-summary-title"
      data-testid="daily-summary"
      className="relative flex flex-col gap-3 rounded-xl bg-card px-4 py-4 ring-1 ring-foreground/10"
    >
      <div className="flex items-start justify-between gap-3">
        <h2
          id="pet-summary-title"
          className="text-xs font-medium tracking-wide text-muted-foreground uppercase"
        >
          {t('summary.title', { name: petName })}
        </h2>
        <Button
          size="icon-xs"
          variant="ghost"
          onClick={onDismiss}
          aria-label={t('actions.dismiss')}
        >
          <X aria-hidden="true" className="size-4" />
        </Button>
      </div>

      <blockquote className="border-l-2 border-primary/40 pl-3 text-sm leading-relaxed text-foreground">
        {message}
      </blockquote>

      <dl className="grid grid-cols-3 gap-2 text-center">
        <SummaryStat
          label={t('summary.xpEarned')}
          value={`+${summary.facts.total_xp}`}
        />
        <SummaryStat
          label={t('stats.level')}
          value={String(summary.facts.level)}
        />
        <SummaryStat
          label={t('streak.label')}
          value={String(summary.facts.streak_days)}
        />
      </dl>

      {advice !== null && advice !== undefined && (
        <div className="flex flex-col gap-2 rounded-lg bg-muted px-3 py-2.5">
          <p className="text-sm text-foreground">
            {adviceText(tCatalog, advice)}
          </p>
          {typeof itemId === 'string' && itemId !== '' && (
            <Button
              size="sm"
              variant="outline"
              render={<Link to={`/items/${itemId}`} />}
            >
              {t('summary.fix')}
            </Button>
          )}
        </div>
      )}
    </section>
  )
}

function SummaryStat({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-lg bg-muted px-2 py-2">
      <dt className="text-xs text-muted-foreground">{label}</dt>
      <dd className="font-mono text-base font-semibold tabular-nums text-foreground">
        {value}
      </dd>
    </div>
  )
}
