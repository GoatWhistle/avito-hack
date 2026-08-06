import { authRepository, type AuthRepository } from '#/features/auth/repository'
import type { SignUpRequest } from '#/features/auth/types'
import { localTokenStorage, type LocalTokenStorage } from '#/shared/storage'

export class SignUpUseCase {
  constructor(
    private readonly repository: AuthRepository,
    private readonly tokenStorage: LocalTokenStorage,
  ) {}

  async execute(request: SignUpRequest) {
    const session = await this.repository.auth(request, 'sign-up')

    this.tokenStorage.set(session.token)

    return session
  }
}

export const signUpUseCase = new SignUpUseCase(
  authRepository,
  localTokenStorage,
)
