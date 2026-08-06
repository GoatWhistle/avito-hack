import { authRepository, type AuthRepository } from '#/features/auth/repository'
import type { SignUpRequest } from '#/features/auth/types'

export class SignUpUseCase {
  constructor(private readonly repository: AuthRepository) {}

  async execute(request: SignUpRequest) {
    const session = await this.repository.auth(request, 'sign-up')

    localStorage.setItem('accessToken', session.token)

    return session
  }
}

export const signUpUseCase = new SignUpUseCase(authRepository)
