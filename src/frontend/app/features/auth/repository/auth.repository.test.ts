import { describe, expect, it, vi } from 'vitest'
import type { AxiosInstance } from 'axios'
import { AuthRepository } from './auth.repository'

const sessionPayload = {
  token: 'jwt-token',
  user: {
    id: 'e6f0d1c2-0000-4000-a000-000000000001',
    email: 'new@demo.avito',
    full_name: 'Новый Пользователь',
    role: 'user',
    created_at: '2026-08-07T10:15:30.123456Z',
  },
}

const clientWith = (post: unknown) => ({ post }) as unknown as AxiosInstance

describe('AuthRepository', () => {
  it('returns a session straight from registration without a second request', async () => {
    const post = vi.fn().mockResolvedValue({ data: sessionPayload })
    const repository = new AuthRepository(clientWith(post))

    const session = await repository.signUp({
      email: 'new@demo.avito',
      password: 'demo1234',
      fullName: 'Новый Пользователь',
    })

    expect(post).toHaveBeenCalledTimes(1)
    expect(post).toHaveBeenCalledWith('/auth/register', {
      email: 'new@demo.avito',
      password: 'demo1234',
      full_name: 'Новый Пользователь',
    })
    expect(session.token).toBe('jwt-token')
    expect(session.user.email).toBe('new@demo.avito')
    expect(session.user.fullName).toBe('Новый Пользователь')
  })

  it('returns a session on sign in', async () => {
    const post = vi.fn().mockResolvedValue({ data: sessionPayload })
    const repository = new AuthRepository(clientWith(post))

    const session = await repository.signIn({
      email: 'new@demo.avito',
      password: 'demo1234',
    })

    expect(post).toHaveBeenCalledWith('/auth/login', {
      email: 'new@demo.avito',
      password: 'demo1234',
    })
    expect(session.token).toBe('jwt-token')
    expect(session.user.id).toBe(sessionPayload.user.id)
  })

  it('propagates registration failures', async () => {
    const failure = new Error('email already registered')
    const repository = new AuthRepository(
      clientWith(vi.fn().mockRejectedValue(failure)),
    )

    await expect(
      repository.signUp({
        email: 'taken@demo.avito',
        password: 'demo1234',
        fullName: 'Занятый Адрес',
      }),
    ).rejects.toThrow(failure)
  })
})
