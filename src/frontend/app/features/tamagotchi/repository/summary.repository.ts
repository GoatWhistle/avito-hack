import { API_ENDPOINTS } from '#/shared/api/endpoints.api'
import { httpClient } from '#/shared/api/http-client.api'
import type { DailySummary } from '#/shared/types/summary.type'
import type { AxiosInstance } from 'axios'

export class SummaryRepository {
  constructor(private readonly httpClient: AxiosInstance) {}

  async getToday() {
    const response = await this.httpClient.get<DailySummary>(
      API_ENDPOINTS.summary.today,
    )

    return response.data
  }
}

export const summaryRepository = new SummaryRepository(httpClient)
