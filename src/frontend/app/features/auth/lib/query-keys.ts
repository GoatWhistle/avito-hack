const AUTH = ['auth'] as const

export const AUTH_QUERY_KEYS = {
  signUp: () => [...AUTH, 'signUp'],
  signIn: () => [...AUTH, 'signIn'],
} as const
