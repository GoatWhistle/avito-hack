import { API_ENDPOINTS } from '#/shared/api/endpoints.api'
import { httpClient } from '#/shared/api/http-client.api'
import type { Leaderboard } from '#/shared/types/leaderboard.type'
import type { AxiosInstance } from 'axios'

export class LeaderboardRepository {
  constructor(private readonly httpClient: AxiosInstance) {}

  async get() {
    const response = await this.httpClient.get<Leaderboard>(
      API_ENDPOINTS.leaderboard.all,
    )

    return response.data
  }
}

export const leaderboardRepository = new LeaderboardRepository(httpClient)
