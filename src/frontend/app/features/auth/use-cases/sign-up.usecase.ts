import { authRepository, type AuthRepository } from '#/features/auth/repository'
import type { SignUpRequest } from '#/features/auth/types'

export class SignUpUseCase {
  constructor(private readonly repository: AuthRepository) {}

  async execute(user: SignUpRequest) {
    const response = await this.repository.auth(user, 'sign-up')

    if (response.token) {
      localStorage.setItem('accessToken', response.token)
    }

    return response.user
  }
}

export const signUpUseCase = new SignUpUseCase(authRepository)
