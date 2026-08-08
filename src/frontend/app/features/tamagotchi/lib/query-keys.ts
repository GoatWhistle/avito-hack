import type { SummaryHistoryRequest } from '#/features/tamagotchi/types/summary-history.request'

const LEADERBOARD = ['items'] as const
const PET = ['pet'] as const
const RACCOON = ['raccoon'] as const
const REWARDS = ['rewards'] as const
const SUMMARY = ['summary'] as const
const BADGES = ['badges'] as const

export const TAMAGOTCHI_QUERY_KEYS = {
  leaderboard: {
    get: () => [...LEADERBOARD, 'get'],
  },
  pet: {
    my: () => [...PET, 'my'],
    checkin: () => [...PET, 'checkin'],
    stroke: () => [...PET, 'stroke'],
  },
  raccoon: {
    my: () => [...RACCOON, 'my'],
  },
  rewards: {
    my: () => [...REWARDS, 'my'],
    all: () => [...REWARDS, 'all'],
    activate: () => [...REWARDS, 'activate'],
  },
  summary: {
    today: () => [...SUMMARY, 'today'],
    history: (params: SummaryHistoryRequest) => [...SUMMARY, 'history', params],
  },
  badges: {
    my: () => [...BADGES, 'my'],
  },
} as const
