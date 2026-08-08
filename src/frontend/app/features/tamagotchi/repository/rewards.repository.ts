import type { RewardActivateResponse } from '#/features/tamagotchi/types/reward-activate.response'
import type { RewardsResponse } from '#/features/tamagotchi/types/rewards.response'
import { API_ENDPOINTS } from '#/shared/api/endpoints.api'
import { httpClient } from '#/shared/api/http-client.api'
import type { AxiosInstance } from 'axios'

export class RewardsRepository {
  constructor(private readonly httpClient: AxiosInstance) {}

  async findAll() {
    const response = await this.httpClient.get<RewardsResponse>(
      API_ENDPOINTS.rewards.all,
    )

    return response.data
  }

  async activate(id: string) {
    const response = await this.httpClient.post<RewardActivateResponse>(
      API_ENDPOINTS.rewards.activate(id),
    )

    return response.data
  }
}

export const rewardsRepository = new RewardsRepository(httpClient)
