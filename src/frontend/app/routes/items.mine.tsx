import { RequireAuth } from '#/features/auth/components'
import { MyItemsScreen } from '#/features/items'

export default function MyItemsRoute() {
  return (
    <RequireAuth>
      <MyItemsScreen />
    </RequireAuth>
  )
}
