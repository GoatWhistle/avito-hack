import { AUTH_QUERY_KEYS } from '#/features/auth/lib/query-keys'
import type { SignUpRequest } from '#/features/auth/types/sign-up.request'
import { signUpUseCase } from '#/features/auth/use-cases/sign-up.usecase'
import { useMutation } from '@tanstack/react-query'

export const useSignUp = () =>
  useMutation({
    mutationKey: AUTH_QUERY_KEYS.signUp,
    mutationFn: (request: SignUpRequest) => signUpUseCase.execute(request),
  })
