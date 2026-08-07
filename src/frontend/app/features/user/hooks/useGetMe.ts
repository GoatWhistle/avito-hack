import { getMeUseCase } from '#/features/user/use-cases/get-me.usecase'
import { useQuery } from '@tanstack/react-query'

export const useGetMe = () =>
  useQuery({
    queryKey: ['user'],
    queryFn: () => getMeUseCase.execute(),
  })
