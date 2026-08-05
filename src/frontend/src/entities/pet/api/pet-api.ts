import { apiClient } from '@/shared/api';

import { toBadge, toClaimedReward, toRaccoonProfile } from '../model/mappers';
import { badgeListSchema, claimRewardDtoSchema, raccoonProfileDtoSchema } from '../model/schemas';
import type { Badge, ClaimedReward, RaccoonProfile } from '../model/types';

export const petApi = {
  profile: (signal?: AbortSignal): Promise<RaccoonProfile> =>
    apiClient
      .get<unknown>('/raccoon/profile', signal === undefined ? {} : { signal })
      .then((response) => toRaccoonProfile(raccoonProfileDtoSchema.parse(response.data))),

  badges: (signal?: AbortSignal): Promise<Badge[]> =>
    apiClient
      .get<unknown>('/badges/', signal === undefined ? {} : { signal })
      .then((response) => badgeListSchema.parse(response.data).map(toBadge)),

  claimReward: (rewardId: string): Promise<ClaimedReward> =>
    apiClient
      .post<unknown>('/rewards/claim', { reward_id: rewardId })
      .then((response) => toClaimedReward(claimRewardDtoSchema.parse(response.data))),
};
