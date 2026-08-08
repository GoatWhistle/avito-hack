const AUTH = 'auth'
const USER = 'users'
const ITEMS = 'items'
const FAVORITES = 'favorites'
const PET = 'pet'
const RACCOON = 'raccoon'
const BADGES = 'badges'
const REWARDS = 'rewards'
const LEADERBOARD = 'leaderboard'
const SUMMARY = 'summary'

export const API_ENDPOINTS = {
  auth: {
    register: `${AUTH}/register`,
    login: `${AUTH}/login`,
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
    my: `${PET}`,
    checkin: '/checkin',
    stroke: `${PET}/actions/stroke`,
    progress: `${PET}/progress`,
  },
  raccoon: {
    my: `${RACCOON}/profile`,
    badges: `/badges`,
  },
  badges: {
    my: `${BADGES}`,
  },
  rewards: {
    my: `${REWARDS}/my`,
    all: `${REWARDS}`,
    claim: `${REWARDS}/claim`,
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
