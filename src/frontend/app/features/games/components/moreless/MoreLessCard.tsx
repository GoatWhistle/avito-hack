import { useTranslation } from 'react-i18next'
import { ArrowUpRight, ImageOff } from 'lucide-react'
import { formatPrice } from '#/features/items/lib'
import { cn } from '#/lib/utils'
import type { MoreLessItem } from '#/features/games/types'

type Highlight = 'correct' | 'wrong' | null

interface MoreLessCardProps {
  item: MoreLessItem
  revealedPrice?: number
  hint?: string
  highlight?: Highlight
  linkable?: boolean
}

export function MoreLessCard({
  item,
  revealedPrice,
  hint,
  highlight = null,
  linkable = false,
}: MoreLessCardProps) {
  const { t, i18n } = useTranslation(['games', 'items'])
  const price = revealedPrice ?? item.price
  const justRevealed = revealedPrice !== undefined

  return (
    <article
      className={cn(
        'game-card-enter relative aspect-4/3 overflow-hidden rounded-2xl bg-card ring-1 transition-all duration-300 sm:aspect-[4/5]',
        highlight === 'correct' && 'game-card-correct ring-2 ring-primary',
        highlight === 'wrong' && 'game-card-wrong ring-2 ring-destructive',
        highlight === null && 'ring-border',
      )}
    >
      {item.photo_url ? (
        <img
          src={item.photo_url}
          alt=""
          className="absolute inset-0 size-full scale-105 object-cover transition-transform duration-500"
        />
      ) : (
        <div className="absolute inset-0 flex items-center justify-center bg-muted">
          <ImageOff
            className="size-8 text-muted-foreground"
            aria-hidden="true"
          />
        </div>
      )}

      <div className="absolute inset-0 bg-gradient-to-t from-black/90 via-black/40 to-transparent" />

      <div className="absolute inset-x-0 bottom-0 flex flex-col gap-2 p-4 text-white sm:p-5">
        <h3 className="line-clamp-2 text-sm leading-snug font-medium drop-shadow sm:text-base">
          {item.title}
        </h3>

        {price === undefined ? (
          <p className="text-xs text-white/70">{hint}</p>
        ) : (
          <p
            aria-live="polite"
            className={cn(
              'text-2xl leading-none font-bold tabular-nums drop-shadow sm:text-3xl',
              justRevealed && 'game-price-reveal',
            )}
          >
            {t('items:price', { value: formatPrice(price, i18n.language) })}
          </p>
        )}

        {linkable && (
          <a
            href={`/items/${item.item_id}`}
            target="_blank"
            rel="noreferrer"
            className="inline-flex w-fit items-center gap-1 rounded-full bg-white/15 px-2.5 py-1 text-xs font-medium text-white no-underline backdrop-blur-sm transition-colors hover:bg-white/25 focus-visible:ring-3 focus-visible:ring-white/40 focus-visible:outline-none"
          >
            {t('games:round.openListing')}
            <ArrowUpRight className="size-3.5" aria-hidden="true" />
          </a>
        )}
      </div>
    </article>
  )
}
