import { useMutation, type UseMutationResult } from '@tanstack/react-query';

import { loginFx, sessionApi, type RegisterPayload, type User } from '@/entities/session';

export function useRegister(): UseMutationResult<User, unknown, RegisterPayload> {
  return useMutation({
    mutationFn: async (payload: RegisterPayload) => {
      const user = await sessionApi.register(payload);
      await loginFx({ email: payload.email, password: payload.password });

      return user;
    },
  });
}
