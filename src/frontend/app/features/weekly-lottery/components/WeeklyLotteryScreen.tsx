import { Gift, LoaderCircle, Play } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Button } from '#/components/ui'
import { translateApiError } from '#/api'
import { GameScreenShell } from '#/features/games/components'
import {
  useRevealWeeklyLotterySlot,
  useStartWeeklyLottery,
  useWeeklyLotteryPrizes,
  useWeeklyLotteryState,
} from '#/features/weekly-lottery/hooks'
import { LotteryBoard } from './LotteryBoard'
import { LotteryOutcome } from './LotteryOutcome'
import { LotteryPrizeLegend } from './LotteryPrizeLegend'
import { WeeklyLotteryHelp } from './WeeklyLotteryHelp'

export function WeeklyLotteryScreen() {
  const { t } = useTranslation(['weeklyLottery', 'errors'])
  const stateQuery = useWeeklyLotteryState()
  const prizesQuery = useWeeklyLotteryPrizes()
  const start = useStartWeeklyLottery()
  const reveal = useRevealWeeklyLotterySlot()
  const state = stateQuery.data
  const run = state?.run ?? start.data
  const openedCount = run?.slots.filter((slot) => slot.opened).length ?? 0
  const error = stateQuery.error ?? start.error ?? reveal.error

  const revealSlot = (slot: number) => {
    if (!run || run.state !== 'active' || reveal.isPending) return
    reveal.mutate({ runId: run.id, slot })
  }

  return (
    <GameScreenShell
      title={t('weeklyLottery:title')}
      subtitle={t('weeklyLottery:subtitle')}
      backLabel={t('weeklyLottery:back')}
      help={<WeeklyLotteryHelp />}
    >
      {stateQuery.isPending ? (
        <div
          role="status"
          className="flex items-center justify-center gap-2 py-20 text-sm text-muted-foreground"
        >
          <LoaderCircle className="size-4 animate-spin" aria-hidden="true" />
          {t('weeklyLottery:loading')}
        </div>
      ) : (
        <div className="flex flex-col gap-4">
          {!run && (
            <section className="flex flex-col items-center gap-4 rounded-2xl bg-card px-5 py-8 text-center ring-1 ring-border">
              <span className="flex size-14 items-center justify-center rounded-full bg-primary-subtle text-primary">
                <Gift className="size-7" aria-hidden="true" />
              </span>
              <div className="flex flex-col gap-1">
                <h2 className="text-lg font-semibold">
                  {t('weeklyLottery:intro.title')}
                </h2>
                <p className="max-w-md text-sm text-muted-foreground">
                  {t('weeklyLottery:intro.hint')}
                </p>
              </div>
              <Button
                type="button"
                size="lg"
                disabled={start.isPending || !state?.available}
                onClick={() => start.mutate()}
              >
                {start.isPending ? (
                  <LoaderCircle
                    className="size-4 animate-spin"
                    aria-hidden="true"
                  />
                ) : (
                  <Play className="size-4" aria-hidden="true" />
                )}
                {t('weeklyLottery:intro.start')}
              </Button>
            </section>
          )}

          {run && (
            <section className="flex flex-col gap-3">
              <div className="flex items-center justify-between gap-3">
                <h2 className="text-sm font-semibold">
                  {t('weeklyLottery:board.title')}
                </h2>
                <p className="text-xs text-muted-foreground tabular-nums">
                  {t('weeklyLottery:board.progress', { count: openedCount })}
                </p>
              </div>
              <LotteryBoard
                run={run}
                revealingSlot={reveal.variables?.slot}
                isRevealing={reveal.isPending}
                onReveal={revealSlot}
              />
            </section>
          )}

          {run && run.state !== 'active' && state && (
            <LotteryOutcome
              run={run}
              nextAvailable={state.next_available_at}
            />
          )}

          {error && (
            <div
              role="alert"
              className="flex items-center justify-between gap-3 rounded-xl bg-destructive-subtle p-3 text-sm text-destructive"
            >
              <span>{translateApiError(error, t)}</span>
              {stateQuery.isError && (
                <Button
                  type="button"
                  size="sm"
                  variant="outline"
                  onClick={() => void stateQuery.refetch()}
                >
                  {t('weeklyLottery:retry')}
                </Button>
              )}
            </div>
          )}

          {prizesQuery.data && (
            <LotteryPrizeLegend prizes={prizesQuery.data} />
          )}
        </div>
      )}
    </GameScreenShell>
  )
}
