export const itemStatuses = [
  'draft',
  'moderation',
  'published',
  'sold',
  'archived',
] as const

export type ItemStatus = (typeof itemStatuses)[number]

export const itemStatusActions = [
  'submit',
  'publish',
  'sell',
  'archive',
  'restore',
] as const

export type ItemStatusAction = (typeof itemStatusActions)[number]

export interface Item {
  id: string
  owner_id: string
  title: string
  description: string
  price: number
  status: ItemStatus
  attributes: Record<string, string> | null
  created_at: string
  updated_at: string
}

export interface ItemListEntry {
  id: string
  owner_id: string
  owner_name?: string
  title: string
  price: number
  status: ItemStatus
  created_at: string
}

export interface ItemPhoto {
  id: string
  item_id: string
  url: string
  position: number
  created_at: string
}

export interface FavoriteEntry {
  item_id: string
  owner_id: string
  title: string
  price: number
  status: ItemStatus
  photo_url?: string
  added_at: string
}

export interface ListResponse<T> {
  items: T[]
  next_cursor?: string
}

export interface ListItemsParams {
  cursor?: string
  limit?: number
  status?: ItemStatus | ''
  owner_id?: string
  search?: string
}

export interface CreateItemRequest {
  title: string
  description: string
  price: number
  attributes?: Record<string, string>
}

export interface UpdateItemRequest {
  title?: string
  description?: string
  price?: number
  attributes?: Record<string, string>
}
