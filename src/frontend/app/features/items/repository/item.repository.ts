import type { ItemsRequest } from '#/features/items/types/items.request'
import type { ItemsResponse } from '#/features/items/types/items.response'
import { API_ENDPOINTS } from '#/shared/api/endpoints.api'
import { httpClient } from '#/shared/api/http-client.api'
import type { AxiosInstance } from 'axios'

export class ItemRepository {
  constructor(private readonly httpClient: AxiosInstance) {}

  async findMy(params: ItemsRequest) {
    const response = await this.httpClient.get<ItemsResponse[]>(
      API_ENDPOINTS.items.my,
      { params },
    )

    return response.data
  }
}

export const itemRepository = new ItemRepository(httpClient)
