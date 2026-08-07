import { useQuery } from '@tanstack/react-query'
import { useSession } from '#/features/auth/session'
import { progressRepository } from './progress.repository'

export const progressQueryKey = ['progress', 'profile'] as const

export const useProgress = () => {
  const { isAuthenticated } = useSession()

  return useQuery({
    queryKey: progressQueryKey,
    queryFn: () => progressRepository.profile(),
    enabled: isAuthenticated,
    staleTime: 30_000,
  })
}
