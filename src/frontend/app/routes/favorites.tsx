import { RequireAuth } from '#/features/auth/components'
import { FavoritesScreen } from '#/features/favorites'

export default function FavoritesRoute() {
  return (
    <RequireAuth>
      <FavoritesScreen />
    </RequireAuth>
  )
}
