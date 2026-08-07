const USER = ['user'] as const

export const USERS_QUERY_KEYS = {
  me: () => [...USER, 'me'],
} as const
