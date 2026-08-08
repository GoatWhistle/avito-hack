import { TAMAGOTCHI_QUERY_KEYS } from '#/features/tamagotchi/lib/query-keys'
import { petRepository } from '#/features/tamagotchi/repository/pet.repository'
import { useMutation, useQueryClient } from '@tanstack/react-query'

export const useStrokeMyPet = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: TAMAGOTCHI_QUERY_KEYS.pet.stroke(),
    mutationFn: () => petRepository.stroke(),
    onSuccess: async () => {
      await Promise.all([
        queryClient.invalidateQueries({
          queryKey: TAMAGOTCHI_QUERY_KEYS.pet.my(),
        }),
        queryClient.invalidateQueries({
          queryKey: TAMAGOTCHI_QUERY_KEYS.raccoon.my(),
        }),
      ])
    },
  })
}
