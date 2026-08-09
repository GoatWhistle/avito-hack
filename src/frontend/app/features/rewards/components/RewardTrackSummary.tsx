import { useTranslation } from 'react-i18next'
import { Card, CardContent, CardHeader, Progress } from '#/components/ui'

export interface RewardTrackSummaryProps {
  claimed: number
  total: number
  percent: number
}

export function RewardTrackSummary({
  claimed,
  total,
  percent,
}: RewardTrackSummaryProps) {
  const { t } = useTranslation('rewards')

  return (
    <Card
      size="sm"
      className="relative bg-achievement-earned ring-achievement-earned-border/30"
    >
      <span
        aria-hidden="true"
        className="pointer-events-none absolute -top-4 -right-4 size-20 rounded-full bg-reward-glow/15 blur-xl"
      />

      <CardHeader className="flex flex-row items-center justify-between gap-2 space-y-0">
        <span className="rounded-lg bg-achievement-earned-tile/70 px-2.5 py-0.5 text-xs font-semibold text-achievement-earned-border ring-1 ring-achievement-earned-border/30">
          {t('track.title')}
        </span>
        <span aria-hidden="true" className="text-xl">
          🎁
        </span>
      </CardHeader>

      <CardContent className="flex flex-col gap-2.5">
        <p
          id="reward-track-title"
          className="text-sm font-bold text-foreground"
        >
          {t('track.counter', { claimed, total })}
        </p>

        <div className="flex items-center gap-2">
          <Progress
            className="h-2 flex-1"
            value={percent}
            indicatorClassName="bg-achievement-earned-border"
            label={t('track.progressLabel', { percent })}
          />
          <span className="shrink-0 text-[11px] font-medium tabular-nums text-muted-foreground">
            {percent}%
          </span>
        </div>
      </CardContent>
    </Card>
  )
}
