import { httpClient } from '#/api'
import type {
  CreateItemRequest,
  Item,
  ItemListEntry,
  ItemPhoto,
  ItemStatusAction,
  ListItemsParams,
  ListResponse,
  UpdateItemRequest,
} from '#/features/items/types'
import type { AxiosInstance } from 'axios'

const BASE = 'items'

export interface UploadPhotoParams {
  itemId: string
  file: File
  onProgress?: (percent: number) => void
  signal?: AbortSignal
}

const toQuery = (params: ListItemsParams) => {
  const query: Record<string, string | number> = {}
  if (params.cursor) query.cursor = params.cursor
  if (params.limit) query.limit = params.limit
  if (params.status) query.status = params.status
  if (params.owner_id) query.owner_id = params.owner_id
  if (params.search) query.search = params.search

  return query
}

export class ItemRepository {
  constructor(private readonly httpClient: AxiosInstance) {}

  async list(params: ListItemsParams) {
    const response = await this.httpClient.get<ListResponse<ItemListEntry>>(
      BASE,
      { params: toQuery(params) },
    )

    return response.data
  }

  async listMine(params: ListItemsParams) {
    const response = await this.httpClient.get<ListResponse<ItemListEntry>>(
      `${BASE}/mine`,
      { params: toQuery(params) },
    )

    return response.data
  }

  async getById(id: string) {
    const response = await this.httpClient.get<Item>(`${BASE}/${id}`)

    return response.data
  }

  async create(payload: CreateItemRequest) {
    const response = await this.httpClient.post<Item>(BASE, payload)

    return response.data
  }

  async update(id: string, payload: UpdateItemRequest) {
    const response = await this.httpClient.patch<Item>(`${BASE}/${id}`, payload)

    return response.data
  }

  async changeStatus(id: string, action: ItemStatusAction) {
    const response = await this.httpClient.post<Item>(`${BASE}/${id}/status`, {
      action,
    })

    return response.data
  }

  async listPhotos(id: string) {
    const response = await this.httpClient.get<ItemPhoto[]>(
      `${BASE}/${id}/photos`,
    )

    return response.data ?? []
  }

  async uploadPhoto({ itemId, file, onProgress, signal }: UploadPhotoParams) {
    const body = new FormData()
    body.append('photo', file)

    const response = await this.httpClient.post<ItemPhoto>(
      `${BASE}/${itemId}/photos`,
      body,
      {
        headers: { 'Content-Type': 'multipart/form-data' },
        signal,
        onUploadProgress: (event) => {
          if (!onProgress) return
          const total = event.total ?? file.size
          if (!total) return
          onProgress(Math.min(100, Math.round((event.loaded / total) * 100)))
        },
      },
    )

    return response.data
  }

  async deletePhoto(itemId: string, photoId: string) {
    await this.httpClient.delete(`${BASE}/${itemId}/photos/${photoId}`)
  }

  async addFavorite(itemId: string) {
    await this.httpClient.post(`${BASE}/${itemId}/favorite`)
  }

  async removeFavorite(itemId: string) {
    await this.httpClient.delete(`${BASE}/${itemId}/favorite`)
  }
}

export const itemRepository = new ItemRepository(httpClient)
