import {
  useMutation,
  useQuery,
  type UseMutationResult,
  type UseQueryResult,
} from '@tanstack/react-query';

import type { Badge, ClaimedReward, RaccoonProfile } from '../model/types';

import { petApi } from './pet-api';
import { petKeys } from './query-keys';

const PROFILE_STALE_TIME_MS = 15_000;
const BADGES_STALE_TIME_MS = 60_000;

export function usePetProfile(enabled: boolean): UseQueryResult<RaccoonProfile> {
  return useQuery({
    queryKey: petKeys.profile(),
    queryFn: ({ signal }) => petApi.profile(signal),
    staleTime: PROFILE_STALE_TIME_MS,
    enabled,
  });
}

export function useBadges(enabled: boolean): UseQueryResult<Badge[]> {
  return useQuery({
    queryKey: petKeys.badges(),
    queryFn: ({ signal }) => petApi.badges(signal),
    staleTime: BADGES_STALE_TIME_MS,
    enabled,
  });
}

export function useClaimReward(): UseMutationResult<ClaimedReward, unknown, string> {
  return useMutation({
    mutationFn: (rewardId: string) => petApi.claimReward(rewardId),
  });
}
