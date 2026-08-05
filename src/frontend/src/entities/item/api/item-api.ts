import { apiClient, type ApiListResponse } from '@/shared/api';

import { toItem, toItemListEntry } from '../model/mappers';
import type {
  Item,
  ItemDto,
  ItemListEntry,
  ItemListEntryDto,
  ItemStatus,
  StatusAction,
} from '../model/types';

export interface ItemListParams {
  status?: ItemStatus;
  ownerId?: string;
  search?: string;
  cursor?: string;
  limit?: number;
}

export interface ItemListPage {
  items: ItemListEntry[];
  nextCursor: string | null;
}

export interface CreateItemPayload {
  title: string;
  description: string;
  priceKopeks: number;
  attributes?: Record<string, string>;
}

export interface UpdateItemPayload {
  id: string;
  title?: string;
  description?: string;
  priceKopeks?: number;
  attributes?: Record<string, string>;
}

function toQuery(params: ItemListParams): Record<string, string | number> {
  const query: Record<string, string | number> = {};

  if (params.status !== undefined) query.status = params.status;
  if (params.ownerId !== undefined) query.owner_id = params.ownerId;
  if (params.search !== undefined && params.search !== '') query.search = params.search;
  if (params.cursor !== undefined && params.cursor !== '') query.cursor = params.cursor;
  if (params.limit !== undefined) query.limit = params.limit;

  return query;
}

export const itemApi = {
  list: (params: ItemListParams, signal?: AbortSignal): Promise<ItemListPage> =>
    apiClient
      .get<ApiListResponse<ItemListEntryDto>>('/items', {
        params: toQuery(params),
        ...(signal === undefined ? {} : { signal }),
      })
      .then((response) => ({
        items: response.data.items.map(toItemListEntry),
        nextCursor: response.data.next_cursor ?? null,
      })),

  listMine: (params: ItemListParams, signal?: AbortSignal): Promise<ItemListPage> =>
    apiClient
      .get<ApiListResponse<ItemListEntryDto>>('/items/mine', {
        params: toQuery(params),
        ...(signal === undefined ? {} : { signal }),
      })
      .then((response) => ({
        items: response.data.items.map(toItemListEntry),
        nextCursor: response.data.next_cursor ?? null,
      })),

  byId: (id: string, signal?: AbortSignal): Promise<Item> =>
    apiClient
      .get<ItemDto>(`/items/${id}`, signal === undefined ? {} : { signal })
      .then((response) => toItem(response.data)),

  create: (payload: CreateItemPayload): Promise<Item> =>
    apiClient
      .post<ItemDto>('/items', {
        title: payload.title,
        description: payload.description,
        price: payload.priceKopeks,
        attributes: payload.attributes ?? {},
      })
      .then((response) => toItem(response.data)),

  update: ({ id, ...payload }: UpdateItemPayload): Promise<Item> =>
    apiClient
      .patch<ItemDto>(`/items/${id}`, {
        ...(payload.title === undefined ? {} : { title: payload.title }),
        ...(payload.description === undefined ? {} : { description: payload.description }),
        ...(payload.priceKopeks === undefined ? {} : { price: payload.priceKopeks }),
        ...(payload.attributes === undefined ? {} : { attributes: payload.attributes }),
      })
      .then((response) => toItem(response.data)),

  changeStatus: (id: string, action: StatusAction): Promise<Item> =>
    apiClient.post<ItemDto>(`/items/${id}/status`, { action }).then((response) => toItem(response.data)),
};
