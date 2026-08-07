import type { AxiosInstance } from 'axios'
import { z } from 'zod'
import { httpClient } from '#/api'
import {
  activateRewardResponseSchema,
  badgeItemSchema,
  listResponseSchema,
  myRewardItemSchema,
  rewardCatalogItemSchema,
  type ActivateRewardResponse,
  type BadgeItem,
  type MyRewardItem,
  type RewardCatalogItem,
} from '#/features/rewards/types'

const catalogResponseSchema = listResponseSchema(rewardCatalogItemSchema)
const mineResponseSchema = listResponseSchema(myRewardItemSchema)

const badgesResponseSchema = z.union([
  z.array(badgeItemSchema),
  listResponseSchema(badgeItemSchema).transform((value) => value.items),
  z.null().transform(() => [] as BadgeItem[]),
])

export class RewardsRepository {
  constructor(private readonly httpClient: AxiosInstance) {}

  async catalog(signal?: AbortSignal): Promise<RewardCatalogItem[]> {
    const response = await this.httpClient.get('/rewards', { signal })

    return catalogResponseSchema.parse(response.data).items
  }

  async mine(signal?: AbortSignal): Promise<MyRewardItem[]> {
    const response = await this.httpClient.get('/rewards/my', { signal })

    return mineResponseSchema.parse(response.data).items
  }

  async badges(signal?: AbortSignal): Promise<BadgeItem[]> {
    const response = await this.httpClient.get('/badges/', { signal })

    return badgesResponseSchema.parse(response.data)
  }

  async activate(rewardId: string): Promise<ActivateRewardResponse> {
    const response = await this.httpClient.post(
      `/rewards/${encodeURIComponent(rewardId)}/activate`,
    )

    return activateRewardResponseSchema.parse(response.data)
  }
}

export const rewardsRepository = new RewardsRepository(httpClient)
