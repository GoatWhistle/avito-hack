import type { ItemStatus, ItemStatusAction } from '#/features/items/types'

const transitions: Record<ItemStatus, ItemStatusAction[]> = {
  draft: ['publish', 'submit', 'archive'],
  moderation: ['publish', 'archive'],
  published: ['sell', 'archive'],
  sold: ['archive'],
  archived: ['restore'],
}

export const availableActions = (status: ItemStatus): ItemStatusAction[] =>
  transitions[status] ?? []

export const isFavoritable = (status: ItemStatus) =>
  status === 'published' || status === 'moderation'

export type StatusTone = 'neutral' | 'info' | 'success' | 'muted'

const tones: Record<ItemStatus, StatusTone> = {
  draft: 'neutral',
  moderation: 'info',
  published: 'success',
  sold: 'info',
  archived: 'muted',
}

export const statusTone = (status: ItemStatus): StatusTone =>
  tones[status] ?? 'neutral'
