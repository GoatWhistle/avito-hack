import { TAMAGOTCHI_QUERY_KEYS } from '#/features/tamagotchi/lib/query-keys'
import { raccoonRepository } from '#/features/tamagotchi/repository/raccoon.repository'
import { useQuery } from '@tanstack/react-query'

export const useGetMyRaccoon = () =>
  useQuery({
    queryKey: TAMAGOTCHI_QUERY_KEYS.raccoon.my(),
    queryFn: () => raccoonRepository.getMy(),
  })
