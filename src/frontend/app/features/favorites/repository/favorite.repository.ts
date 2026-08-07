import { httpClient } from '#/api'
import type { FavoriteEntry, ListResponse } from '#/features/items/types'
import type { AxiosInstance } from 'axios'

const BASE = 'favorites'

export interface ListFavoritesParams {
  cursor?: string
  limit?: number
}

export class FavoriteRepository {
  constructor(private readonly httpClient: AxiosInstance) {}

  async list(params: ListFavoritesParams) {
    const query: Record<string, string | number> = {}
    if (params.cursor) query.cursor = params.cursor
    if (params.limit) query.limit = params.limit

    const response = await this.httpClient.get<ListResponse<FavoriteEntry>>(
      BASE,
      { params: query },
    )

    return response.data
  }
}

export const favoriteRepository = new FavoriteRepository(httpClient)
