import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { ImageOff } from 'lucide-react'
import { formatPrice } from '#/features/items/lib'
import type { BukovkiListing } from '#/features/games/types'
import type { CSSProperties } from 'react'

interface BukovkiListingsProps {
  listings: BukovkiListing[]
}

export function BukovkiListings({ listings }: BukovkiListingsProps) {
  const { t, i18n } = useTranslation(['games', 'items'])

  if (listings.length === 0) {
    return (
      <section className="flex flex-col gap-3 rounded-2xl bg-card p-4 text-center ring-1 ring-border">
        <h2 className="text-sm font-semibold">
          {t('games:bukovki.listingsTitle')}
        </h2>

        <p className="text-sm text-muted-foreground">
          {t('games:bukovki.listingsEmpty')}
        </p>

        <Link
          to="/items"
          className="text-sm font-semibold text-primary no-underline hover:underline"
        >
          {t('games:bukovki.listingsBrowseAll')}
        </Link>
      </section>
    )
  }

  return (
    <section className="flex flex-col gap-3 rounded-2xl bg-card p-4 ring-1 ring-border">
      <h2 className="text-sm font-semibold">
        {t('games:bukovki.listingsTitle')}
      </h2>

      <ul className="grid list-none gap-2 sm:grid-cols-2">
        {listings.map((listing, index) => (
          <li
            key={listing.display_id}
            className="game-stagger"
            style={{ '--game-stagger-index': index } as CSSProperties}
          >
            <Link
              to={`/items/${listing.display_id}`}
              className="flex items-center gap-3 rounded-xl bg-muted/40 p-2.5 no-underline ring-1 ring-border transition-colors hover:bg-muted"
            >
              <span className="flex size-12 shrink-0 items-center justify-center overflow-hidden rounded-lg bg-muted">
                {listing.photo_url ? (
                  <img
                    src={listing.photo_url}
                    alt=""
                    className="size-full object-cover"
                  />
                ) : (
                  <ImageOff
                    className="size-4 text-muted-foreground"
                    aria-hidden="true"
                  />
                )}
              </span>

              <span className="flex min-w-0 flex-1 flex-col">
                <span className="line-clamp-2 text-xs leading-snug font-medium">
                  {listing.title}
                </span>
                <span className="text-sm font-semibold tabular-nums">
                  {t('items:price', {
                    value: formatPrice(listing.price_kopeks, i18n.language),
                  })}
                </span>
              </span>
            </Link>
          </li>
        ))}
      </ul>
    </section>
  )
}
