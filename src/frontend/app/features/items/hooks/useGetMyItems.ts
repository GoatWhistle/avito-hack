import { ITEMS_QUERY_KEYS } from '#/features/items/lib/query-keys'
import type { ItemsRequest } from '#/features/items/types/items.request'
import { getMyItemsUseCase } from '#/features/items/use-cases/get-my-items.usecase'
import { useQuery } from '@tanstack/react-query'

export const useGetMyItems = (params: ItemsRequest) =>
  useQuery({
    queryKey: ITEMS_QUERY_KEYS.my(),
    queryFn: () => getMyItemsUseCase.execute(params),
  })
