import type { SignInRequest } from '#/features/auth/types/sign-in.request.type'
import type { SignUpRequest } from '#/features/auth/types/sign-up.request.type'
import { httpClient } from '#/shared/api/http-client.api'
import type { Session } from '#/shared/types/session.type'
import type { AxiosInstance } from 'axios'

type AuthRequest = SignUpRequest | SignInRequest
type AuthMode = 'login' | 'register'

export class AuthRepository {
  constructor(private readonly httpClient: AxiosInstance) {}

  async auth(request: AuthRequest, mode: AuthMode) {
    const { data } = await this.httpClient.post<Session>(mode, request)

    return data
  }
}

export const authRepository = new AuthRepository(httpClient)
