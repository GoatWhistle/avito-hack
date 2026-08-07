import { FAVORITES_QUERY_KEYS } from '#/features/favorites/lib/query-keys'
import type { FavoritesRequest } from '#/features/favorites/types/favorites.request'
import { findMyFavoritesUseCase } from '#/features/favorites/use-cases/find-my-favorites.usecase'
import { useQuery } from '@tanstack/react-query'

export const useFindMyFavorites = (params: FavoritesRequest) =>
  useQuery({
    queryKey: FAVORITES_QUERY_KEYS.my(),
    queryFn: () => findMyFavoritesUseCase.execute(params),
  })
