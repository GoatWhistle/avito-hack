import {
  itemRepository,
  type ItemRepository,
} from '#/features/items/repository/item.repository'
import type { ItemRequest } from '#/features/items/types/item.request'

export class GetMyItemsUseCase {
  constructor(private readonly repository: ItemRepository) {}

  execute(params: ItemRequest) {
    return this.repository.findMy(params)
  }
}

export const getMyItemsUseCase = new GetMyItemsUseCase(itemRepository)
