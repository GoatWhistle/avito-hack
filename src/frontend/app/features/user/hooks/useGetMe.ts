import { USERS_QUERY_KEYS } from '#/features/user/lib/query-keys'
import { getMeUseCase } from '#/features/user/use-cases/get-me.usecase'
import { useQuery } from '@tanstack/react-query'

export const useGetMe = () =>
  useQuery({
    queryKey: USERS_QUERY_KEYS.me(),
    queryFn: () => getMeUseCase.execute(),
  })
