import { PlusOutlined } from '@ant-design/icons';
import { Button } from 'antd';
import { useUnit } from 'effector-react';
import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';

import { useInfiniteMyItems } from '@/entities/item';
import { $isAuthenticated } from '@/entities/session';
import { $filters, ItemsFilterPanel } from '@/features/items-filter';
import { ROUTES } from '@/shared/config/routes';
import { useDebouncedValue } from '@/shared/lib/use-debounced-value';
import { EmptyState, ErrorState, PageSkeleton } from '@/shared/ui';
import { ItemTable } from '@/widgets/item-table';

import './my-items-page.css';

export function MyItemsPage() {
  const { t } = useTranslation(['item', 'common']);
  const isAuthenticated = useUnit($isAuthenticated);
  const filters = useUnit($filters);
  const search = useDebouncedValue(filters.search);

  const query = useInfiniteMyItems(
    {
      search,
      ...(filters.status === null ? {} : { status: filters.status }),
    },
    isAuthenticated,
  );

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
      <header className="my-items-head">
        <h1 className="app-page-title">{t('item:list.myTitle')}</h1>

        <Link to={ROUTES.itemCreate}>
          <Button type="primary" icon={<PlusOutlined />}>
            {t('common:actions.create')}
          </Button>
        </Link>
      </header>

      <ItemsFilterPanel showStatusFilter />

      {query.isPending ? (
        <PageSkeleton />
      ) : items.length === 0 ? (
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
        <ItemTable items={items} loading={false} />
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
    </div>
  );
}
