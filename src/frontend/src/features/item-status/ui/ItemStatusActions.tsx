import { App, Button, Space } from 'antd';
import { useUnit } from 'effector-react';
import { useTranslation } from 'react-i18next';

import { type Item, type ItemStatus, type StatusAction } from '@/entities/item';
import { $user, isModerator } from '@/entities/session';
import { useApiErrorMessage } from '@/shared/lib/use-api-error';

import { useChangeStatus } from '../model/use-change-status';

const OWNER_ACTIONS: Record<ItemStatus, StatusAction[]> = {
  draft: ['submit', 'archive'],
  moderation: [],
  published: ['archive'],
  archived: ['restore'],
};

const MODERATOR_ACTIONS: Record<ItemStatus, StatusAction[]> = {
  draft: ['submit', 'archive'],
  moderation: ['publish', 'archive'],
  published: ['archive'],
  archived: ['restore'],
};

interface ItemStatusActionsProps {
  item: Item;
}

export function ItemStatusActions({ item }: ItemStatusActionsProps) {
  const { t } = useTranslation('item');
  const { notification } = App.useApp();
  const user = useUnit($user);
  const changeStatus = useChangeStatus();
  const getErrorMessage = useApiErrorMessage();

  if (user === null) {
    return null;
  }

  const moderator = isModerator(user);
  const isOwner = item.ownerId === user.id;

  if (!moderator && !isOwner) {
    return null;
  }

  const actions = moderator ? MODERATOR_ACTIONS[item.status] : OWNER_ACTIONS[item.status];

  const run = (action: StatusAction) => {
    changeStatus.mutate(
      { id: item.id, action },
      {
        onSuccess: () => {
          notification.success({ message: t('notifications.statusChanged') });
        },
        onError: (error) => {
          notification.error({ message: getErrorMessage(error) });
        },
      },
    );
  };

  return (
    <Space wrap>
      {actions.map((action) => (
        <Button
          key={action}
          type={action === 'publish' ? 'primary' : 'default'}
          loading={changeStatus.isPending}
          onClick={() => {
            run(action);
          }}
        >
          {t(`actions.${action}`)}
        </Button>
      ))}
    </Space>
  );
}
