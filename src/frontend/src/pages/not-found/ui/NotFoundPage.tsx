import { Button, Result } from 'antd';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';

import { ROUTES } from '@/shared/config/routes';

export function NotFoundPage() {
  const { t } = useTranslation(['errors', 'common']);

  return (
    <Result
      status="404"
      title={t('errors:not_found')}
      extra={
        <Link to={ROUTES.items}>
          <Button type="primary">{t('common:actions.back')}</Button>
        </Link>
      }
    />
  );
}
