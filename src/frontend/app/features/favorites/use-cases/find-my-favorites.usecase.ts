import {
  favoritesRepository,
  type FavoritesRepository,
} from '#/features/favorites/repository/favorites.repository'
import type { FavoritesRequest } from '#/features/favorites/types/favorites.request'

export class FindMyFavoritesUseCase {
  constructor(private readonly repository: FavoritesRepository) {}

  execute(params: FavoritesRequest) {
    return this.repository.findMy(params)
  }
}

export const findMyFavoritesUseCase = new FindMyFavoritesUseCase(
  favoritesRepository,
)
