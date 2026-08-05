import { Empty } from 'antd';
import { useTranslation } from 'react-i18next';

interface EmptyStateProps {
  description?: string;
}

export function EmptyState({ description }: EmptyStateProps) {
  const { t } = useTranslation('common');

  return <Empty description={description ?? t('app.empty')} />;
}
