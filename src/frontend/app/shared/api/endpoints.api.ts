const AUTH = 'auth'
const USER = 'users'
const ITEM = 'items'
const FAVORITE = 'favorites'
const PET = 'pet'
const REWARD = 'rewards'
const LEADERBOARD = 'leaderboard'
const SUMMARY = 'summary'

export const API_ENDPOINTS = {
  auth: {
    signUp: `${AUTH}/register`,
    signIn: `${AUTH}/login`,
  },
  user: {
    me: `${USER}/me`,
  },
  item: {
    my: `${ITEM}/mine`,
  },
  favorite: {
    my: `${FAVORITE}`,
  },
  pet: {
    my: `${PET}/me`,
    checkin: `${PET}/checkin`,
    stroke: `${PET}/actions/stroke`,
  },
  rewards: {
    my: `${REWARD}/my`,
    all: `${REWARD}`,
    activate: (id: string) => `${REWARD}/${id}/activate`,
  },
  leaderboard: {
    all: `${LEADERBOARD}`,
  },
  summary: {
    history: `${SUMMARY}/history`,
    today: `${SUMMARY}/today`,
  },
} as const
