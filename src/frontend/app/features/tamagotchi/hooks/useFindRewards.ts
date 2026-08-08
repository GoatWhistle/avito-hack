import { TAMAGOTCHI_QUERY_KEYS } from '#/features/tamagotchi/lib/query-keys'
import { rewardsRepository } from '#/features/tamagotchi/repository/rewards.repository'
import { useQuery } from '@tanstack/react-query'

export const useFindRewards = () =>
  useQuery({
    queryKey: TAMAGOTCHI_QUERY_KEYS.rewards.all(),
    queryFn: () => rewardsRepository.findAll(),
  })
