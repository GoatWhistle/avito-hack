import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { formatDate, formatPrice } from '#/features/items/lib'
import { ItemPhotoThumb, ItemStatusBadge } from '#/features/items/components'
import { FavoriteButton } from './FavoriteButton'
import type { FavoriteEntry } from '#/features/items/types'

interface FavoriteCardProps {
  entry: FavoriteEntry
}

export function FavoriteCard({ entry }: FavoriteCardProps) {
  const { t, i18n } = useTranslation('items')

  return (
    <article className="relative flex gap-3 rounded-xl bg-card p-3 ring-1 ring-foreground/10 transition-shadow focus-within:ring-2 focus-within:ring-ring hover:shadow-md">
      <div className="w-24 shrink-0 sm:w-32">
        <ItemPhotoThumb url={entry.photo_url} alt={entry.title} />
      </div>

      <div className="flex min-w-0 flex-1 flex-col gap-1.5">
        <h3 className="text-sm leading-snug font-medium break-words">
          <Link
            to={`/items/${entry.item_id}`}
            className="outline-none after:absolute after:inset-0 after:content-['']"
          >
            {entry.title}
          </Link>
        </h3>

        <p className="text-base font-semibold">
          {t('price', { value: formatPrice(entry.price, i18n.language) })}
        </p>

        <div className="mt-auto flex flex-wrap items-center gap-2">
          <ItemStatusBadge status={entry.status} />
          <span className="text-xs text-muted-foreground">
            {formatDate(entry.added_at, i18n.language)}
          </span>
        </div>
      </div>

      <div className="z-10 shrink-0 self-start">
        <FavoriteButton
          itemId={entry.item_id}
          isFavorite
          variant="outline"
          entry={entry}
        />
      </div>
    </article>
  )
}
