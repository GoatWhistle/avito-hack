import { TAMAGOTCHI_QUERY_KEYS } from '#/features/tamagotchi/lib/query-keys'
import { petRepository } from '#/features/tamagotchi/repository/pet.repository'
import { useQuery } from '@tanstack/react-query'

export const useGetMyPet = () =>
  useQuery({
    queryKey: TAMAGOTCHI_QUERY_KEYS.pet.my(),
    queryFn: () => petRepository.getMy(),
  })
