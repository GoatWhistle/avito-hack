import { authRepository, type AuthRepository } from '#/features/auth/repository'
import type { SignInRequest } from '#/features/auth/types'

export class SignInUseCase {
  constructor(private readonly repository: AuthRepository) {}

  execute(user: SignInRequest) {
    return this.repository.auth(user, 'sign-in')
  }
}

export const signInUseCase = new SignInUseCase(authRepository)
