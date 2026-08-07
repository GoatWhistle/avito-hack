import type { ItemStatus } from '#/shared/types/item.type'

export interface ItemsRequest {
  limit: number
  cursor: string
  status: ItemStatus
}
