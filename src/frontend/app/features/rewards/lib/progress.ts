import type {
  RewardCatalogItem,
  RewardGroup,
  RewardGroupKey,
  RewardProgress,
} from '#/features/rewards/types'

export const CLOSE_THRESHOLD = 0.5

export const groupOrder: RewardGroupKey[] = ['close', 'available', 'far']

const clampRatio = (value: number) => Math.min(1, Math.max(0, value))

export const progressRatio = (current: number, target: number): number => {
  if (!Number.isFinite(target) || target <= 0) return current > 0 ? 1 : 0

  return clampRatio(current / target)
}

export const groupOf = (
  reward: Pick<RewardCatalogItem, 'unlocked' | 'claimed'>,
  ratio: number,
): RewardGroupKey => {
  if (reward.unlocked || reward.claimed) return 'available'

  return ratio >= CLOSE_THRESHOLD ? 'close' : 'far'
}

export const toRewardProgress = (reward: RewardCatalogItem): RewardProgress => {
  const target = Math.max(0, reward.progress_target || reward.condition_value)
  const current = Math.max(
    0,
    Math.min(reward.progress_current, target || Infinity),
  )
  const ratio =
    reward.unlocked || reward.claimed ? 1 : progressRatio(current, target)

  return {
    reward,
    current,
    target,
    ratio,
    percent: Math.round(ratio * 100),
    remaining: Math.max(0, target - current),
    group: groupOf(reward, ratio),
  }
}

const byProximity = (a: RewardProgress, b: RewardProgress) => {
  if (b.ratio !== a.ratio) return b.ratio - a.ratio

  return a.target - b.target
}

export const groupRewards = (rewards: RewardCatalogItem[]): RewardGroup[] => {
  const entries = rewards.map(toRewardProgress)

  return groupOrder
    .map((key) => ({
      key,
      entries: entries.filter((entry) => entry.group === key).sort(byProximity),
    }))
    .filter((group) => group.entries.length > 0)
}
