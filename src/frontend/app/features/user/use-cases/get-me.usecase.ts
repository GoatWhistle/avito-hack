import {
  userRepository,
  type UserRepository,
} from '#/features/user/repository/user.repository'
import {
  localTokenStorage,
  type LocalTokenStorage,
} from '#/shared/storage/local-token.storage'

export class GetMeUseCase {
  constructor(
    private readonly repository: UserRepository,
    private readonly tokenStorage: LocalTokenStorage,
  ) {}

  execute() {
    if (!this.tokenStorage.get()) return null

    return this.repository.getMe()
  }
}

export const getMeUseCase = new GetMeUseCase(userRepository, localTokenStorage)
