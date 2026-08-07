import { useInfiniteQuery, useQuery } from '@tanstack/react-query'
import { useSession } from '#/features/auth/session'
import {
  loadLeaderboardUseCase,
  loadMyNeighborhoodUseCase,
} from '#/features/leaderboard/use-cases'
import type { LeaderboardPage } from '#/features/leaderboard/types'

export const leaderboardKeys = {
  list: (limit: number) => ['leaderboard', 'list', limit] as const,
  neighborhood: ['leaderboard', 'around-me'] as const,
}

export const useLeaderboard = (limit = 20) => {
  const { isAuthenticated } = useSession()

  return useInfiniteQuery({
    queryKey: leaderboardKeys.list(limit),
    initialPageParam: '',
    queryFn: ({ pageParam, signal }) =>
      loadLeaderboardUseCase.execute({ cursor: pageParam, limit }, signal),
    getNextPageParam: (lastPage: LeaderboardPage) =>
      lastPage.nextCursor ? lastPage.nextCursor : undefined,
    enabled: isAuthenticated,
  })
}

export const useMyNeighborhood = () => {
  const { isAuthenticated } = useSession()

  return useQuery({
    queryKey: leaderboardKeys.neighborhood,
    queryFn: ({ signal }) => loadMyNeighborhoodUseCase.execute(signal),
    enabled: isAuthenticated,
  })
}
