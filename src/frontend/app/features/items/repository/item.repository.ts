import type { ItemRequest } from '#/features/items/types/item.request'
import type { ItemResponse } from '#/features/items/types/item.response'
import { API_ENDPOINTS } from '#/shared/api/endpoints.api'
import { httpClient } from '#/shared/api/http-client.api'
import type { AxiosInstance } from 'axios'

export class ItemRepository {
  constructor(private readonly httpClient: AxiosInstance) {}

  async findMy(params: ItemRequest) {
    const response = await this.httpClient.get<ItemResponse[]>(
      API_ENDPOINTS.items.my,
      { params },
    )

    return response.data
  }
}

export const itemRepository = new ItemRepository(httpClient)
