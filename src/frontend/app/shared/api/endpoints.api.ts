const AUTH = 'auth'
const USER = 'users'
const ITEMS = 'items'
const FAVORITES = 'favorites'
const PET = 'pet'
const REWARDS = 'rewards'
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
  items: {
    my: `${ITEMS}/mine`,
  },
  favorites: {
    my: `${FAVORITES}`,
  },
  pet: {
    my: `${PET}/me`,
    checkin: `${PET}/checkin`,
    stroke: `${PET}/actions/stroke`,
  },
  rewards: {
    my: `${REWARDS}/my`,
    all: `${REWARDS}`,
    activate: (id: string) => `${REWARDS}/${id}/activate`,
  },
  leaderboard: {
    all: `${LEADERBOARD}`,
  },
  summary: {
    history: `${SUMMARY}/history`,
    today: `${SUMMARY}/today`,
  },
} as const
