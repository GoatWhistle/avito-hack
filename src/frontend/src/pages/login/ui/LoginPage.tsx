import { Card, Space, Typography } from 'antd';
import { useTranslation } from 'react-i18next';
import { Link, useNavigate } from 'react-router-dom';

import { LoginForm } from '@/features/auth-login';
import { ROUTES } from '@/shared/config/routes';

const FORM_MAX_WIDTH = 420;

export function LoginPage() {
  const { t } = useTranslation('auth');
  const navigate = useNavigate();

  return (
    <Space
      direction="vertical"
      size="middle"
      style={{ width: '100%', maxWidth: FORM_MAX_WIDTH, margin: '0 auto' }}
    >
      <Typography.Title level={3}>{t('login.title')}</Typography.Title>

      <Card>
        <LoginForm
          onSuccess={() => {
            navigate(ROUTES.myItems);
          }}
        />

        <Space size="small">
          <Typography.Text type="secondary">{t('login.noAccount')}</Typography.Text>
          <Link to={ROUTES.register}>{t('register.submit')}</Link>
        </Space>
      </Card>
    </Space>
  );
}
