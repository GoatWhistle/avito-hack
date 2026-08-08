import { TAMAGOTCHI_QUERY_KEYS } from '#/features/tamagotchi/lib/query-keys'
import { summaryRepository } from '#/features/tamagotchi/repository/summary.repository'
import { useQuery } from '@tanstack/react-query'

export const useGetSummaryToday = () =>
  useQuery({
    queryKey: TAMAGOTCHI_QUERY_KEYS.summary.today(),
    queryFn: () => summaryRepository.getToday(),
  })
