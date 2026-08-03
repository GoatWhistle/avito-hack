import type { TableColumnsType } from 'antd';
import type { TFunction } from 'i18next';
import { Link } from 'react-router-dom';

import { ItemStatusTag, type ItemListEntry } from '@/entities/item';
import { ROUTES } from '@/shared/config/routes';
import { formatDate, formatPrice } from '@/shared/lib/format';

interface BuildColumnsParams {
  t: TFunction<'item'>;
  locale: string;
  showOwner: boolean;
}

export function buildItemColumns({
  t,
  locale,
  showOwner,
}: BuildColumnsParams): TableColumnsType<ItemListEntry> {
  const columns: TableColumnsType<ItemListEntry> = [
    {
      key: 'title',
      title: t('columns.title'),
      dataIndex: 'title',
      render: (_value, record) => <Link to={ROUTES.itemDetail(record.id)}>{record.title}</Link>,
    },
    {
      key: 'price',
      title: t('columns.price'),
      dataIndex: 'priceKopeks',
      width: 160,
      render: (_value, record) => formatPrice(record.priceKopeks, locale),
    },
    {
      key: 'status',
      title: t('columns.status'),
      dataIndex: 'status',
      width: 160,
      render: (_value, record) => <ItemStatusTag status={record.status} />,
    },
    {
      key: 'createdAt',
      title: t('columns.createdAt'),
      dataIndex: 'createdAt',
      width: 160,
      render: (_value, record) => formatDate(record.createdAt, locale),
    },
  ];

  if (!showOwner) {
    return columns;
  }

  return [
    ...columns.slice(0, 1),
    {
      key: 'owner',
      title: t('columns.owner'),
      dataIndex: 'ownerName',
      width: 200,
    },
    ...columns.slice(1),
  ];
}
