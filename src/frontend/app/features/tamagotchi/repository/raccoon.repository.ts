import { API_ENDPOINTS } from '#/shared/api/endpoints.api'
import { httpClient } from '#/shared/api/http-client.api'
import type { Raccoon } from '#/shared/types/raccoon.type'
import type { AxiosInstance } from 'axios'

export class RaccoonRepository {
  constructor(private readonly httpClient: AxiosInstance) {}

  async getMy() {
    const response = await this.httpClient.get<Raccoon>(
      API_ENDPOINTS.raccoon.my,
    )

    return response.data
  }
}

export const raccoonRepository = new RaccoonRepository(httpClient)
