import { Button, Empty } from 'antd';
import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';

import './shared-ui.css';

interface EmptyStateProps {
  description?: string;
  hint?: string;
  action?: ReactNode;
  actionLabel?: string;
  onAction?: () => void;
}

export function EmptyState({ description, hint, action, actionLabel, onAction }: EmptyStateProps) {
  const { t } = useTranslation('common');

  return (
    <Empty
      image={Empty.PRESENTED_IMAGE_SIMPLE}
      description={
        <span className="empty-state">
          <span className="empty-state__title">{description ?? t('app.empty')}</span>
          {hint !== undefined && <span className="empty-state__hint">{hint}</span>}
        </span>
      }
    >
      {action ??
        (actionLabel !== undefined && (
          <Button type="primary" onClick={onAction}>
            {actionLabel}
          </Button>
        ))}
    </Empty>
  );
}
