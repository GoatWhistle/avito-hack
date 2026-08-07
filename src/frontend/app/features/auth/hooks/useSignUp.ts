import { useMutation } from '@tanstack/react-query'
import { signUpUseCase } from '#/features/auth/use-cases'
import type { SignUpRequest } from '#/features/auth/types'

export const useSignUp = () =>
  useMutation({
    mutationKey: ['signUp'],
    mutationFn: (user: SignUpRequest) => signUpUseCase.execute(user),
  })
