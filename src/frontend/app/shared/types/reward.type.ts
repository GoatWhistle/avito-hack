export type RewardConditionType = 'level' | 'streak' | 'achievement'
export type RewardKind = 'promo' | 'utility' | 'cosmetic'
export type RewardStatus = 'granted' | 'activated' | 'expired'

export interface Reward {
  id: string
  title: string
  description: string
  kind: RewardKind
  condition_type: RewardConditionType
  condition_value: number
  unlocked: boolean
  claimed: boolean
  progress_current: number
  progress_target: number
  status?: RewardStatus
}
