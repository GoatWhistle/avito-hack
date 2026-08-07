const AUTH = 'auth'
const USER = 'users'
const ITEM = 'items'
const FAVORITE = 'favorites'
const PET = 'pet'
const REWARD = 'rewards'
const LEADERBOARD = 'leaderboard'

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
    me: `${PET}/me`,
    checkin: `${PET}/checkin`,
    stroke: `${PET}/actions/stroke`,
  },
  rewards: {
    me: `${REWARD}/my`,
    all: `${REWARD}`,
    activate: (id: string) => `${REWARD}/${id}/activate`,
  },
  leaderboard: {
    all: `${LEADERBOARD}`,
  },
} as const
