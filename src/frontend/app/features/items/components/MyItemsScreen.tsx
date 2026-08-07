import { useCallback, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { Button } from '#/components/ui'
import { flattenPages, useMyItemsQuery } from '#/features/items/hooks'
import { ItemFilters } from './ItemFilters'
import { MyItemRow } from './MyItemRow'
import { LoadMore } from './LoadMore'
import { EmptyState, ErrorState, ItemsSkeleton } from './ListStates'
import type { ItemStatus } from '#/features/items/types'

export function MyItemsScreen() {
  const { t } = useTranslation('items')
  const [status, setStatus] = useState<ItemStatus | ''>('')

  const query = useMyItemsQuery({ status })
  const items = flattenPages(query.data?.pages)

  const loadMore = useCallback(() => {
    if (query.hasNextPage && !query.isFetchingNextPage) {
      void query.fetchNextPage()
    }
  }, [query])

  return (
    <section className="flex flex-col gap-6">
      <header className="flex flex-wrap items-center justify-between gap-3">
        <h1 className="text-xl font-semibold">{t('myTitle')}</h1>
        <Button render={<Link to="/items/new" />}>{t('actions.create')}</Button>
      </header>

      <ItemFilters status={status} onStatusChange={setStatus} />

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

      {query.isSuccess && items.length === 0 && (
        <EmptyState
          title={t('empty.mine')}
          hint={t('empty.mineHint')}
          actionLabel={t('actions.create')}
          actionTo="/items/new"
        />
      )}

      {query.isSuccess && items.length > 0 && (
        <div className="flex flex-col gap-4">
          <ul
            aria-label={t('myTitle')}
            className="flex list-none flex-col gap-3"
          >
            {items.map((item) => (
              <li key={item.id}>
                <MyItemRow item={item} />
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
