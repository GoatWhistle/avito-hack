import { RequireAuth } from '#/features/auth/components'
import { ItemCreateScreen } from '#/features/items'

export default function ItemCreateRoute() {
  return (
    <RequireAuth>
      <ItemCreateScreen />
    </RequireAuth>
  )
}
