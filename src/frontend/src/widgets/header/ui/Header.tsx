import { Button, Layout, Menu } from 'antd';
import { useUnit } from 'effector-react';
import { useTranslation } from 'react-i18next';
import { Link, useLocation, useNavigate } from 'react-router-dom';

import { $isAuthenticated, $user, logoutRequested } from '@/entities/session';
import { ROUTES } from '@/shared/config/routes';
import { LanguageSwitcher, ThemeSwitcher } from '@/shared/ui';

import './header.css';
import { LevelBadge } from './LevelBadge';

export function Header() {
  const { t } = useTranslation('common');
  const navigate = useNavigate();
  const location = useLocation();
  const [isAuthenticated, user] = useUnit([$isAuthenticated, $user]);

  const items = [
    { key: ROUTES.items, label: <Link to={ROUTES.items}>{t('nav.items')}</Link> },
    ...(isAuthenticated
      ? [
          { key: ROUTES.pet, label: <Link to={ROUTES.pet}>{t('nav.pet')}</Link> },
          { key: ROUTES.rewards, label: <Link to={ROUTES.rewards}>{t('nav.rewards')}</Link> },
          {
            key: ROUTES.leaderboard,
            label: <Link to={ROUTES.leaderboard}>{t('nav.leaderboard')}</Link>,
          },
          { key: ROUTES.myItems, label: <Link to={ROUTES.myItems}>{t('nav.myItems')}</Link> },
        ]
      : []),
  ];

  return (
    <Layout.Header className="app-header">
      <Link className="app-header__brand" to={ROUTES.home}>
        {t('app.title')}
      </Link>

      <nav className="app-header__nav" aria-label={t('nav.primary')}>
        <Menu
          mode="horizontal"
          theme="dark"
          selectedKeys={[location.pathname]}
          items={items}
          disabledOverflow
        />
      </nav>

      <div className="app-header__side">
        {isAuthenticated && <LevelBadge />}

        <LanguageSwitcher />
        <ThemeSwitcher />

        {isAuthenticated ? (
          <>
            <Link className="app-header__profile" to={ROUTES.profile}>
              <Button type="text">
                {user?.fullName === undefined || user.fullName === ''
                  ? t('nav.profile')
                  : user.fullName}
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
          </>
        ) : (
          <>
            <Link to={ROUTES.login}>
              <Button>{t('actions.login')}</Button>
            </Link>
            <Link to={ROUTES.register}>
              <Button type="primary">{t('actions.register')}</Button>
            </Link>
          </>
        )}
      </div>
    </Layout.Header>
  );
}
