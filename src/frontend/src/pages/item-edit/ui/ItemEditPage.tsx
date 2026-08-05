import { Card, Space, Typography } from 'antd';
import { useTranslation } from 'react-i18next';
import { useNavigate, useParams } from 'react-router-dom';

import { useItem } from '@/entities/item';
import { EditItemForm } from '@/features/item-edit';
import { ROUTES } from '@/shared/config/routes';
import { ErrorState, PageSkeleton } from '@/shared/ui';

const FORM_MAX_WIDTH = 720;

export function ItemEditPage() {
  const { t } = useTranslation('item');
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
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

  return (
    <Space direction="vertical" size="middle" style={{ width: '100%', maxWidth: FORM_MAX_WIDTH }}>
      <Typography.Title level={3}>{t('form.editTitle')}</Typography.Title>

      <Card>
        <EditItemForm
          item={query.data}
          onSuccess={() => {
            navigate(ROUTES.itemDetail(query.data.id));
          }}
        />
      </Card>
    </Space>
  );
}
