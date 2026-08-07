import {
  leaderboardRepository,
  type LeaderboardRepository,
} from '#/features/leaderboard/repository'
import type { LeaderboardQuery } from '#/features/leaderboard/types'

export class LoadLeaderboardUseCase {
  constructor(private readonly repository: LeaderboardRepository) {}

  execute(query: LeaderboardQuery = {}, signal?: AbortSignal) {
    return this.repository.list(query, signal)
  }
}

export class LoadMyNeighborhoodUseCase {
  constructor(private readonly repository: LeaderboardRepository) {}

  execute(signal?: AbortSignal) {
    return this.repository.list({ around: true }, signal)
  }
}

export const loadLeaderboardUseCase = new LoadLeaderboardUseCase(
  leaderboardRepository,
)
export const loadMyNeighborhoodUseCase = new LoadMyNeighborhoodUseCase(
  leaderboardRepository,
)
