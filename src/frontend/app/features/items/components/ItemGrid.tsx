import { useTranslation } from 'react-i18next'
import { ItemCard } from './ItemCard'
import { LoadMore } from './LoadMore'
import type { ItemListEntry } from '#/features/items/types'

interface ItemGridProps {
  items: ItemListEntry[]
  favoriteIds?: Set<string>
  showFavorite?: boolean
  hasNextPage: boolean
  isFetchingNextPage: boolean
  onLoadMore: () => void
}

export function ItemGrid({
  items,
  favoriteIds,
  showFavorite = true,
  hasNextPage,
  isFetchingNextPage,
  onLoadMore,
}: ItemGridProps) {
  const { t } = useTranslation('items')

  return (
    <div className="flex flex-col gap-6">
      <ul
        aria-label={t('title')}
        className="grid list-none auto-rows-fr grid-cols-1 items-stretch gap-4 sm:grid-cols-2 lg:grid-cols-3"
      >
        {items.map((item) => (
          <li key={item.id} className="flex h-full">
            <ItemCard
              item={item}
              isFavorite={favoriteIds?.has(item.id) ?? false}
              showFavorite={showFavorite}
            />
          </li>
        ))}
      </ul>

      <LoadMore
        hasNextPage={hasNextPage}
        isFetching={isFetchingNextPage}
        onLoadMore={onLoadMore}
      />
    </div>
  )
}
