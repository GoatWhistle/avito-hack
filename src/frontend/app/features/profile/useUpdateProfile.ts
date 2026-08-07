import { useMutation, useQueryClient } from '@tanstack/react-query'
import { authRepository } from '#/features/auth/repository'
import { sessionQueryKey } from '#/features/auth/session'

export const useUpdateProfile = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: ['profile', 'update'],
    mutationFn: (fullName: string) => authRepository.updateProfile(fullName),
    onSuccess: (user) => {
      queryClient.setQueryData(sessionQueryKey, user)
    },
  })
}
