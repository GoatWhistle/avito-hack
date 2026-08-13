import { ChevronRight, Gift, LockKeyhole, Sparkles } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { useWeeklyLotteryState } from '#/features/weekly-lottery/hooks'

export function WeeklyLotteryCard() {
  const { t, i18n } = useTranslation('weeklyLottery')
  const query = useWeeklyLotteryState()
  const run = query.data?.run
  const finished = run?.state === 'won' || run?.state === 'lost'
  const nextDate = query.data
    ? new Intl.DateTimeFormat(i18n.language, {
        day: 'numeric',
        month: 'short',
      }).format(new Date(query.data.next_available_at))
    : ''

  const status = query.isPending
    ? t('card.loading')
    : run?.state === 'won'
      ? t('card.won')
      : run?.state === 'lost'
        ? t('card.finished', { date: nextDate })
        : run?.state === 'active'
          ? t('card.continue')
          : t('card.available')

  return (
    <section className="px-3 pt-3">
      <Link
        to="/play/weekly-lottery"
        className="group flex items-center gap-3 rounded-xl bg-card p-3 no-underline ring-1 ring-border transition-colors hover:bg-muted"
      >
        <span className="flex size-10 shrink-0 items-center justify-center rounded-lg bg-primary-subtle text-primary">
          {finished ? (
            run?.state === 'won' ? (
              <Sparkles className="size-5" aria-hidden="true" />
            ) : (
              <LockKeyhole className="size-5" aria-hidden="true" />
            )
          ) : (
            <Gift className="size-5" aria-hidden="true" />
          )}
        </span>
        <span className="flex min-w-0 flex-1 flex-col">
          <span className="text-sm font-medium">{t('title')}</span>
          <span className="truncate text-xs text-muted-foreground">
            {status}
          </span>
        </span>
        <ChevronRight
          className="size-4 shrink-0 text-muted-foreground transition-transform group-hover:translate-x-0.5"
          aria-hidden="true"
        />
      </Link>
    </section>
  )
}
