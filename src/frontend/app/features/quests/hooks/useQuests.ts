import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {
  claimQuestRewardsUseCase,
  loadDailyQuestsUseCase,
} from '#/features/quests/use-cases'

export const questKeys = {
  today: ['quests', 'today'] as const,
}

export const useDailyQuests = () =>
  useQuery({
    queryKey: questKeys.today,
    queryFn: ({ signal }) => loadDailyQuestsUseCase.execute(signal),
    staleTime: 30_000,
  })

export const useClaimQuestRewards = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: () => claimQuestRewardsUseCase.execute(),
    onSuccess: (quests) => {
      queryClient.setQueryData(questKeys.today, quests)
      void queryClient.invalidateQueries({ queryKey: ['pet'] })
      void queryClient.invalidateQueries({ queryKey: ['rewards'] })
    },
  })
}
