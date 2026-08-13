import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { ImageOff, Sparkles, Trophy } from 'lucide-react'
import { Button } from '#/components/ui'
import { translateApiError } from '#/api'
import { formatPrice } from '#/features/items/lib'
import type { BukovkiListing, RaccoonJumpReveal } from '#/features/games/types'
import type { CSSProperties } from 'react'

interface RaccoonJumpResultProps {
  reveal: RaccoonJumpReveal | null
  minStreakScore: number
  isSubmitting: boolean
  submitError: unknown
  onRetry: () => void
}

export function RaccoonJumpResult({
  reveal,
  minStreakScore,
  isSubmitting,
  submitError,
  onRetry,
}: RaccoonJumpResultProps) {
  const { t } = useTranslation(['games', 'errors'])

  if (isSubmitting) {
    return (
      <p role="status" className="text-center text-sm text-muted-foreground">
        {t('games:raccoonjump.saving')}
      </p>
    )
  }

  if (submitError) {
    return (
      <div className="flex flex-col items-center gap-2">
        <p role="alert" className="text-center text-sm text-destructive">
          {t('games:raccoonjump.saveFailed')}: {translateApiError(submitError, t)}
        </p>

        <Button type="button" variant="ghost" size="sm" onClick={onRetry}>
          {t('games:raccoonjump.retrySave')}
        </Button>
      </div>
    )
  }

  if (!reveal) return null

  const listings = reveal.listings ?? []

  return (
    <div className="flex flex-col gap-3">
      <p className="flex items-center justify-center gap-1.5 text-sm text-muted-foreground tabular-nums">
        <Trophy className="size-4" aria-hidden="true" />
        {t('games:raccoonjump.best', { score: reveal.best_score })}
      </p>

      {reveal.new_best && (
        <p className="flex items-center justify-center gap-1.5 rounded-xl bg-primary-subtle py-2 text-sm font-semibold text-primary-subtle-foreground">
          <Sparkles className="size-4" aria-hidden="true" />
          {t('games:raccoonjump.newBest')}
        </p>
      )}

      <p className="text-center text-xs text-muted-foreground">
        {reveal.counts_toward_streak
          ? t('games:raccoonjump.counted')
          : t('games:raccoonjump.notCounted', { count: minStreakScore })}
      </p>

      {listings.length > 0 && <CollectedListings listings={listings} />}
    </div>
  )
}

interface CollectedListingsProps {
  listings: BukovkiListing[]
}

function CollectedListings({ listings }: CollectedListingsProps) {
  const { t, i18n } = useTranslation(['games', 'items'])

  return (
    <section className="flex flex-col gap-3 rounded-2xl bg-card p-4 ring-1 ring-border">
      <h2 className="text-sm font-semibold">
        {t('games:raccoonjump.listingsTitle')}
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
