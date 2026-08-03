import { Button, Layout, Menu, Space, Typography } from 'antd';
import { useUnit } from 'effector-react';
import { useTranslation } from 'react-i18next';
import { Link, useLocation, useNavigate } from 'react-router-dom';

import { $isAuthenticated, $user, logoutRequested } from '@/entities/session';
import { ROUTES } from '@/shared/config/routes';
import { layout } from '@/shared/design';
import { LanguageSwitcher, ThemeSwitcher } from '@/shared/ui';

export function Header() {
  const { t } = useTranslation('common');
  const navigate = useNavigate();
  const location = useLocation();
  const [isAuthenticated, user] = useUnit([$isAuthenticated, $user]);

  const items = [
    { key: ROUTES.items, label: <Link to={ROUTES.items}>{t('nav.items')}</Link> },
    ...(isAuthenticated
      ? [{ key: ROUTES.myItems, label: <Link to={ROUTES.myItems}>{t('nav.myItems')}</Link> }]
      : []),
  ];

  return (
    <Layout.Header
      style={{
        display: 'flex',
        alignItems: 'center',
        gap: 'var(--spacing-md)',
        height: layout.headerHeight,
        paddingInline: 'var(--spacing-md)',
      }}
    >
      <Link to={ROUTES.home}>
        <Typography.Text strong style={{ color: 'inherit' }}>
          {t('app.title')}
        </Typography.Text>
      </Link>

      <Menu
        mode="horizontal"
        theme="dark"
        selectedKeys={[location.pathname]}
        items={items}
        style={{ flex: 1, minWidth: 0 }}
      />

      <Space size="small">
        <LanguageSwitcher />
        <ThemeSwitcher />

        {isAuthenticated ? (
          <Space size="small">
            <Link to={ROUTES.profile}>
              <Button type="text" style={{ color: 'inherit' }}>
                {user?.displayName ?? t('nav.profile')}
              </Button>
            </Link>
            <Button
              onClick={() => {
                logoutRequested();
                navigate(ROUTES.items);
              }}
            >
              {t('actions.logout')}
            </Button>
          </Space>
        ) : (
          <Space size="small">
            <Link to={ROUTES.login}>
              <Button>{t('actions.login')}</Button>
            </Link>
            <Link to={ROUTES.register}>
              <Button type="primary">{t('actions.register')}</Button>
            </Link>
          </Space>
        )}
      </Space>
    </Layout.Header>
  );
}
