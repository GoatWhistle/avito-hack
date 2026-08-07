import { setToken } from '#/api'
import { authRepository, type AuthRepository } from '#/features/auth/repository'
import type { SignInRequest } from '#/features/auth/types'
import type { Session } from '#/types'

export class SignInUseCase {
  constructor(private readonly repository: AuthRepository) {}

  async execute(user: SignInRequest): Promise<Session> {
    const session = await this.repository.signIn(user)
    setToken(session.token)

    return session
  }
}

export const signInUseCase = new SignInUseCase(authRepository)
