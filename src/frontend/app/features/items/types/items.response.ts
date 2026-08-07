import type { Item } from '#/shared/types/item.type'

export interface ItemsResponse {
  items: Item[]
  nextCursor: string
}
