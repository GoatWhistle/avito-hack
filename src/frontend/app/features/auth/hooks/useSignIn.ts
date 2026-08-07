import type { SignInRequest } from '#/features/auth/types/sign-in.request'
import { signInUseCase } from '#/features/auth/use-cases/sign-in.usecase'
import { useMutation } from '@tanstack/react-query'

export const useSignIn = () =>
  useMutation({
    mutationKey: ['signIn'],
    mutationFn: (request: SignInRequest) => signInUseCase.execute(request),
  })
