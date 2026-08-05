import { App } from 'antd';
import { useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { useClaimReward } from '@/entities/pet';
import type { RewardId } from '@/entities/reward';
import { useApiErrorMessage } from '@/shared/lib/use-api-error';

export interface RewardClaim {
  promocodes: Partial<Record<RewardId, string>>;
  pendingId: RewardId | null;
  claim: (rewardId: RewardId) => void;
}

export function useRewardClaim(): RewardClaim {
  const { t } = useTranslation('reward');
  const { notification } = App.useApp();
  const mutation = useClaimReward();
  const getErrorMessage = useApiErrorMessage();
  const [promocodes, setPromocodes] = useState<Partial<Record<RewardId, string>>>({});
  const [pendingId, setPendingId] = useState<RewardId | null>(null);

  const claim = useCallback(
    (rewardId: RewardId) => {
      setPendingId(rewardId);
      mutation.mutate(rewardId, {
        onSuccess: (result) => {
          setPromocodes((current) => ({ ...current, [rewardId]: result.promocode }));
          notification.success({ message: t('promocode.claimed') });
        },
        onError: (error) => {
          notification.error({ message: getErrorMessage(error) });
        },
        onSettled: () => {
          setPendingId(null);
        },
      });
    },
    [mutation, notification, t, getErrorMessage],
  );

  return { promocodes, pendingId, claim };
}
