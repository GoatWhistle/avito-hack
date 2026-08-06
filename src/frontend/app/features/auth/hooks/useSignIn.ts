import type { SignInRequest } from '#/features/auth/types'
import { signInUseCase } from '#/features/auth/use-cases'
import { useMutation } from '@tanstack/react-query'

export const useSignIn = () =>
  useMutation({
    mutationKey: ['signIn'],
    mutationFn: (request: SignInRequest) => signInUseCase.execute(request),
  })
