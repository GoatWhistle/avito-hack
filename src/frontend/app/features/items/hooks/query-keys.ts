import type { ListItemsParams } from '#/features/items/types'

export type ItemListFilters = Pick<
  ListItemsParams,
  'status' | 'search' | 'category' | 'condition' | 'sort'
>

export const itemKeys = {
  all: ['items'] as const,
  lists: () => [...itemKeys.all, 'list'] as const,
  list: (filters: ItemListFilters) =>
    [
      ...itemKeys.lists(),
      filters.status ?? '',
      filters.search ?? '',
      filters.category ?? '',
      filters.condition ?? '',
      filters.sort ?? 'newest',
    ] as const,
  mine: (filters: ItemListFilters) =>
    [
      ...itemKeys.all,
      'mine',
      filters.status ?? '',
      filters.category ?? '',
      filters.condition ?? '',
      filters.sort ?? 'newest',
    ] as const,
  details: () => [...itemKeys.all, 'detail'] as const,
  detail: (id: string) => [...itemKeys.details(), id] as const,
  photos: (id: string) => [...itemKeys.all, 'photos', id] as const,
}

export const favoriteKeys = {
  all: ['favorites'] as const,
  list: () => [...favoriteKeys.all, 'list'] as const,
}
