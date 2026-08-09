import { beforeEach, describe, expect, it, vi } from 'vitest'
import { getToken, setToken } from '#/api/token-store'
import type { AuthRepository } from '#/features/auth/repository'
import { SignUpUseCase } from './sign-up.usecase'

const session = {
  token: 'jwt-from-register',
  user: {
    id: 'user-1',
    email: 'new@example.com',
    fullName: 'Новый Пользователь',
    role: 'user',
    createdAt: '2026-08-07T10:15:30.123456Z',
  },
}

beforeEach(() => {
  setToken(null)
})

describe('SignUpUseCase', () => {
  it('stores the token returned by registration so the user is authenticated at once', async () => {
    const signUp = vi.fn().mockResolvedValue(session)
    const useCase = new SignUpUseCase({
      signUp,
    } as unknown as AuthRepository)

    const result = await useCase.execute({
      email: 'new@example.com',
      password: 'demo1234',
      fullName: 'Новый Пользователь',
    })

    expect(result.token).toBe('jwt-from-register')
    expect(getToken()).toBe('jwt-from-register')
  })

  it('leaves no token behind when registration fails', async () => {
    const signUp = vi.fn().mockRejectedValue(new Error('conflict'))
    const useCase = new SignUpUseCase({
      signUp,
    } as unknown as AuthRepository)

    await expect(
      useCase.execute({
        email: 'taken@example.com',
        password: 'demo1234',
        fullName: 'Занятый Адрес',
      }),
    ).rejects.toThrow('conflict')

    expect(getToken()).toBeNull()
  })
})
