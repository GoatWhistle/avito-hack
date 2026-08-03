import { Button, Result } from 'antd';
import { useTranslation } from 'react-i18next';

import { useApiErrorMessage } from '@/shared/lib/use-api-error';

interface ErrorStateProps {
  error: unknown;
  onRetry?: () => void;
}

export function ErrorState({ error, onRetry }: ErrorStateProps) {
  const { t } = useTranslation('common');
  const getErrorMessage = useApiErrorMessage();

  return (
    <Result
      status="warning"
      title={getErrorMessage(error)}
      extra={
        onRetry === undefined ? null : (
          <Button type="primary" onClick={onRetry}>
            {t('app.retry')}
          </Button>
        )
      }
    />
  );
}
