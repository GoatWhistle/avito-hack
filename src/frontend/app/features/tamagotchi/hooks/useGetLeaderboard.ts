import { TAMAGOTCHI_QUERY_KEYS } from '#/features/tamagotchi/lib/query-keys'
import { leaderboardRepository } from '#/features/tamagotchi/repository/leaderboard.repository'
import { useQuery } from '@tanstack/react-query'

export const useGetLeaderboard = () =>
  useQuery({
    queryKey: TAMAGOTCHI_QUERY_KEYS.leaderboard.get(),
    queryFn: () => leaderboardRepository.get(),
  })
