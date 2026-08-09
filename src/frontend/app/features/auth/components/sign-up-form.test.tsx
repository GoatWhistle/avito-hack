import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import type * as routerModule from 'react-router'
import { ApiError } from '#/api/api-error'
import { setToken } from '#/api/token-store'
import { SignUpForm } from './sign-up-form'
import { makeUser, renderWithShell } from '#/features/layout/test-utils'
import {
  execute,
  fill,
  navigate,
  submit,
  validPassword,
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

describe('SignUpForm', () => {
  it('sends full name, email and password and lands on onboarding', async () => {
    renderWithShell(<SignUpForm />)

    await fill({
      fullName: 'Иван Иванов',
      email: 'ivan@example.com',
      password: validPassword,
    })
    await submit()

    await waitFor(() => {
      expect(execute).toHaveBeenCalledWith({
        fullName: 'Иван Иванов',
        email: 'ivan@example.com',
        password: validPassword,
      })
    })

    await waitFor(() => {
      expect(navigate).toHaveBeenCalledWith('/onboarding', { replace: true })
    })
  })

  it('reports a taken email from the backend error envelope', async () => {
    execute.mockRejectedValue(
      new ApiError({
        kind: 'conflict',
        message: 'state conflict',
        status: 409,
      }),
    )

    renderWithShell(<SignUpForm />)

    await fill({
      fullName: 'Иван',
      email: 'taken@example.com',
      password: validPassword,
    })
    await submit()

    expect(await screen.findByText('Эта почта уже занята')).toBeInTheDocument()
    expect(navigate).not.toHaveBeenCalled()
  })

  it('focuses the email field when the server reports a duplicate email', async () => {
    execute.mockRejectedValue(
      new ApiError({
        kind: 'conflict',
        message: 'state conflict',
        status: 409,
      }),
    )

    renderWithShell(<SignUpForm />)

    await fill({
      fullName: 'Иван',
      email: 'taken@example.com',
      password: validPassword,
    })
    await submit()

    await waitFor(() => {
      expect(screen.getByLabelText('Электронная почта')).toHaveFocus()
    })
  })

  it('shows an inline email format error after blur', async () => {
    renderWithShell(<SignUpForm />)

    await userEvent.type(
      screen.getByLabelText('Электронная почта'),
      'not-an-email',
    )
    await userEvent.tab()

    expect(
      await screen.findByText('Введите корректный адрес почты'),
    ).toBeInTheDocument()
  })

  it('moves focus to the first invalid field on submit', async () => {
    renderWithShell(<SignUpForm />)

    await submit()

    await waitFor(() => {
      expect(screen.getByLabelText('Имя и фамилия')).toHaveFocus()
    })
  })

  it('labels every field and keeps them valid until touched', () => {
    renderWithShell(<SignUpForm />)

    for (const label of ['Имя и фамилия', 'Электронная почта', 'Пароль']) {
      const input = screen.getByLabelText(label)
      expect(input).toBeInTheDocument()
      expect(input).not.toHaveAttribute('aria-invalid')
    }
  })
})
