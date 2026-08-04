import type { ItemListParams } from './item-api';

export const itemKeys = {
  all: ['items'] as const,
  lists: () => [...itemKeys.all, 'list'] as const,
  list: (params: ItemListParams) => [...itemKeys.lists(), params] as const,
  mineLists: () => [...itemKeys.all, 'mine'] as const,
  mine: (params: ItemListParams) => [...itemKeys.mineLists(), params] as const,
  details: () => [...itemKeys.all, 'detail'] as const,
  detail: (id: string) => [...itemKeys.details(), id] as const,
};
