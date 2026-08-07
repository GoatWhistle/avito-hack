const AUTH = 'auth'
const USER = 'users'
const ITEM = 'items'

export const API_ENDPOINTS = {
  auth: {
    signUp: `${AUTH}/register`,
    signIn: `${AUTH}/login`,
  },
  user: {
    me: `${USER}/me`,
  },
  item: {
    me: `${ITEM}/mine`,
  },
} as const
