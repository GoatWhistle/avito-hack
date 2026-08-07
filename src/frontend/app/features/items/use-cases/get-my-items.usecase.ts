import {
  itemRepository,
  type ItemRepository,
} from '#/features/items/repository/item.repository'
import type { ItemsRequest } from '#/features/items/types/items.request'

export class GetMyItemsUseCase {
  constructor(private readonly repository: ItemRepository) {}

  execute(params: ItemsRequest) {
    return this.repository.findMy(params)
  }
}

export const getMyItemsUseCase = new GetMyItemsUseCase(itemRepository)
