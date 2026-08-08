import type { PetCheckinResponse } from '#/features/tamagotchi/types/pet-checkin.response'
import { API_ENDPOINTS } from '#/shared/api/endpoints.api'
import { httpClient } from '#/shared/api/http-client.api'
import type { Pet } from '#/shared/types/pet.type'
import type { AxiosInstance } from 'axios'

export class PetRepository {
  constructor(private readonly httpClient: AxiosInstance) {}

  async getMy() {
    const response = await this.httpClient.get<Pet>(API_ENDPOINTS.pet.my)

    return response.data
  }

  async checkin() {
    const response = await this.httpClient.post<PetCheckinResponse>(
      API_ENDPOINTS.pet.checkin,
    )

    return response.data
  }

  async stroke() {
    const response = await this.httpClient.post<Pet>(API_ENDPOINTS.pet.stroke)

    return response.data
  }
}

export const petRepository = new PetRepository(httpClient)
