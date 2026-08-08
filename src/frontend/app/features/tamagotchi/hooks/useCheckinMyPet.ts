import { TAMAGOTCHI_QUERY_KEYS } from '#/features/tamagotchi/lib/query-keys'
import { petRepository } from '#/features/tamagotchi/repository/pet.repository'
import { useMutation, useQueryClient } from '@tanstack/react-query'

export const useCheckinMyPet = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: TAMAGOTCHI_QUERY_KEYS.pet.checkin(),
    mutationFn: () => petRepository.checkin(),
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: await Promise.all([
          queryClient.invalidateQueries({
            queryKey: TAMAGOTCHI_QUERY_KEYS.pet.my(),
          }),
          queryClient.invalidateQueries({
            queryKey: TAMAGOTCHI_QUERY_KEYS.raccoon.my(),
          }),
          queryClient.invalidateQueries({
            queryKey: TAMAGOTCHI_QUERY_KEYS.badges.my(),
          }),
          queryClient.invalidateQueries({
            queryKey: TAMAGOTCHI_QUERY_KEYS.rewards.my(),
          }),
        ]),
      })
    },
  })
}
