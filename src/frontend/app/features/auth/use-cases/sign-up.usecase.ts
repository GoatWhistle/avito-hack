import {
  authRepository,
  type AuthRepository,
} from '#/features/auth/repository/auth.repository'
import type { SignUpRequest } from '#/features/auth/types/sign-up.request'
import {
  localTokenStorage,
  type LocalTokenStorage,
} from '#/shared/storage/local-token.storage'

export class SignUpUseCase {
  constructor(
    private readonly repository: AuthRepository,
    private readonly tokenStorage: LocalTokenStorage,
  ) {}

  async execute(request: SignUpRequest) {
    const session = await this.repository.auth(request, 'register')

    this.tokenStorage.set(session.token)

    return session
  }
}

export const signUpUseCase = new SignUpUseCase(
  authRepository,
  localTokenStorage,
)
