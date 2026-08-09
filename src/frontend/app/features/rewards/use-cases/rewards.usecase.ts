import {
  rewardsRepository,
  type RewardsRepository,
} from '#/features/rewards/repository'
import { groupRewards } from '#/features/rewards/lib/progress'
import { buildRewardTrack } from '#/features/rewards/lib/track'
import type { RewardGroup, RewardTrackEntry } from '#/features/rewards/types'

export class LoadRewardCatalogUseCase {
  constructor(private readonly repository: RewardsRepository) {}

  async execute(signal?: AbortSignal): Promise<RewardGroup[]> {
    const rewards = await this.repository.catalog(signal)

    return groupRewards(rewards)
  }
}

export class LoadRewardTrackUseCase {
  constructor(private readonly repository: RewardsRepository) {}

  async execute(signal?: AbortSignal): Promise<RewardTrackEntry[]> {
    const [catalog, mine] = await Promise.all([
      this.repository.catalog(signal),
      this.repository.mine(signal).catch(() => []),
    ])

    return buildRewardTrack(catalog, mine)
  }
}

export class LoadMyRewardsUseCase {
  constructor(private readonly repository: RewardsRepository) {}

  execute(signal?: AbortSignal) {
    return this.repository.mine(signal)
  }
}

export class LoadBadgesUseCase {
  constructor(private readonly repository: RewardsRepository) {}

  execute(signal?: AbortSignal) {
    return this.repository.badges(signal)
  }
}

export class ActivateRewardUseCase {
  constructor(private readonly repository: RewardsRepository) {}

  execute(rewardId: string) {
    return this.repository.activate(rewardId)
  }
}

export const loadRewardCatalogUseCase = new LoadRewardCatalogUseCase(
  rewardsRepository,
)
export const loadRewardTrackUseCase = new LoadRewardTrackUseCase(
  rewardsRepository,
)
export const loadMyRewardsUseCase = new LoadMyRewardsUseCase(rewardsRepository)
export const loadBadgesUseCase = new LoadBadgesUseCase(rewardsRepository)
export const activateRewardUseCase = new ActivateRewardUseCase(
  rewardsRepository,
)
