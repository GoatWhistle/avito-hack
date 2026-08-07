import type { AxiosInstance } from 'axios'
import { httpClient } from '#/api'
import {
  checkInResultSchema,
  dailySummarySchema,
  petSchema,
  progressSchema,
  type CheckInResult,
  type DailySummary,
  type Pet,
  type Progress,
} from '#/features/pet/types'

export class PetRepository {
  constructor(private readonly httpClient: AxiosInstance) {}

  async state(): Promise<Pet> {
    const response = await this.httpClient.get<unknown>('/pet')

    return petSchema.parse(response.data)
  }

  async stroke(): Promise<Pet> {
    const response = await this.httpClient.post<unknown>('/pet/actions/stroke')

    return petSchema.parse(response.data)
  }

  async checkIn(): Promise<CheckInResult> {
    const response = await this.httpClient.post<unknown>('/checkin')

    return checkInResultSchema.parse(response.data)
  }

  async progress(): Promise<Progress> {
    const response = await this.httpClient.get<unknown>('/progress')

    return progressSchema.parse(response.data)
  }

  async summaryToday(): Promise<DailySummary | null> {
    try {
      const response = await this.httpClient.get<unknown>('/summary/today')

      return dailySummarySchema.parse(response.data)
    } catch (error) {
      if (isNotFound(error)) return null
      throw error
    }
  }
}

const isNotFound = (error: unknown): boolean => {
  if (typeof error !== 'object' || error === null) return false
  const status = error as { status?: number; response?: { status?: number } }
  return status.status === 404 || status.response?.status === 404
}

export const petRepository = new PetRepository(httpClient)
