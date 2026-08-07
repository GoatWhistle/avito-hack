const ITEMS = ['items'] as const

export const ITEMS_QUERY_KEYS = {
  my: () => [...ITEMS, 'my'],
} as const
