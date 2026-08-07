import { useQuery } from '@tanstack/react-query'
import { httpClient } from '#/api'
import { useSession } from '#/features/auth/session'

interface MineResponse {
  items?: unknown[]
}

export const itemCountQueryKey = ['profile', 'item-count'] as const

export const useItemCount = () => {
  const { isAuthenticated } = useSession()

  return useQuery({
    queryKey: itemCountQueryKey,
    queryFn: async () => {
      const response = await httpClient.get<MineResponse>('/items/mine', {
        params: { limit: 100 },
      })

      return response.data.items?.length ?? 0
    },
    enabled: isAuthenticated,
    staleTime: 60_000,
  })
}
