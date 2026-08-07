import type { Item } from '#/shared/types/item.type'

export interface ItemResponse {
  item: Item[]
  nextCursor: string
}
