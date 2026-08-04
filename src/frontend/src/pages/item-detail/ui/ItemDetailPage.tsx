import { Button, Card, Descriptions, Space, Typography } from 'antd';
import { useUnit } from 'effector-react';
import { useTranslation } from 'react-i18next';
import { Link, useParams } from 'react-router-dom';

import { ItemStatusTag, useItem } from '@/entities/item';
import { $user } from '@/entities/session';
import { ItemStatusActions } from '@/features/item-status';
import { ROUTES } from '@/shared/config/routes';
import { formatDateTime, formatPrice } from '@/shared/lib/format';
import { ErrorState, PageSkeleton } from '@/shared/ui';

export function ItemDetailPage() {
  const { t, i18n } = useTranslation(['item', 'common']);
  const { id } = useParams<{ id: string }>();
  const user = useUnit($user);
  const query = useItem(id);

  if (query.isPending) {
    return <PageSkeleton />;
  }

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

  const item = query.data;
  const canEdit = user !== null && user.id === item.ownerId && item.status !== 'archived';

  return (
    <Space direction="vertical" size="middle" style={{ width: '100%' }}>
      <Space style={{ width: '100%', justifyContent: 'space-between' }} align="start">
        <Typography.Title level={3} style={{ margin: 0 }}>
          {item.title}
        </Typography.Title>

        {canEdit && (
          <Link to={ROUTES.itemEdit(item.id)}>
            <Button>{t('common:actions.edit')}</Button>
          </Link>
        )}
      </Space>

      <Card>
        <Descriptions column={1} size="middle">
          <Descriptions.Item label={t('item:columns.price')}>
            {formatPrice(item.priceKopeks, i18n.language)}
          </Descriptions.Item>
          <Descriptions.Item label={t('item:columns.status')}>
            <ItemStatusTag status={item.status} />
          </Descriptions.Item>
          <Descriptions.Item label={t('item:columns.createdAt')}>
            {formatDateTime(item.createdAt, i18n.language)}
          </Descriptions.Item>
          <Descriptions.Item label={t('item:form.description')}>
            <Typography.Paragraph style={{ marginBottom: 0, whiteSpace: 'pre-wrap' }}>
              {item.description}
            </Typography.Paragraph>
          </Descriptions.Item>
        </Descriptions>
      </Card>

      <ItemStatusActions item={item} />
    </Space>
  );
}
