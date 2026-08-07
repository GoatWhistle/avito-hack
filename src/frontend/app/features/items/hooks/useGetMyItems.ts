import type { ItemRequest } from '#/features/items/types/item.request'
import { getMyItemsUseCase } from '#/features/items/use-cases/get-my-items.usecase'
import { useQuery } from '@tanstack/react-query'

export const useGetMyItems = (params: ItemRequest) =>
  useQuery({
    queryKey: ['myItems'],
    queryFn: () => getMyItemsUseCase.execute(params),
  })
