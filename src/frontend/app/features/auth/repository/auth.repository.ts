import type { SignInRequest, SignUpRequest } from '#/features/auth/types'
import { httpClient } from '#/shared/api'
import type { Session } from '#/shared/types'
import type { AxiosInstance } from 'axios'

export type AuthRequest = SignUpRequest | SignInRequest
export type AuthMode = 'sign-in' | 'sign-up'

export class AuthRepository {
  constructor(private readonly httpClient: AxiosInstance) {}

  async auth(request: AuthRequest, url: AuthMode) {
    const { data } = await this.httpClient.post<Session>(url, request)

    return data
  }
}

export const authRepository = new AuthRepository(httpClient)
