import { useTranslation } from 'react-i18next'
import { LotterySymbolIcon } from './LotterySymbolIcon'
import type { LotteryPrizeCatalogItem } from '#/features/weekly-lottery/types'

export function LotteryPrizeLegend({
  prizes,
}: {
  prizes: LotteryPrizeCatalogItem[]
}) {
  const { t } = useTranslation('weeklyLottery')

  return (
    <section className="flex flex-col gap-2" aria-labelledby="lottery-prizes">
      <h2 id="lottery-prizes" className="text-sm font-semibold">
        {t('prizes.title')}
      </h2>
      <ul className="grid list-none gap-2 sm:grid-cols-2">
        {prizes.map((prize) => (
          <li
            key={prize.id}
            className="flex items-center gap-3 rounded-xl bg-card p-3 ring-1 ring-border"
          >
            <span className="flex size-9 shrink-0 items-center justify-center rounded-lg bg-primary-subtle text-primary">
              <LotterySymbolIcon symbol={prize.symbol} className="size-5" />
            </span>
            <span className="flex min-w-0 flex-col">
              <span className="text-sm font-medium">{prize.title}</span>
              <span className="truncate text-xs text-muted-foreground">
                {prize.description}
              </span>
            </span>
          </li>
        ))}
      </ul>
    </section>
  )
}
