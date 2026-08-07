const AUTH = 'auth'
const USER = 'users'
const ITEM = 'items'
const FAVORITE = 'favorites'
const PET = 'pet'

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
  favorite: {
    me: `${FAVORITE}`,
  },
  pet: {
    checkin: `${PET}/checkin`,
    me: `${PET}/me`,
    stroke: `${PET}/actions/stroke`,
  },
} as const
