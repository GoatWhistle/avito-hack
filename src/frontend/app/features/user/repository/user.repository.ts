import { API_ENDPOINTS } from '#/shared/api/endpoints.api'
import { httpClient } from '#/shared/api/http-client.api'
import type { User } from '#/shared/types/user.type'
import type { AxiosInstance } from 'axios'

export class UserRepository {
  constructor(private readonly httpClient: AxiosInstance) {}

  async getMe() {
    const response = await this.httpClient.get<User>(API_ENDPOINTS.user.me)

    return response.data
  }
}

export const userRepository = new UserRepository(httpClient)
