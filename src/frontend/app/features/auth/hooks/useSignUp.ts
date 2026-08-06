import type { SignUpRequest } from '#/features/auth/types'
import { signUpUseCase } from '#/features/auth/use-cases'
import { useMutation } from '@tanstack/react-query'

export const useSignUp = () =>
  useMutation({
    mutationKey: ['signUp'],
    mutationFn: (request: SignUpRequest) => signUpUseCase.execute(request),
  })
