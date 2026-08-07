import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import type * as routerModule from 'react-router'
import { ApiError } from '#/api/api-error'
import { SignUpForm } from './sign-up-form'
import { makeUser, renderWithShell } from '#/features/layout/test-utils'

type RouterModule = typeof routerModule

const execute = vi.fn()
const navigate = vi.fn()

vi.mock('#/features/auth/use-cases', () => ({
  signUpUseCase: { execute: (value: unknown) => execute(value) },
  signInUseCase: { execute: vi.fn() },
}))

vi.mock('react-router', async () => {
  const actual = await vi.importActual<RouterModule>('react-router')

  return { ...actual, useNavigate: () => navigate }
})

const fill = async (values: {
  fullName?: string
  email?: string
  password?: string
}) => {
  if (values.fullName !== undefined) {
    await userEvent.type(
      screen.getByLabelText('Имя и фамилия'),
      values.fullName,
    )
  }
  if (values.email !== undefined) {
    await userEvent.type(
      screen.getByLabelText('Электронная почта'),
      values.email,
    )
  }
  if (values.password !== undefined) {
    await userEvent.type(screen.getByLabelText('Пароль'), values.password)
  }
}

beforeEach(() => {
  vi.clearAllMocks()
  execute.mockResolvedValue({ token: 'jwt', user: makeUser() })
})

describe('SignUpForm', () => {
  it('sends full name, email and password and lands on onboarding', async () => {
    renderWithShell(<SignUpForm />)

    await fill({
      fullName: 'Иван Иванов',
      email: 'ivan@example.com',
      password: 'a',
    })
    await userEvent.click(
      screen.getByRole('button', { name: 'Создать аккаунт' }),
    )

    await waitFor(() => {
      expect(execute).toHaveBeenCalledWith({
        fullName: 'Иван Иванов',
        email: 'ivan@example.com',
        password: 'a',
      })
    })

    await waitFor(() => {
      expect(navigate).toHaveBeenCalledWith('/onboarding', { replace: true })
    })
  })

  it('accepts a one character password', async () => {
    renderWithShell(<SignUpForm />)

    await fill({ fullName: 'Иван', email: 'ivan@example.com', password: 'a' })
    await userEvent.click(
      screen.getByRole('button', { name: 'Создать аккаунт' }),
    )

    await waitFor(() => {
      expect(execute).toHaveBeenCalled()
    })
    expect(screen.queryByText(/не короче/)).not.toBeInTheDocument()
  })

  it('reports a taken email from the backend error envelope', async () => {
    execute.mockRejectedValue(
      new ApiError({
        kind: 'conflict',
        message: 'user already exists',
        status: 409,
      }),
    )

    renderWithShell(<SignUpForm />)

    await fill({ fullName: 'Иван', email: 'taken@example.com', password: 'a' })
    await userEvent.click(
      screen.getByRole('button', { name: 'Создать аккаунт' }),
    )

    expect(await screen.findByText('Эта почта уже занята')).toBeInTheDocument()
    expect(navigate).not.toHaveBeenCalled()
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
})
