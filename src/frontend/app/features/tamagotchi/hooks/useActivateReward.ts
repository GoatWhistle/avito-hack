import { TAMAGOTCHI_QUERY_KEYS } from '#/features/tamagotchi/lib/query-keys'
import { rewardsRepository } from '#/features/tamagotchi/repository/rewards.repository'
import { useMutation, useQueryClient } from '@tanstack/react-query'

export const useActivateReward = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: TAMAGOTCHI_QUERY_KEYS.rewards.activate(),
    mutationFn: (id: string) => rewardsRepository.activate(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: TAMAGOTCHI_QUERY_KEYS.rewards.all(),
      })
    },
  })
}
