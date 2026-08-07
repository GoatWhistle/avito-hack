import { useParams } from 'react-router'
import { ItemDetailScreen } from '#/features/items'

export default function ItemDetailRoute() {
  const { itemId } = useParams()

  if (!itemId) return null

  return <ItemDetailScreen itemId={itemId} />
}
