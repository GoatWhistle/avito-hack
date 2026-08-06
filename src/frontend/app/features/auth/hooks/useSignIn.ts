import { useMutation } from '@tanstack/react-query'
import type { SignInRequest } from '#/features/auth/types'
import { signInUseCase } from '#/features/auth/use-cases'

export const useSignIn = () =>
  useMutation({
    mutationKey: ['signIn'],
    mutationFn: (user: SignInRequest) => signInUseCase.execute(user),
  })
