import { useQuery } from '@tanstack/react-query'
import { itemRepository } from '#/features/items/repository'
import { itemKeys } from './query-keys'

export const useItemQuery = (id: string | undefined) =>
  useQuery({
    queryKey: itemKeys.detail(id ?? ''),
    queryFn: () => itemRepository.getById(id as string),
    enabled: Boolean(id),
  })

export const useItemPhotosQuery = (id: string | undefined) =>
  useQuery({
    queryKey: itemKeys.photos(id ?? ''),
    queryFn: () => itemRepository.listPhotos(id as string),
    enabled: Boolean(id),
  })
