import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {
  activateRewardUseCase,
  loadBadgesUseCase,
  loadMyRewardsUseCase,
  loadRewardCatalogUseCase,
  loadRewardTrackUseCase,
} from '#/features/rewards/use-cases'

export const rewardKeys = {
  catalog: ['rewards', 'catalog'] as const,
  track: ['rewards', 'track'] as const,
  mine: ['rewards', 'mine'] as const,
  badges: ['rewards', 'badges'] as const,
}

export const useRewardCatalog = () =>
  useQuery({
    queryKey: rewardKeys.catalog,
    queryFn: ({ signal }) => loadRewardCatalogUseCase.execute(signal),
  })

export const useRewardTrack = () =>
  useQuery({
    queryKey: rewardKeys.track,
    queryFn: ({ signal }) => loadRewardTrackUseCase.execute(signal),
  })

export const useMyRewards = () =>
  useQuery({
    queryKey: rewardKeys.mine,
    queryFn: ({ signal }) => loadMyRewardsUseCase.execute(signal),
  })

export const useBadges = () =>
  useQuery({
    queryKey: rewardKeys.badges,
    queryFn: ({ signal }) => loadBadgesUseCase.execute(signal),
  })

export const useActivateReward = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: ['rewards', 'activate'],
    mutationFn: (rewardId: string) => activateRewardUseCase.execute(rewardId),
    onSuccess: async () => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: rewardKeys.mine }),
        queryClient.invalidateQueries({ queryKey: rewardKeys.catalog }),
        queryClient.invalidateQueries({ queryKey: rewardKeys.track }),
      ])
    },
  })
}
