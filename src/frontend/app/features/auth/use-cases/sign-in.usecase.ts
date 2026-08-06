import { authRepository, type AuthRepository } from '#/features/auth/repository'
import type { SignInRequest } from '#/features/auth/types'

export class SignInUseCase {
  constructor(private readonly repository: AuthRepository) {}

  async execute(user: SignInRequest) {
    const response = await this.repository.auth(user, 'sign-in')

    if (response.token) {
      localStorage.setItem('accessToken', response.token)
    }

    return response.user
  }
}

export const signInUseCase = new SignInUseCase(authRepository)
