import { AUTH_QUERY_KEYS } from '#/features/auth/lib/query-keys'
import type { SignInRequest } from '#/features/auth/types/sign-in.request'
import { signInUseCase } from '#/features/auth/use-cases/sign-in.usecase'
import { useMutation } from '@tanstack/react-query'

export const useSignIn = () =>
  useMutation({
    mutationKey: AUTH_QUERY_KEYS.signIn(),
    mutationFn: (request: SignInRequest) => signInUseCase.execute(request),
  })
