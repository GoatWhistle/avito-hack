import {
  rewardsRepository,
  type RewardsRepository,
} from '#/features/rewards/repository'
import { groupRewards } from '#/features/rewards/lib/progress'
import type { RewardGroup } from '#/features/rewards/types'

export class LoadRewardCatalogUseCase {
  constructor(private readonly repository: RewardsRepository) {}

  async execute(signal?: AbortSignal): Promise<RewardGroup[]> {
    const rewards = await this.repository.catalog(signal)

    return groupRewards(rewards)
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
export const loadMyRewardsUseCase = new LoadMyRewardsUseCase(rewardsRepository)
export const loadBadgesUseCase = new LoadBadgesUseCase(rewardsRepository)
export const activateRewardUseCase = new ActivateRewardUseCase(
  rewardsRepository,
)
