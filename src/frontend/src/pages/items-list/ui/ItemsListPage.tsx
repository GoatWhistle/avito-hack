import { Button } from 'antd';
import { useUnit } from 'effector-react';
import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';

import { useInfiniteItems } from '@/entities/item';
import { $filters, ItemsFilterPanel } from '@/features/items-filter';
import { ROUTES } from '@/shared/config/routes';
import { useDebouncedValue } from '@/shared/lib/use-debounced-value';
import { EmptyState, ErrorState, PageSkeleton } from '@/shared/ui';
import { ItemTable } from '@/widgets/item-table';

export function ItemsListPage() {
  const { t } = useTranslation(['item', 'common']);
  const filters = useUnit($filters);
  const search = useDebouncedValue(filters.search);

  const query = useInfiniteItems({ status: 'published', search });

  const items = useMemo(() => query.data?.pages.flatMap((page) => page.items) ?? [], [query.data]);

  if (query.isError) {
    return (
      <ErrorState
        error={query.error}
        onRetry={() => {
          void query.refetch();
        }}
      />
    );
  }

  return (
    <div className="app-stack">
      <header>
        <h1 className="app-page-title">{t('item:list.title')}</h1>
      </header>

      <ItemsFilterPanel />

      {query.isPending ? (
        <PageSkeleton />
      ) : (
        <>
          <p className="app-muted">{t('item:list.count', { count: items.length })}</p>

          {items.length === 0 ? (
            <EmptyState
              description={t('item:list.empty')}
              hint={t('item:list.emptyHint')}
              action={
                <Link to={ROUTES.itemCreate}>
                  <Button type="primary">{t('item:list.create')}</Button>
                </Link>
              }
            />
          ) : (
            <ItemTable items={items} loading={false} showOwner />
          )}

          {query.hasNextPage && (
            <Button
              loading={query.isFetchingNextPage}
              onClick={() => {
                void query.fetchNextPage();
              }}
            >
              {t('common:pagination.loadMore')}
            </Button>
          )}
        </>
      )}
    </div>
  );
}
