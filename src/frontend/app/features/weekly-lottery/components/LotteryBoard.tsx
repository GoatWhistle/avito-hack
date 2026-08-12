import { Gift, LoaderCircle } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { cn } from '#/lib/utils'
import { LotterySymbolIcon } from './LotterySymbolIcon'
import type {
  LotteryRun,
  LotterySymbol,
} from '#/features/weekly-lottery/types'

interface LotteryBoardProps {
  run: LotteryRun
  revealingSlot?: number
  isRevealing: boolean
  onReveal: (slot: number) => void
}

const winningSymbolOf = (run: LotteryRun): LotterySymbol | null => {
  if (run.state !== 'won') return null

  const counts = new Map<LotterySymbol, number>()
  for (const slot of run.slots) {
    if (!slot.opened || !slot.symbol) continue
    counts.set(slot.symbol, (counts.get(slot.symbol) ?? 0) + 1)
  }

  for (const [symbol, count] of counts) {
    if (count >= 3) return symbol
  }

  return null
}

export function LotteryBoard({
  run,
  revealingSlot,
  isRevealing,
  onReveal,
}: LotteryBoardProps) {
  const { t } = useTranslation('weeklyLottery')
  const winningSymbol = winningSymbolOf(run)

  return (
    <div
      className="grid grid-cols-3 gap-2 sm:gap-3"
      aria-label={t('board.label')}
    >
      {run.slots.map((slot) => {
        const revealing = isRevealing && revealingSlot === slot.index
        const matched =
          slot.opened &&
          slot.symbol !== undefined &&
          slot.symbol === winningSymbol
        const disabled = slot.opened || isRevealing || run.state !== 'active'
        const label = slot.opened
          ? t('board.opened', {
              position: slot.index + 1,
              symbol: t(`symbols.${slot.symbol}`),
            })
          : t('board.closed', { position: slot.index + 1 })

        return (
          <button
            key={slot.index}
            type="button"
            aria-label={label}
            aria-busy={revealing}
            disabled={disabled}
            className={cn(
              'lottery-slot group relative aspect-square overflow-hidden rounded-2xl border text-foreground shadow-sm outline-none transition-all focus-visible:ring-3 focus-visible:ring-ring/50',
              slot.opened
                ? 'lottery-slot-open border-primary/25 bg-card'
                : 'border-primary/20 bg-primary-subtle hover:-translate-y-0.5 hover:border-primary/40 hover:shadow-md',
              matched && 'lottery-slot-match border-primary ring-2 ring-primary/40',
            )}
            onClick={() => onReveal(slot.index)}
          >
            <span className="absolute inset-0 flex flex-col items-center justify-center gap-2 p-2">
              {revealing ? (
                <LoaderCircle
                  className="size-7 animate-spin text-primary"
                  aria-hidden="true"
                />
              ) : slot.opened && slot.symbol ? (
                <>
                  <LotterySymbolIcon
                    symbol={slot.symbol}
                    className="lottery-symbol size-8 text-primary sm:size-10"
                  />
                  <span className="text-[0.6875rem] leading-tight font-medium sm:text-xs">
                    {t(`symbols.${slot.symbol}`)}
                  </span>
                </>
              ) : (
                <Gift
                  className="size-7 text-primary transition-transform group-hover:scale-110 sm:size-9"
                  aria-hidden="true"
                />
              )}
            </span>
          </button>
        )
      })}
    </div>
  )
}
