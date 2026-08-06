import { authRepository, type AuthRepository } from '#/features/auth/repository'
import type { SignUpRequest } from '#/features/auth/types'

export class SignUpUseCase {
  constructor(private readonly repository: AuthRepository) {}

  async execute(request: SignUpRequest) {
    const { token, user } = await this.repository.auth(request, 'sign-up')

    if (token) {
      localStorage.setItem('accessToken', token)
    }

    return user
  }
}

export const signUpUseCase = new SignUpUseCase(authRepository)
