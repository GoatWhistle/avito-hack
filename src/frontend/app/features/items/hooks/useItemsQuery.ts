import { useInfiniteQuery } from '@tanstack/react-query'
import { itemRepository } from '#/features/items/repository'
import { PAGE_LIMIT } from '#/features/items/lib'
import { itemKeys, type ItemListFilters } from './query-keys'
import type { ItemListEntry, ListResponse } from '#/features/items/types'

const nextPageParam = (last: ListResponse<ItemListEntry>) =>
  last.next_cursor && last.items.length > 0 ? last.next_cursor : undefined

export const useItemsQuery = (filters: ItemListFilters) =>
  useInfiniteQuery({
    queryKey: itemKeys.list(filters),
    queryFn: ({ pageParam }) =>
      itemRepository.list({
        cursor: pageParam,
        limit: PAGE_LIMIT,
        status: filters.status,
        search: filters.search,
      }),
    initialPageParam: undefined as string | undefined,
    getNextPageParam: nextPageParam,
  })

export const useMyItemsQuery = (filters: ItemListFilters) =>
  useInfiniteQuery({
    queryKey: itemKeys.mine(filters),
    queryFn: ({ pageParam }) =>
      itemRepository.listMine({
        cursor: pageParam,
        limit: PAGE_LIMIT,
        status: filters.status,
      }),
    initialPageParam: undefined as string | undefined,
    getNextPageParam: nextPageParam,
  })

export const flattenPages = <T>(pages: ListResponse<T>[] | undefined): T[] =>
  pages?.flatMap((page) => page.items) ?? []
