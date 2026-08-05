import { Card, Space, Typography } from 'antd';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';

import { CreateItemForm } from '@/features/item-create';
import { ROUTES } from '@/shared/config/routes';

const FORM_MAX_WIDTH = 720;

export function ItemCreatePage() {
  const { t } = useTranslation('item');
  const navigate = useNavigate();

  return (
    <Space direction="vertical" size="middle" style={{ width: '100%', maxWidth: FORM_MAX_WIDTH }}>
      <Typography.Title level={3}>{t('form.createTitle')}</Typography.Title>

      <Card>
        <CreateItemForm
          onSuccess={(itemId) => {
            navigate(ROUTES.itemDetail(itemId));
          }}
        />
      </Card>
    </Space>
  );
}
