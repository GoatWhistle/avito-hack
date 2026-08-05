import { HeartOutlined } from '@ant-design/icons';
import { Button, Space, Tooltip } from 'antd';
import { useTranslation } from 'react-i18next';

import type { ConnectionStatus } from '@/entities/pet';

export interface PetActionsProps {
  status: ConnectionStatus;
  onStroke: () => void;
}

export function PetActions({ status, onStroke }: PetActionsProps) {
  const { t } = useTranslation('pet');
  const isOnline = status === 'open';

  return (
    <Space size="middle" wrap>
      <Tooltip title={isOnline ? t('actions.strokeHint') : t('connection.offline')}>
        <Button
          type="primary"
          size="large"
          icon={<HeartOutlined />}
          disabled={!isOnline}
          onClick={onStroke}
        >
          {t('actions.stroke')}
        </Button>
      </Tooltip>
    </Space>
  );
}
