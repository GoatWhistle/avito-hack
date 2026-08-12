import { useCallback, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { Button } from '#/components/ui'
import { translateApiError } from '#/api'
import { useFavoriteIds } from '#/features/favorites/hooks'
import { useSession } from '#/features/auth/session'
import { flattenPages, useItemsQuery } from '#/features/items/hooks'
import { useDebouncedValue } from '#/features/items/hooks/useDebouncedValue'
import { ItemFilters } from './ItemFilters'
import { ItemGrid } from './ItemGrid'
import { PetHintBanner } from './PetHintBanner'
import { EmptyState, ErrorState, ItemsSkeleton } from './ListStates'
import type {
  ItemCategory,
  ItemCondition,
  ItemSort,
  ItemStatus,
} from '#/features/items/types'

const catalogStatuses: readonly ItemStatus[] = ['published', 'sold']

export function ItemsScreen() {
  const { t } = useTranslation(['items', 'errors'])
  const { isAuthenticated } = useSession()
  const [search, setSearch] = useState('')
  const [status, setStatus] = useState<ItemStatus | ''>('')
  const [category, setCategory] = useState<ItemCategory | ''>('')
  const [condition, setCondition] = useState<ItemCondition | ''>('')
  const [sort, setSort] = useState<ItemSort>('newest')
  const debouncedSearch = useDebouncedValue(search)
  const favoriteIds = useFavoriteIds()

  const query = useItemsQuery({
    status,
    search: debouncedSearch,
    category,
    condition,
    sort,
  })
  const items = flattenPages(query.data?.pages)

  const loadMore = useCallback(() => {
    if (query.hasNextPage && !query.isFetchingNextPage) {
      void query.fetchNextPage()
    }
  }, [query])

  return (
    <section className="flex flex-col gap-6">
      <header className="flex flex-wrap items-center justify-between gap-3">
        <h1 className="text-xl font-semibold">{t('title')}</h1>
        {isAuthenticated && (
          <Button render={<Link to="/items/new" />}>
            {t('actions.create')}
          </Button>
        )}
      </header>

      <PetHintBanner />

      <ItemFilters
        search={search}
        onSearchChange={setSearch}
        status={status}
        onStatusChange={setStatus}
        statuses={catalogStatuses}
        category={category}
        onCategoryChange={setCategory}
        condition={condition}
        onConditionChange={setCondition}
        sort={sort}
        onSortChange={setSort}
      />

      {query.isPending && <ItemsSkeleton />}

      {query.isError && (
        <ErrorState
          message={translateApiError(query.error, t)}
          onRetry={() => void query.refetch()}
        />
      )}

      {query.isSuccess && items.length === 0 && (
        <EmptyState
          title={t('empty.list')}
          hint={t('empty.listHint')}
          actionLabel={isAuthenticated ? t('actions.create') : undefined}
          actionTo={isAuthenticated ? '/items/new' : undefined}
        />
      )}

      {query.isSuccess && items.length > 0 && (
        <ItemGrid
          items={items}
          favoriteIds={favoriteIds}
          showFavorite={isAuthenticated}
          hasNextPage={query.hasNextPage}
          isFetchingNextPage={query.isFetchingNextPage}
          onLoadMore={loadMore}
        />
      )}
    </section>
  )
}
