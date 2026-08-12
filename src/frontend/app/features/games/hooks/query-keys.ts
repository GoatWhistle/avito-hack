export const gameKeys = {
  all: ['games'] as const,
  list: () => [...gameKeys.all, 'list'] as const,
  state: (slug: string) => [...gameKeys.all, 'state', slug] as const,
}
