import { screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import type * as routerModule from 'react-router'
import { setToken } from '#/api/token-store'
import { SignUpForm } from './sign-up-form'
import { makeUser, renderWithShell } from '#/features/layout/test-utils'
import {
  execute,
  fill,
  navigate,
  submit,
  validPassword,
  validationError,
} from './sign-up-form.harness'

type RouterModule = typeof routerModule

vi.mock('#/features/auth/use-cases', () => ({
  signUpUseCase: { execute: (value: unknown) => execute(value) },
  signInUseCase: { execute: vi.fn() },
}))

vi.mock('react-router', async () => {
  const actual = await vi.importActual<RouterModule>('react-router')

  return { ...actual, useNavigate: () => navigate }
})

beforeEach(() => {
  vi.clearAllMocks()
  setToken(null)
  execute.mockResolvedValue({ token: 'jwt', user: makeUser() })
})

describe('SignUpForm password rules', () => {
  it('blocks a password shorter than eight characters before reaching the server', async () => {
    renderWithShell(<SignUpForm />)

    await fill({
      fullName: 'Иван',
      email: 'ivan@example.com',
      password: '1234567',
    })
    await submit()

    expect(
      await screen.findByText('Пароль должен быть не короче 8 символов'),
    ).toBeInTheDocument()
    expect(execute).not.toHaveBeenCalled()
    expect(navigate).not.toHaveBeenCalled()
  })

  it('blocks a password longer than seventy two characters before reaching the server', async () => {
    renderWithShell(<SignUpForm />)

    await fill({
      fullName: 'Иван',
      email: 'ivan@example.com',
      password: 'a'.repeat(73),
    })
    await submit()

    expect(
      await screen.findByText('Пароль должен быть не длиннее 72 символов'),
    ).toBeInTheDocument()
    expect(execute).not.toHaveBeenCalled()
  }, 20_000)

  it('explains a short password rejected by the server and marks the field', async () => {
    execute.mockRejectedValue(
      validationError('value is shorter than minimum: 8', 'password'),
    )

    renderWithShell(<SignUpForm />)

    await fill({
      fullName: 'Иван',
      email: 'ivan@example.com',
      password: validPassword,
    })
    await submit()

    expect(
      await screen.findByText('Пароль должен быть не короче 8 символов'),
    ).toBeInTheDocument()
    expect(screen.getByLabelText('Пароль')).toHaveAttribute(
      'aria-invalid',
      'true',
    )
    expect(navigate).not.toHaveBeenCalled()
  })

  it('explains a long password rejected by the server', async () => {
    execute.mockRejectedValue(
      validationError('value is longer than maximum: 72', 'password'),
    )

    renderWithShell(<SignUpForm />)

    await fill({
      fullName: 'Иван',
      email: 'ivan@example.com',
      password: validPassword,
    })
    await submit()

    expect(
      await screen.findByText('Пароль должен быть не длиннее 72 символов'),
    ).toBeInTheDocument()
  })

  it('states the password rule before the user types', () => {
    renderWithShell(<SignUpForm />)

    expect(screen.getByText('От 8 до 72 символов')).toBeInTheDocument()
    expect(screen.getByLabelText('Пароль')).toHaveAccessibleDescription(
      'От 8 до 72 символов',
    )
  })
})
