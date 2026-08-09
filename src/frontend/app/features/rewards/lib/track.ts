import type {
  MyRewardItem,
  RewardCatalogItem,
  RewardTrackEntry,
  RewardTrackState,
} from '#/features/rewards/types'
import { isExpired } from './format'

export const trackStateOf = (
  reward: RewardCatalogItem,
  granted: MyRewardItem | null,
): RewardTrackState => {
  if (granted) {
    if (granted.status === 'activated') return 'activated'
    if (granted.status === 'expired' || isExpired(granted.expires_at)) {
      return 'expired'
    }

    return 'granted'
  }

  if (reward.claimed) return 'activated'
  if (reward.unlocked) return 'available'

  return 'locked'
}

export const isActivatable = (state: RewardTrackState) =>
  state === 'available' || state === 'granted'

const levelOf = (reward: RewardCatalogItem) =>
  reward.condition_type === 'level'
    ? reward.condition_value || reward.progress_target
    : reward.progress_target || reward.condition_value

export const buildRewardTrack = (
  rewards: RewardCatalogItem[],
  mine: MyRewardItem[] = [],
): RewardTrackEntry[] => {
  const grantedById = new Map(mine.map((item) => [item.reward_id, item]))

  return rewards
    .map((reward) => {
      const granted = grantedById.get(reward.id) ?? null
      const state = trackStateOf(reward, granted)

      return {
        reward,
        granted,
        state,
        level: levelOf(reward),
        code: granted?.code ? granted.code : '',
        current: Math.max(0, reward.progress_current),
        target: Math.max(0, reward.progress_target || reward.condition_value),
      }
    })
    .sort((a, b) => a.level - b.level)
}

export const trackSummary = (entries: RewardTrackEntry[]) => {
  const claimed = entries.filter(
    (entry) => entry.state === 'activated' || entry.state === 'granted',
  ).length
  const available = entries.filter(
    (entry) => entry.state === 'available',
  ).length

  return { claimed, available, total: entries.length }
}
