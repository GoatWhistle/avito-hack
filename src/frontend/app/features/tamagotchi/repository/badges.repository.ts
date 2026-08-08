import { API_ENDPOINTS } from '#/shared/api/endpoints.api'
import { httpClient } from '#/shared/api/http-client.api'
import type { Badge } from '#/shared/types/badge.type'
import type { AxiosInstance } from 'axios'

export class BadgesRepository {
  constructor(private readonly httpClient: AxiosInstance) {}

  async findMy() {
    const response = await this.httpClient.get<Badge[]>(API_ENDPOINTS.badges.my)

    return response.data
  }
}

export const badgesRepository = new BadgesRepository(httpClient)
