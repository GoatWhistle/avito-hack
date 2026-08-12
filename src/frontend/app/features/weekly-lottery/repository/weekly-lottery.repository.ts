import type { AxiosInstance } from 'axios'
import { z } from 'zod'
import { httpClient } from '#/api'
import {
  lotteryPrizeCatalogItemSchema,
  lotteryRevealSchema,
  lotteryRunSchema,
  lotteryStateSchema,
  type LotteryPrizeCatalogItem,
  type LotteryReveal,
  type LotteryRun,
  type LotteryState,
} from '#/features/weekly-lottery/types'

const BASE = 'weekly-lottery'

export class WeeklyLotteryRepository {
  constructor(private readonly client: AxiosInstance) {}

  async state(signal?: AbortSignal): Promise<LotteryState> {
    const response = await this.client.get(`${BASE}/state`, { signal })

    return lotteryStateSchema.parse(response.data)
  }

  async prizes(signal?: AbortSignal): Promise<LotteryPrizeCatalogItem[]> {
    const response = await this.client.get(`${BASE}/prizes`, { signal })

    return z.array(lotteryPrizeCatalogItemSchema).parse(response.data)
  }

  async start(): Promise<LotteryRun> {
    const response = await this.client.post(`${BASE}/runs`)

    return lotteryRunSchema.parse(response.data)
  }

  async reveal(runId: string, slot: number): Promise<LotteryReveal> {
    const response = await this.client.post(
      `${BASE}/runs/${runId}/slots/${slot}/reveal`,
    )

    return lotteryRevealSchema.parse(response.data)
  }
}

export const weeklyLotteryRepository = new WeeklyLotteryRepository(httpClient)
