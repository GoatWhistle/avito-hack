import { PlusOutlined } from '@ant-design/icons';
import { Button, Space, Typography } from 'antd';
import { useUnit } from 'effector-react';
import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';

import { useInfiniteMyItems } from '@/entities/item';
import { $isAuthenticated } from '@/entities/session';
import { $filters, ItemsFilterPanel } from '@/features/items-filter';
import { ROUTES } from '@/shared/config/routes';
import { useDebouncedValue } from '@/shared/lib/use-debounced-value';
import { ErrorState } from '@/shared/ui';
import { ItemTable } from '@/widgets/item-table';

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
    <Space direction="vertical" size="middle" style={{ width: '100%' }}>
      <Space style={{ width: '100%', justifyContent: 'space-between' }}>
        <Typography.Title level={3} style={{ margin: 0 }}>
          {t('item:list.myTitle')}
        </Typography.Title>

        <Link to={ROUTES.itemCreate}>
          <Button type="primary" icon={<PlusOutlined />}>
            {t('common:actions.create')}
          </Button>
        </Link>
      </Space>

      <ItemsFilterPanel showStatusFilter />

      <ItemTable items={items} loading={query.isPending} />

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
    </Space>
  );
}
