import { Card, Space, Typography } from 'antd';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';

import { ROUTES } from '@/shared/config/routes';
import { formatDate, formatPrice } from '@/shared/lib/format';

import type { ItemListEntry } from '../model/types';

import { ItemStatusTag } from './ItemStatusTag';


interface ItemCardProps {
  item: ItemListEntry;
}

export function ItemCard({ item }: ItemCardProps) {
  const { i18n } = useTranslation();

  return (
    <Card size="small" hoverable>
      <Space direction="vertical" size="small" style={{ width: '100%' }}>
        <Link to={ROUTES.itemDetail(item.id)}>
          <Typography.Title level={5} style={{ margin: 0 }}>
            {item.title}
          </Typography.Title>
        </Link>

        <Typography.Text strong>{formatPrice(item.priceKopeks, i18n.language)}</Typography.Text>

        <Space size="small" wrap>
          <ItemStatusTag status={item.status} />
          <Typography.Text type="secondary">
            {formatDate(item.createdAt, i18n.language)}
          </Typography.Text>
        </Space>

        {item.ownerName !== '' && (
          <Typography.Text type="secondary">{item.ownerName}</Typography.Text>
        )}
      </Space>
    </Card>
  );
}
