import { authRepository, type AuthRepository } from '#/features/auth/repository'
import type { SignUpRequest } from '#/features/auth/types'
import type { Session } from '#/types'

export class SignUpUseCase {
  constructor(private readonly repository: AuthRepository) {}

  async execute(user: SignUpRequest): Promise<Session> {
    return this.repository.signUp(user)
  }
}

export const signUpUseCase = new SignUpUseCase(authRepository)
