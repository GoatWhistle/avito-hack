import { Button, Space, Typography } from 'antd';
import { useUnit } from 'effector-react';
import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';

import { useInfiniteItems } from '@/entities/item';
import { $filters, ItemsFilterPanel } from '@/features/items-filter';
import { useDebouncedValue } from '@/shared/lib/use-debounced-value';
import { ErrorState } from '@/shared/ui';
import { ItemTable } from '@/widgets/item-table';

export function ItemsListPage() {
  const { t } = useTranslation(['item', 'common']);
  const filters = useUnit($filters);
  const search = useDebouncedValue(filters.search);

  const query = useInfiniteItems({ status: 'published', search });

  const items = useMemo(
    () => query.data?.pages.flatMap((page) => page.items) ?? [],
    [query.data],
  );

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
      <Typography.Title level={3}>{t('item:list.title')}</Typography.Title>

      <ItemsFilterPanel />

      <Typography.Text type="secondary">
        {t('item:list.count', { count: items.length })}
      </Typography.Text>

      <ItemTable items={items} loading={query.isPending} showOwner />

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
