import { Sparkles, X } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { RewardCodeReveal } from '#/features/rewards'
import { cn } from '#/lib/utils'
import type { LotteryRun } from '#/features/weekly-lottery/types'

export function LotteryOutcome({
  run,
  nextAvailable,
}: {
  run: LotteryRun
  nextAvailable: string
}) {
  const { t, i18n } = useTranslation('weeklyLottery')
  const won = run.state === 'won'
  const nextDate = new Intl.DateTimeFormat(i18n.language, {
    dateStyle: 'long',
    timeStyle: 'short',
  }).format(new Date(nextAvailable))

  return (
    <section
      role="status"
      className={cn(
        'lottery-outcome flex flex-col gap-3 rounded-2xl p-4 ring-1 sm:p-5',
        won
          ? 'bg-primary-subtle ring-primary/30'
          : 'bg-muted ring-border',
      )}
    >
      <div className="flex items-start gap-3">
        <span
          className={cn(
            'flex size-10 shrink-0 items-center justify-center rounded-full',
            won
              ? 'bg-primary text-primary-foreground'
              : 'bg-background text-muted-foreground',
          )}
        >
          {won ? (
            <Sparkles className="size-5" aria-hidden="true" />
          ) : (
            <X className="size-5" aria-hidden="true" />
          )}
        </span>
        <div className="flex min-w-0 flex-col gap-1">
          <h2 className="text-base font-semibold">
            {won ? t('outcome.won') : t('outcome.lost')}
          </h2>
          <p className="text-sm text-muted-foreground">
            {won && run.prize ? run.prize.title : t('outcome.lostHint')}
          </p>
        </div>
      </div>

      {won && run.prize && <RewardCodeReveal code={run.prize.code} />}

      <p className="text-xs text-muted-foreground">
        {t('outcome.next', { date: nextDate })}
      </p>
    </section>
  )
}
