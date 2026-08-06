import { authRepository, type AuthRepository } from '#/features/auth/repository'
import type { SignInRequest } from '#/features/auth/types'

export class SignInUseCase {
  constructor(private readonly repository: AuthRepository) {}

  async execute(request: SignInRequest) {
    const { token, user } = await this.repository.auth(request, 'sign-in')

    if (token) {
      localStorage.setItem('accessToken', token)
    }

    return user
  }
}

export const signInUseCase = new SignInUseCase(authRepository)
