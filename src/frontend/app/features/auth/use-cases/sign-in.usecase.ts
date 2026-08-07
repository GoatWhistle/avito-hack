import {
  authRepository,
  type AuthRepository,
} from '#/features/auth/repository/auth.repository'
import type { SignInRequest } from '#/features/auth/types/sign-in.request.type'
import {
  localTokenStorage,
  type LocalTokenStorage,
} from '#/shared/storage/local-token.storage'

export class SignInUseCase {
  constructor(
    private readonly repository: AuthRepository,
    private readonly tokenStorage: LocalTokenStorage,
  ) {}

  async execute(request: SignInRequest) {
    const session = await this.repository.auth(request, 'login')

    this.tokenStorage.set(session.token)

    return session
  }
}

export const signInUseCase = new SignInUseCase(
  authRepository,
  localTokenStorage,
)
