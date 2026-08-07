const FAVORITES = ['favorites'] as const

export const FAVORITES_QUERY_KEYS = {
  my: () => [...FAVORITES],
}
