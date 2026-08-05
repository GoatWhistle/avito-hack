import {
  keepPreviousData,
  useInfiniteQuery,
  useQuery,
  type UseInfiniteQueryResult,
  type InfiniteData,
  type UseQueryResult,
} from '@tanstack/react-query';

import type { Item } from '../model/types';

import { itemApi, type ItemListPage, type ItemListParams } from './item-api';
import { itemKeys } from './query-keys';


const LIST_STALE_TIME_MS = 30_000;
const DETAIL_STALE_TIME_MS = 60_000;

export function useItems(params: ItemListParams): UseQueryResult<ItemListPage> {
  return useQuery({
    queryKey: itemKeys.list(params),
    queryFn: ({ signal }) => itemApi.list(params, signal),
    staleTime: LIST_STALE_TIME_MS,
    placeholderData: keepPreviousData,
  });
}

export function useMyItems(params: ItemListParams, enabled: boolean): UseQueryResult<ItemListPage> {
  return useQuery({
    queryKey: itemKeys.mine(params),
    queryFn: ({ signal }) => itemApi.listMine(params, signal),
    staleTime: LIST_STALE_TIME_MS,
    placeholderData: keepPreviousData,
    enabled,
  });
}

export function useInfiniteItems(
  params: ItemListParams,
): UseInfiniteQueryResult<InfiniteData<ItemListPage>> {
  return useInfiniteQuery({
    queryKey: itemKeys.list(params),
    queryFn: ({ signal, pageParam }) =>
      itemApi.list({ ...params, ...(pageParam === '' ? {} : { cursor: pageParam }) }, signal),
    initialPageParam: '',
    getNextPageParam: (lastPage) => lastPage.nextCursor ?? undefined,
    staleTime: LIST_STALE_TIME_MS,
  });
}

export function useInfiniteMyItems(
  params: ItemListParams,
  enabled: boolean,
): UseInfiniteQueryResult<InfiniteData<ItemListPage>> {
  return useInfiniteQuery({
    queryKey: itemKeys.mine(params),
    queryFn: ({ signal, pageParam }) =>
      itemApi.listMine({ ...params, ...(pageParam === '' ? {} : { cursor: pageParam }) }, signal),
    initialPageParam: '',
    getNextPageParam: (lastPage) => lastPage.nextCursor ?? undefined,
    staleTime: LIST_STALE_TIME_MS,
    enabled,
  });
}

export function useItem(id: string | undefined): UseQueryResult<Item> {
  return useQuery({
    queryKey: itemKeys.detail(id ?? ''),
    queryFn: ({ signal }) => itemApi.byId(id ?? '', signal),
    staleTime: DETAIL_STALE_TIME_MS,
    enabled: id !== undefined && id !== '',
  });
}
