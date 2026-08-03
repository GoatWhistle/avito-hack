import { apiClient } from '@/shared/api';

import { toSession, toUser, type Session, type SessionDto, type User, type UserDto } from '../model/types';

export interface LoginPayload {
  email: string;
  password: string;
}

export interface RegisterPayload {
  email: string;
  password: string;
  displayName: string;
}

export interface UpdateProfilePayload {
  displayName: string;
}

export const sessionApi = {
  login: (payload: LoginPayload): Promise<Session> =>
    apiClient
      .post<SessionDto>('/auth/login', { email: payload.email, password: payload.password })
      .then((response) => toSession(response.data)),

  register: (payload: RegisterPayload): Promise<User> =>
    apiClient
      .post<UserDto>('/auth/register', {
        email: payload.email,
        password: payload.password,
        display_name: payload.displayName,
      })
      .then((response) => toUser(response.data)),

  me: (signal?: AbortSignal): Promise<User> =>
    apiClient
      .get<UserDto>('/users/me', signal === undefined ? {} : { signal })
      .then((response) => toUser(response.data)),

  updateProfile: (payload: UpdateProfilePayload): Promise<User> =>
    apiClient
      .patch<UserDto>('/users/me', { display_name: payload.displayName })
      .then((response) => toUser(response.data)),
};
