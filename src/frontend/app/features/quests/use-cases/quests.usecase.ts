import {
  questsRepository,
  type QuestsRepository,
} from '#/features/quests/repository'
import { toQuestProgressList } from '#/features/quests/lib/progress'
import type { QuestProgressView } from '#/features/quests/types'

export class LoadDailyQuestsUseCase {
  constructor(private readonly repository: QuestsRepository) {}

  async execute(signal?: AbortSignal): Promise<QuestProgressView[]> {
    return toQuestProgressList(await this.repository.today(signal))
  }
}

export class ClaimQuestRewardsUseCase {
  constructor(private readonly repository: QuestsRepository) {}

  async execute(): Promise<QuestProgressView[]> {
    return toQuestProgressList(await this.repository.claim())
  }
}

export const loadDailyQuestsUseCase = new LoadDailyQuestsUseCase(
  questsRepository,
)
export const claimQuestRewardsUseCase = new ClaimQuestRewardsUseCase(
  questsRepository,
)
