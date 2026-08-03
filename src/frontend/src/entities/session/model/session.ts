import { createEffect, createEvent, createStore, sample } from 'effector';

import { tokenStorage } from '@/shared/api';

import { sessionApi, type LoginPayload } from '../api/session-api';


import type { Session, User } from './types';

export const sessionExpired = createEvent();
export const logoutRequested = createEvent();
export const profileUpdated = createEvent<User>();

export const loginFx = createEffect((payload: LoginPayload): Promise<Session> =>
  sessionApi.login(payload),
);

export const fetchMeFx = createEffect((): Promise<User> => sessionApi.me());

const persistTokenFx = createEffect((token: string) => {
  tokenStorage.set(token);
});

const clearTokenFx = createEffect(() => {
  tokenStorage.clear();
});

export const $user = createStore<User | null>(null)
  .on(loginFx.doneData, (_, session) => session.user)
  .on(fetchMeFx.doneData, (_, user) => user)
  .on(profileUpdated, (_, user) => user)
  .reset(sessionExpired, logoutRequested, fetchMeFx.fail);

export const $isAuthenticated = $user.map((user) => user !== null);
export const $isSessionPending = fetchMeFx.pending;
export const $isLoginPending = loginFx.pending;

sample({
  clock: loginFx.doneData,
  fn: (session) => session.token,
  target: persistTokenFx,
});

sample({
  clock: [sessionExpired, logoutRequested],
  target: clearTokenFx,
});

export function bootstrapSession(): void {
  if (tokenStorage.get() !== null) {
    void fetchMeFx();
  }
}
