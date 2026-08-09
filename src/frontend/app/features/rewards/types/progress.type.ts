import type { MyRewardItem, RewardCatalogItem } from './reward.type'

export type RewardGroupKey = 'close' | 'available' | 'far'

export type RewardTrackState =
  'locked' | 'available' | 'granted' | 'activated' | 'expired'

export type RewardTrackEntry = {
  reward: RewardCatalogItem
  granted: MyRewardItem | null
  state: RewardTrackState
  level: number
  code: string
  current: number
  target: number
}

export type RewardProgress = {
  reward: RewardCatalogItem
  current: number
  target: number
  ratio: number
  percent: number
  remaining: number
  group: RewardGroupKey
}

export type RewardGroup = {
  key: RewardGroupKey
  entries: RewardProgress[]
}
