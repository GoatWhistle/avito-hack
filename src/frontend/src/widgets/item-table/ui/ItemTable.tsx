import { Table } from 'antd';
import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';

import type { ItemListEntry } from '@/entities/item';
import { EmptyState } from '@/shared/ui';

import { buildItemColumns } from './columns';

const VIRTUAL_THRESHOLD = 100;

interface ItemTableProps {
  items: ItemListEntry[];
  loading: boolean;
  showOwner?: boolean;
}

const VIRTUAL_SCROLL_HEIGHT = 600;

export function ItemTable({ items, loading, showOwner = false }: ItemTableProps) {
  const { t, i18n } = useTranslation('item');

  const columns = useMemo(
    () => buildItemColumns({ t, locale: i18n.language, showOwner }),
    [t, i18n.language, showOwner],
  );

  const isVirtual = items.length > VIRTUAL_THRESHOLD;

  return (
    <Table<ItemListEntry>
      rowKey="id"
      columns={columns}
      dataSource={items}
      loading={loading}
      pagination={false}
      virtual={isVirtual}
      {...(isVirtual ? { scroll: { y: VIRTUAL_SCROLL_HEIGHT } } : {})}
      locale={{ emptyText: <EmptyState description={t('list.empty')} /> }}
    />
  );
}
