export const itemStatuses = [
  'draft',
  'moderation',
  'published',
  'sold',
  'archived',
] as const

export type ItemStatus = (typeof itemStatuses)[number]

export const itemCategories = [
  'electronics',
  'appliances',
  'furniture',
  'clothes',
  'kids',
  'sport',
  'hobby',
  'music',
  'books',
  'auto',
  'realty',
  'beauty',
  'animals',
  'other',
] as const

export type ItemCategory = (typeof itemCategories)[number]

export const itemConditions = ['new', 'used'] as const

export type ItemCondition = (typeof itemConditions)[number]

export const itemSorts = ['newest', 'price_asc', 'price_desc'] as const

export type ItemSort = (typeof itemSorts)[number]

export const itemStatusActions = [
  'submit',
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
  moderation_reason?: string
  is_seed: boolean
  ai_verified: boolean
}

export interface ItemListEntry {
  id: string
  owner_id: string
  owner_name?: string
  title: string
  price: number
  status: ItemStatus
  category?: string
  condition?: string
  created_at: string
  is_seed: boolean
  ai_verified: boolean
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
  category?: ItemCategory | ''
  condition?: ItemCondition | ''
  sort?: ItemSort
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
