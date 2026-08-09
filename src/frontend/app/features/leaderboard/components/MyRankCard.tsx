import { TrendingUp } from 'lucide-react'
import { useTranslation } from 'react-i18next'

export function MyRankCard({ rank }: { rank: number | null }) {
  const { t } = useTranslation('leaderboard')

  if (rank === null) {
    return (
      <p
        data-testid="my-rank-card"
        className="rounded-xl bg-muted/50 px-4 py-3 text-sm text-muted-foreground"
      >
        {t('notRanked')}
      </p>
    )
  }

  return (
    <div
      data-testid="my-rank-card"
      className="flex items-center justify-between gap-3 rounded-xl bg-primary-subtle px-4 py-3 ring-1 ring-primary/20"
    >
      <span className="flex items-center gap-2 text-sm font-medium text-primary-subtle-foreground">
        <TrendingUp aria-hidden="true" className="size-4" />
        {t('myRankShort')}
      </span>
      <span className="font-mono text-2xl leading-none font-bold tabular-nums text-primary-subtle-foreground">
        {t('rankValue', { rank })}
      </span>
    </div>
  )
}
