import { useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import { useFavoritesQuery } from '#/features/favorites/hooks'
import { flattenPages } from '#/features/items/hooks'
import {
  EmptyState,
  ErrorState,
  ItemsSkeleton,
  LoadMore,
} from '#/features/items/components'
import { FavoriteCard } from './FavoriteCard'

export function FavoritesScreen() {
  const { t } = useTranslation('items')

  const query = useFavoritesQuery()
  const entries = flattenPages(query.data?.pages)

  const loadMore = useCallback(() => {
    if (query.hasNextPage && !query.isFetchingNextPage) {
      void query.fetchNextPage()
    }
  }, [query])

  return (
    <section className="flex flex-col gap-6">
      <h1 className="text-xl font-semibold">{t('favoritesTitle')}</h1>

      {query.isPending && (
        <ItemsSkeleton count={3} className="sm:grid-cols-1" />
      )}

      {query.isError && (
        <ErrorState
          message={
            query.error instanceof Error ? query.error.message : undefined
          }
          onRetry={() => void query.refetch()}
        />
      )}

      {query.isSuccess && entries.length === 0 && (
        <EmptyState
          title={t('empty.favorites')}
          hint={t('empty.favoritesHint')}
          actionLabel={t('title')}
          actionTo="/items"
        />
      )}

      {query.isSuccess && entries.length > 0 && (
        <div className="flex flex-col gap-4">
          <ul
            aria-label={t('favoritesTitle')}
            className="flex list-none flex-col gap-3"
          >
            {entries.map((entry) => (
              <li key={entry.item_id}>
                <FavoriteCard entry={entry} />
              </li>
            ))}
          </ul>

          <LoadMore
            hasNextPage={query.hasNextPage}
            isFetching={query.isFetchingNextPage}
            onLoadMore={loadMore}
          />
        </div>
      )}
    </section>
  )
}
