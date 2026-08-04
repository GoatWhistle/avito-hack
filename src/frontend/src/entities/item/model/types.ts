export const ITEM_STATUSES = ['draft', 'moderation', 'published', 'archived'] as const;

export type ItemStatus = (typeof ITEM_STATUSES)[number];

export const STATUS_ACTIONS = ['submit', 'publish', 'archive', 'restore'] as const;

export type StatusAction = (typeof STATUS_ACTIONS)[number];

export interface Item {
  id: string;
  ownerId: string;
  title: string;
  description: string;
  priceKopeks: number;
  status: ItemStatus;
  attributes: Record<string, string>;
  createdAt: string;
  updatedAt: string;
}

export interface ItemListEntry {
  id: string;
  ownerId: string;
  ownerName: string;
  title: string;
  priceKopeks: number;
  status: ItemStatus;
  createdAt: string;
}

export interface ItemDto {
  id: string;
  owner_id: string;
  title: string;
  description: string;
  price: number;
  status: ItemStatus;
  attributes: Record<string, string>;
  created_at: string;
  updated_at: string;
}

export interface ItemListEntryDto {
  id: string;
  owner_id: string;
  owner_name?: string;
  title: string;
  price: number;
  status: ItemStatus;
  created_at: string;
}

export function isItemStatus(value: string): value is ItemStatus {
  return (ITEM_STATUSES as readonly string[]).includes(value);
}
