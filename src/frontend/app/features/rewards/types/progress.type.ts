import type { RewardCatalogItem } from './reward.type'

export type RewardGroupKey = 'close' | 'available' | 'far'

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
