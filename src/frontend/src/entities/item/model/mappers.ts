import type { Item, ItemDto, ItemListEntry, ItemListEntryDto } from './types';

export function toItem(dto: ItemDto): Item {
  return {
    id: dto.id,
    ownerId: dto.owner_id,
    title: dto.title,
    description: dto.description,
    priceKopeks: dto.price,
    status: dto.status,
    attributes: dto.attributes,
    createdAt: dto.created_at,
    updatedAt: dto.updated_at,
  };
}

export function toItemListEntry(dto: ItemListEntryDto): ItemListEntry {
  return {
    id: dto.id,
    ownerId: dto.owner_id,
    ownerName: dto.owner_name ?? '',
    title: dto.title,
    priceKopeks: dto.price,
    status: dto.status,
    createdAt: dto.created_at,
  };
}
