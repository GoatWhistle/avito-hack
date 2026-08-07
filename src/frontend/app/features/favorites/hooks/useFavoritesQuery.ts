import { useInfiniteQuery } from '@tanstack/react-query'
import { favoriteRepository } from '#/features/favorites/repository'
import { favoriteKeys } from '#/features/items/hooks'
import { PAGE_LIMIT } from '#/features/items/lib'
import type { FavoriteEntry, ListResponse } from '#/features/items/types'

export const useFavoritesQuery = () =>
  useInfiniteQuery({
    queryKey: favoriteKeys.list(),
    queryFn: ({ pageParam }) =>
      favoriteRepository.list({ cursor: pageParam, limit: PAGE_LIMIT }),
    initialPageParam: undefined as string | undefined,
    getNextPageParam: (last: ListResponse<FavoriteEntry>) =>
      last.next_cursor && last.items.length > 0 ? last.next_cursor : undefined,
  })
