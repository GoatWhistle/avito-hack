import type { ItemStatus } from '#/shared/types/item.type'

export interface ItemRequest {
  limit: number
  cursor: string
  status: ItemStatus
}
