import { Tag } from 'antd';
import { useTranslation } from 'react-i18next';

import type { ItemStatus } from '../model/types';

const STATUS_COLORS: Record<ItemStatus, string> = {
  draft: 'default',
  moderation: 'processing',
  published: 'success',
  archived: 'warning',
};

interface ItemStatusTagProps {
  status: ItemStatus;
}

export function ItemStatusTag({ status }: ItemStatusTagProps) {
  const { t } = useTranslation('item');

  return <Tag color={STATUS_COLORS[status]}>{t(`status.${status}`)}</Tag>;
}
