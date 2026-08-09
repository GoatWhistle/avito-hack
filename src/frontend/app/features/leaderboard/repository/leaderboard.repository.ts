import type { AxiosInstance } from 'axios'
import { httpClient } from '#/api'
import {
  leaderboardResponseSchema,
  type LeaderboardPage,
  type LeaderboardQuery,
} from '#/features/leaderboard/types'

export class LeaderboardRepository {
  constructor(private readonly httpClient: AxiosInstance) {}

  async list(
    query: LeaderboardQuery = {},
    signal?: AbortSignal,
  ): Promise<LeaderboardPage> {
    const params: Record<string, string | number | boolean> = {
      with_my_rank: true,
    }

    if (query.cursor) params.cursor = query.cursor
    if (query.limit) params.limit = query.limit
    if (query.around) params.around = 'me'

    const response = await this.httpClient.get('/leaderboard', {
      params,
      signal,
    })

    const parsed = leaderboardResponseSchema.parse(response.data)

    return {
      items: parsed.items,
      myRank: parsed.my_rank ?? null,
      nextCursor: parsed.next_cursor,
    }
  }
}

export const leaderboardRepository = new LeaderboardRepository(httpClient)
