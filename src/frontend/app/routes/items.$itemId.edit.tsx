import { useParams } from 'react-router'
import { ItemEditScreen } from '#/features/items'

export default function ItemEditRoute() {
  const { itemId } = useParams()

  if (!itemId) return null

  return <ItemEditScreen itemId={itemId} />
}
