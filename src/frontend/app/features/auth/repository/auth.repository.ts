import { httpClient } from '#/api'
import type { SignInRequest, SignUpRequest } from '#/features/auth/types'
import type { Session } from '#/types'
import type { AxiosInstance } from 'axios'

export type AuthRequest = SignUpRequest | SignInRequest
export type UrlRequest = 'sign-in' | 'sign-up'

export class AuthRepository {
  constructor(private readonly httpClient: AxiosInstance) {}

  async auth(user: AuthRequest, url: UrlRequest) {
    const response = await this.httpClient.post<Session>(url, user)

    return response.data
  }
}

export const authRepository = new AuthRepository(httpClient)
