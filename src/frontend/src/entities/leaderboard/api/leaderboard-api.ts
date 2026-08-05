import { useQuery, type UseQueryResult } from '@tanstack/react-query';

import { apiClient } from '@/shared/api';

import { leaderboardDtoSchema, toLeaderboardPage, type LeaderboardPage } from '../model/types';

const STALE_TIME_MS = 30_000;

export const leaderboardKeys = {
  all: ['leaderboard'] as const,
  list: () => [...leaderboardKeys.all, 'list'] as const,
};

export const leaderboardApi = {
  list: (myUserId: string | null, signal?: AbortSignal): Promise<LeaderboardPage> =>
    apiClient
      .get<unknown>('/leaderboard', signal === undefined ? {} : { signal })
      .then((response) => toLeaderboardPage(leaderboardDtoSchema.parse(response.data), myUserId)),
};

export function useLeaderboard(
  myUserId: string | null,
  enabled: boolean,
): UseQueryResult<LeaderboardPage> {
  return useQuery({
    queryKey: leaderboardKeys.list(),
    queryFn: ({ signal }) => leaderboardApi.list(myUserId, signal),
    staleTime: STALE_TIME_MS,
    retry: false,
    enabled,
  });
}
