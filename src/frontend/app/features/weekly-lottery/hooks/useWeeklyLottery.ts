import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { weeklyLotteryRepository } from '#/features/weekly-lottery/repository'
import type {
  LotteryRun,
  LotteryState,
} from '#/features/weekly-lottery/types'

export const weeklyLotteryKeys = {
  all: ['weekly-lottery'] as const,
  state: () => [...weeklyLotteryKeys.all, 'state'] as const,
  prizes: () => [...weeklyLotteryKeys.all, 'prizes'] as const,
}

export const useWeeklyLotteryState = () =>
  useQuery({
    queryKey: weeklyLotteryKeys.state(),
    queryFn: ({ signal }) => weeklyLotteryRepository.state(signal),
  })

export const useWeeklyLotteryPrizes = () =>
  useQuery({
    queryKey: weeklyLotteryKeys.prizes(),
    queryFn: ({ signal }) => weeklyLotteryRepository.prizes(signal),
    staleTime: Number.POSITIVE_INFINITY,
  })

const updateRun = (
  current: LotteryState | undefined,
  run: LotteryRun,
): LotteryState | undefined =>
  current
    ? { ...current, available: run.state === 'active', run }
    : undefined

export const useStartWeeklyLottery = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: () => weeklyLotteryRepository.start(),
    onSuccess: (run) => {
      queryClient.setQueryData<LotteryState>(
        weeklyLotteryKeys.state(),
        (current) => updateRun(current, run),
      )
      void queryClient.invalidateQueries({ queryKey: weeklyLotteryKeys.state() })
    },
  })
}

export const useRevealWeeklyLotterySlot = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ runId, slot }: { runId: string; slot: number }) =>
      weeklyLotteryRepository.reveal(runId, slot),
    onSuccess: ({ run }) => {
      queryClient.setQueryData<LotteryState>(
        weeklyLotteryKeys.state(),
        (current) => updateRun(current, run),
      )
    },
  })
}
