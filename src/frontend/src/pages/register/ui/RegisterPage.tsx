import { Card, Space, Typography } from 'antd';
import { useTranslation } from 'react-i18next';
import { Link, useNavigate } from 'react-router-dom';

import { RegisterForm } from '@/features/auth-register';
import { ROUTES } from '@/shared/config/routes';

const FORM_MAX_WIDTH = 420;

export function RegisterPage() {
  const { t } = useTranslation('auth');
  const navigate = useNavigate();

  return (
    <Space
      direction="vertical"
      size="middle"
      style={{ width: '100%', maxWidth: FORM_MAX_WIDTH, margin: '0 auto' }}
    >
      <Typography.Title level={3}>{t('register.title')}</Typography.Title>

      <Card>
        <RegisterForm
          onSuccess={() => {
            navigate(ROUTES.myItems);
          }}
        />

        <Space size="small">
          <Typography.Text type="secondary">{t('register.hasAccount')}</Typography.Text>
          <Link to={ROUTES.login}>{t('login.submit')}</Link>
        </Space>
      </Card>
    </Space>
  );
}
