import { httpClient } from '#/api'
import type { SignInRequest, SignUpRequest } from '#/features/auth/types'
import {
  toUser,
  type Session,
  type SessionResponse,
  type User,
  type UserResponse,
} from '#/types'
import type { AxiosInstance } from 'axios'

export class AuthRepository {
  constructor(private readonly httpClient: AxiosInstance) {}

  async signUp(user: SignUpRequest): Promise<Session> {
    const response = await this.httpClient.post<SessionResponse>(
      '/auth/register',
      {
        email: user.email,
        password: user.password,
        full_name: user.fullName,
      },
    )

    return {
      token: response.data.token,
      user: toUser(response.data.user),
    }
  }

  async signIn(user: SignInRequest): Promise<Session> {
    const response = await this.httpClient.post<SessionResponse>(
      '/auth/login',
      { email: user.email, password: user.password },
    )

    return {
      token: response.data.token,
      user: toUser(response.data.user),
    }
  }

  async me(): Promise<User> {
    const response = await this.httpClient.get<UserResponse>('/users/me')

    return toUser(response.data)
  }

  async updateProfile(fullName: string): Promise<User> {
    const response = await this.httpClient.patch<UserResponse>('/users/me', {
      full_name: fullName,
    })

    return toUser(response.data)
  }
}

export const authRepository = new AuthRepository(httpClient)
