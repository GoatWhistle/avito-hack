import { authRepository, type AuthRepository } from '#/features/auth/repository'
import type { SignUpRequest } from '#/features/auth/types'

export class SignUpUseCase {
  constructor(private readonly repository: AuthRepository) {}

  execute(user: SignUpRequest) {
    return this.repository.auth(user, 'sign-up')
  }
}

export const signUpUseCase = new SignUpUseCase(authRepository)
