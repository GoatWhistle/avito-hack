import type { FavoritesRequest } from '#/features/favorites/types/favorites.request'
import type { FavoritesResponse } from '#/features/favorites/types/favorites.response'
import { API_ENDPOINTS } from '#/shared/api/endpoints.api'
import { httpClient } from '#/shared/api/http-client.api'
import type { AxiosInstance } from 'axios'

export class FavoritesRepository {
  constructor(private readonly httpClient: AxiosInstance) {}

  async findMy(params: FavoritesRequest) {
    const response = await this.httpClient.get<FavoritesResponse[]>(
      API_ENDPOINTS.favorites.my,
      { params },
    )

    return response.data
  }
}

export const favoritesRepository = new FavoritesRepository(httpClient)
