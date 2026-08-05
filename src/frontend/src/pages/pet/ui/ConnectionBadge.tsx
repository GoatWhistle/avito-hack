import { Badge } from 'antd';
import { useTranslation } from 'react-i18next';

import type { ConnectionStatus } from '@/entities/pet';

export interface ConnectionBadgeProps {
  status: ConnectionStatus;
}

const STATUS_TO_BADGE: Record<ConnectionStatus, 'success' | 'processing' | 'default' | 'warning'> =
  {
    open: 'success',
    connecting: 'processing',
    closed: 'warning',
    idle: 'default',
  };

const STATUS_KEY = {
  open: 'connection.open',
  connecting: 'connection.connecting',
  closed: 'connection.closed',
  idle: 'connection.idle',
} as const;

export function ConnectionBadge({ status }: ConnectionBadgeProps) {
  const { t } = useTranslation('pet');

  return (
    <Badge
      status={STATUS_TO_BADGE[status]}
      text={<span className="app-muted">{t(STATUS_KEY[status])}</span>}
      role="status"
    />
  );
}
