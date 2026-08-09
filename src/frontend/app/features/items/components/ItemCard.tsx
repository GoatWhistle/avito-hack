import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { formatPrice } from '#/features/items/lib'
import { useItemPhotosQuery } from '#/features/items/hooks'
import { FavoriteButton } from '#/features/favorites/components'
import { ItemStatusBadge } from './ItemStatusBadge'
import { ItemPhotoThumb } from './ItemPhotoThumb'
import type { ItemListEntry } from '#/features/items/types'

interface ItemCardProps {
  item: ItemListEntry
  photoUrl?: string
  isFavorite?: boolean
  showFavorite?: boolean
}

export function ItemCard({
  item,
  photoUrl,
  isFavorite = false,
  showFavorite = true,
}: ItemCardProps) {
  const { t, i18n } = useTranslation('items')
  const photosQuery = useItemPhotosQuery(photoUrl ? undefined : item.id)
  const cover = photoUrl ?? photosQuery.data?.[0]?.url

  return (
    <article className="group/item relative flex h-full w-full flex-col gap-3 rounded-xl bg-card p-3 ring-1 shadow-sm ring-foreground/8 transition-shadow focus-within:ring-2 focus-within:ring-ring hover:shadow-lg dark:shadow-none dark:ring-foreground/10 dark:hover:shadow-md">
      <ItemPhotoThumb url={cover} alt={item.title} />

      <div className="flex flex-1 flex-col gap-2">
        <h3 className="line-clamp-2 min-h-10 text-sm leading-snug font-medium break-words">
          <Link
            to={`/items/${item.id}`}
            className="outline-none after:absolute after:inset-0 after:content-['']"
          >
            {item.title}
          </Link>
        </h3>

        <p className="text-base font-semibold">
          {t('price', { value: formatPrice(item.price, i18n.language) })}
        </p>

        <div className="mt-auto flex flex-wrap items-center gap-2">
          <ItemStatusBadge status={item.status} />
          {item.owner_name && (
            <span className="truncate text-xs text-muted-foreground">
              {item.owner_name}
            </span>
          )}
        </div>
      </div>

      {showFavorite && (
        <div className="absolute top-4 right-4 z-10">
          <FavoriteButton
            itemId={item.id}
            isFavorite={isFavorite}
            variant="outline"
            entry={{
              item_id: item.id,
              owner_id: item.owner_id,
              title: item.title,
              price: item.price,
              status: item.status,
              photo_url: cover,
            }}
          />
        </div>
      )}
    </article>
  )
}
