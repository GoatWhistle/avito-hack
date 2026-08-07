import { useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { favoriteRepository } from '#/features/favorites/repository'
import { favoriteKeys } from '#/features/items/hooks'
import { PAGE_LIMIT } from '#/features/items/lib'
import { getToken } from '#/api'

export const useFavoriteIds = () => {
  const { data } = useQuery({
    queryKey: [...favoriteKeys.all, 'ids'],
    queryFn: () => favoriteRepository.list({ limit: PAGE_LIMIT }),
    enabled: Boolean(getToken()),
    staleTime: 30_000,
  })

  return useMemo(
    () => new Set((data?.items ?? []).map((item) => item.item_id)),
    [data],
  )
}
