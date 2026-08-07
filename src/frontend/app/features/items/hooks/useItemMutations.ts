import { useMutation, useQueryClient } from '@tanstack/react-query'
import { itemRepository } from '#/features/items/repository'
import { favoriteKeys, itemKeys } from './query-keys'
import type {
  CreateItemRequest,
  Item,
  ItemStatusAction,
  UpdateItemRequest,
} from '#/features/items/types'

export const useCreateItem = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: ['items', 'create'],
    mutationFn: (payload: CreateItemRequest) => itemRepository.create(payload),
    onSuccess: (item: Item) => {
      queryClient.setQueryData(itemKeys.detail(item.id), item)
      void queryClient.invalidateQueries({ queryKey: itemKeys.all })
    },
  })
}

export const useUpdateItem = (id: string) => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: ['items', 'update', id],
    mutationFn: (payload: UpdateItemRequest) =>
      itemRepository.update(id, payload),
    onSuccess: (item: Item) => {
      queryClient.setQueryData(itemKeys.detail(item.id), item)
      void queryClient.invalidateQueries({ queryKey: itemKeys.lists() })
      void queryClient.invalidateQueries({
        queryKey: [...itemKeys.all, 'mine'],
      })
    },
  })
}

export const useChangeItemStatus = (id: string) => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: ['items', 'status', id],
    mutationFn: (action: ItemStatusAction) =>
      itemRepository.changeStatus(id, action),
    onSuccess: (item: Item) => {
      queryClient.setQueryData(itemKeys.detail(item.id), item)
      void queryClient.invalidateQueries({ queryKey: itemKeys.lists() })
      void queryClient.invalidateQueries({
        queryKey: [...itemKeys.all, 'mine'],
      })
      void queryClient.invalidateQueries({ queryKey: favoriteKeys.all })
    },
  })
}
