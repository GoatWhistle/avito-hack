import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { gameRepository } from '#/features/games/repository'
import { gameKeys } from './query-keys'

export const useGamesQuery = () =>
  useQuery({
    queryKey: gameKeys.list(),
    queryFn: () => gameRepository.list(),
  })

export const useGameStateQuery = (slug: string) =>
  useQuery({
    queryKey: gameKeys.state(slug),
    queryFn: () => gameRepository.state(slug),
    enabled: slug !== '',
  })

export const useClaimGameReward = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: () => gameRepository.claimReward(),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: gameKeys.all })
    },
  })
}
