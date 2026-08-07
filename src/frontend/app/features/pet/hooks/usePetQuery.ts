import { useQuery } from '@tanstack/react-query'
import { petRepository } from '#/features/pet/repository'

export const petQueryKey = ['pet'] as const
export const summaryTodayQueryKey = ['pet', 'summary', 'today'] as const

export const usePetQuery = () =>
  useQuery({
    queryKey: petQueryKey,
    queryFn: () => petRepository.state(),
    staleTime: 10_000,
  })

export const useSummaryTodayQuery = (enabled = true) =>
  useQuery({
    queryKey: summaryTodayQueryKey,
    queryFn: () => petRepository.summaryToday(),
    enabled,
    retry: false,
    staleTime: 5 * 60_000,
  })
