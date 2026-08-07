import type { Item } from '#/shared/types/item.type'

export interface ItemResponse {
  items: Item[]
  nextCursor: string
}
