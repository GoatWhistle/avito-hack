import { httpClient } from '#/api'
import type {
  GameGuessResult,
  GameList,
  GameRewardCode,
  GameRound,
  GameState,
} from '#/features/games/types'
import type { AxiosInstance } from 'axios'

const BASE = 'games'

export class GameRepository {
  constructor(private readonly httpClient: AxiosInstance) {}

  async list() {
    const response = await this.httpClient.get<GameList>(BASE)

    return response.data
  }

  async state(slug: string) {
    const response = await this.httpClient.get<GameState>(`${BASE}/${slug}/state`)

    return response.data
  }

  async startRound(slug: string) {
    const response = await this.httpClient.post<GameRound>(
      `${BASE}/${slug}/rounds`,
    )

    return response.data
  }

  async guess(slug: string, roundId: string, move: unknown) {
    const response = await this.httpClient.post<GameGuessResult>(
      `${BASE}/${slug}/rounds/${roundId}/guess`,
      { move },
    )

    return response.data
  }

  async claimReward() {
    const response = await this.httpClient.post<GameRewardCode>(
      `${BASE}/reward/claim`,
    )

    return response.data
  }
}

export const gameRepository = new GameRepository(httpClient)
