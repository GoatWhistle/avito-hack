import type { AxiosInstance } from 'axios'
import { httpClient } from '#/api'
import { questListResponseSchema, type QuestItem } from '#/features/quests/types'

export class QuestsRepository {
  constructor(private readonly httpClient: AxiosInstance) {}

  async today(signal?: AbortSignal): Promise<QuestItem[]> {
    const response = await this.httpClient.get('/quests', { signal })

    return questListResponseSchema.parse(response.data).items
  }

  async claim(): Promise<QuestItem[]> {
    const response = await this.httpClient.post('/quests/claim')

    return questListResponseSchema.parse(response.data).items
  }
}

export const questsRepository = new QuestsRepository(httpClient)
