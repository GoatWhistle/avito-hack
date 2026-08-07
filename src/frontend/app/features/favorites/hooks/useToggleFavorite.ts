import { useMutation, useQueryClient } from '@tanstack/react-query'
import { itemRepository } from '#/features/items/repository'
import { favoriteKeys } from '#/features/items/hooks'
import type { FavoriteEntry, ListResponse } from '#/features/items/types'

export interface ToggleFavoriteVariables {
  itemId: string
  isFavorite: boolean
  entry?: Omit<FavoriteEntry, 'added_at'>
}

type FavoritePages = {
  pages: ListResponse<FavoriteEntry>[]
  pageParams: unknown[]
}

const removeFromPages = (data: FavoritePages, itemId: string) => ({
  ...data,
  pages: data.pages.map((page) => ({
    ...page,
    items: page.items.filter((item) => item.item_id !== itemId),
  })),
})

const prependToPages = (data: FavoritePages, entry: FavoriteEntry) => {
  const [first, ...rest] = data.pages
  if (!first) return data

  return {
    ...data,
    pages: [{ ...first, items: [entry, ...first.items] }, ...rest],
  }
}

export const useToggleFavorite = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: ['favorites', 'toggle'],
    mutationFn: ({ itemId, isFavorite }: ToggleFavoriteVariables) =>
      isFavorite
        ? itemRepository.removeFavorite(itemId)
        : itemRepository.addFavorite(itemId),
    onMutate: async (variables) => {
      await queryClient.cancelQueries({ queryKey: favoriteKeys.list() })
      const previous = queryClient.getQueryData<FavoritePages>(
        favoriteKeys.list(),
      )

      if (previous) {
        if (variables.isFavorite) {
          queryClient.setQueryData(
            favoriteKeys.list(),
            removeFromPages(previous, variables.itemId),
          )
        } else if (variables.entry) {
          queryClient.setQueryData(
            favoriteKeys.list(),
            prependToPages(previous, {
              ...variables.entry,
              added_at: new Date().toISOString(),
            }),
          )
        }
      }

      return { previous }
    },
    onError: (_error, _variables, context) => {
      if (context?.previous) {
        queryClient.setQueryData(favoriteKeys.list(), context.previous)
      }
    },
    onSettled: () => {
      void queryClient.invalidateQueries({ queryKey: favoriteKeys.all })
    },
  })
}
