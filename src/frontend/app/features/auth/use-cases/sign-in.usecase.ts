import { authRepository, type AuthRepository } from '#/features/auth/repository'
import type { SignInRequest } from '#/features/auth/types'

export class SignInUseCase {
  constructor(private readonly repository: AuthRepository) {}

  async execute(request: SignInRequest) {
    const session = await this.repository.auth(request, 'sign-in')

    localStorage.setItem('accessToken', session.token)

    return session
  }
}

export const signInUseCase = new SignInUseCase(authRepository)
