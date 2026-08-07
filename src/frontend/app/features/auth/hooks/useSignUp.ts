import type { SignUpRequest } from '#/features/auth/types/sign-up.request'
import { signUpUseCase } from '#/features/auth/use-cases/sign-up.usecase'
import { useMutation } from '@tanstack/react-query'

export const useSignUp = () =>
  useMutation({
    mutationKey: ['signUp'],
    mutationFn: (request: SignUpRequest) => signUpUseCase.execute(request),
  })
