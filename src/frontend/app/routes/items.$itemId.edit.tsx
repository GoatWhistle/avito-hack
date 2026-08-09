import { useParams } from 'react-router'
import { RequireAuth } from '#/features/auth/components'
import { ItemEditScreen } from '#/features/items'

export default function ItemEditRoute() {
  const { itemId } = useParams()

  if (!itemId) return null

  return (
    <RequireAuth>
      <ItemEditScreen itemId={itemId} />
    </RequireAuth>
  )
}
