export type ItemStatus =
  'draft' | 'moderation' | 'published' | 'sold' | 'archived'

export interface Item {
  id: string
  title: string
  price: number
  status: ItemStatus
}
