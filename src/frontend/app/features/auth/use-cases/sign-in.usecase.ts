import { authRepository, type AuthRepository } from '#/features/auth/repository'
import type { SignInRequest } from '#/features/auth/types'
import { localTokenStorage, type LocalTokenStorage } from '#/shared/storage'

export class SignInUseCase {
  constructor(
    private readonly repository: AuthRepository,
    private readonly tokenStorage: LocalTokenStorage,
  ) {}

  async execute(request: SignInRequest) {
    const session = await this.repository.auth(request, 'sign-in')

    this.tokenStorage.set(session.token)

    return session
  }
}

export const signInUseCase = new SignInUseCase(
  authRepository,
  localTokenStorage,
)
