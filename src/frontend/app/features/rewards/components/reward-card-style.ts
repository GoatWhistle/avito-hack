import type { RewardTrackState } from '#/features/rewards/types'

export type RewardCardState = 'locked' | 'current' | 'completed'

export const rewardCardState = (state: RewardTrackState): RewardCardState => {
  if (state === 'activated' || state === 'granted') return 'completed'
  if (state === 'available') return 'current'

  return 'locked'
}

export const REWARD_CARD_TONE: Record<RewardCardState, string> = {
  locked: 'opacity-75 hover:opacity-100',
  current:
    'bg-achievement-earned ring-2 ring-achievement-earned-border/70 shadow-sm',
  completed: 'bg-achievement-earned ring-achievement-earned-border/30',
}

export const REWARD_MARKER_TONE: Record<RewardCardState, string> = {
  locked: 'bg-muted text-muted-foreground',
  current: 'bg-success text-success-foreground',
  completed: 'bg-success text-success-foreground',
}

export const REWARD_KIND_ICON: Record<string, string> = {
  promo: '％',
  cosmetic: '🎨',
  utility: '🛠',
}

export const REWARD_KIND_TONE: Record<string, string> = {
  promo: 'text-reward-promo',
  cosmetic: 'text-reward-cosmetic',
  utility: 'text-reward-utility',
}

export const REWARD_UNIT_KEY: Record<string, string> = {
  level: 'track.unit.level',
  streak: 'track.unit.streak',
  achievement: 'track.unit.achievement',
}
