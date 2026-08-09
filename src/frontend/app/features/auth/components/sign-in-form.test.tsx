import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import type * as routerModule from 'react-router'
import { ApiError } from '#/api/api-error'
import { setToken } from '#/api/token-store'
import { SignInForm } from './sign-in.form'
import { makeUser, renderWithShell } from '#/features/layout/test-utils'

type RouterModule = typeof routerModule

const execute = vi.fn()
const navigate = vi.fn()

vi.mock('#/features/auth/use-cases', () => ({
  signInUseCase: { execute: (value: unknown) => execute(value) },
  signUpUseCase: { execute: vi.fn() },
}))

vi.mock('react-router', async () => {
  const actual = await vi.importActual<RouterModule>('react-router')

  return { ...actual, useNavigate: () => navigate }
})

const validPassword = 'demo1234'

const fill = async (values: { email?: string; password?: string }) => {
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

const submit = () =>
  userEvent.click(screen.getByRole('button', { name: 'Войти' }))

beforeEach(() => {
  vi.clearAllMocks()
  setToken(null)
  execute.mockResolvedValue({ token: 'jwt', user: makeUser() })
})

describe('SignInForm', () => {
  it('signs in and lands on the pet screen', async () => {
    renderWithShell(<SignInForm />)

    await fill({ email: 'ivan@example.com', password: validPassword })
    await submit()

    await waitFor(() => {
      expect(execute).toHaveBeenCalledWith({
        email: 'ivan@example.com',
        password: validPassword,
      })
    })

    await waitFor(() => {
      expect(navigate).toHaveBeenCalledWith('/pet', { replace: true })
    })
  })

  it('explains wrong credentials on a 401 without naming which field failed', async () => {
    execute.mockRejectedValue(
      new ApiError({
        kind: 'unauthorized',
        message: 'authentication required',
        status: 401,
      }),
    )

    renderWithShell(<SignInForm />)

    await fill({ email: 'ivan@example.com', password: 'wrongpassword' })
    await submit()

    expect(
      await screen.findByText('Неверная почта или пароль'),
    ).toBeInTheDocument()
    expect(navigate).not.toHaveBeenCalled()
  })

  it('accepts a legacy short password so existing accounts can sign in', async () => {
    renderWithShell(<SignInForm />)

    await fill({ email: 'ivan@example.com', password: 'old' })
    await submit()

    await waitFor(() => {
      expect(execute).toHaveBeenCalledWith({
        email: 'ivan@example.com',
        password: 'old',
      })
    })
  })

  it('shows an inline email format error after blur', async () => {
    renderWithShell(<SignInForm />)

    await userEvent.type(
      screen.getByLabelText('Электронная почта'),
      'not-an-email',
    )
    await userEvent.tab()

    const message = await screen.findByText('Введите корректный адрес почты')
    expect(message).toBeInTheDocument()
    expect(screen.getByLabelText('Электронная почта')).toHaveAttribute(
      'aria-invalid',
      'true',
    )
  })

  it('moves focus to the first invalid field on submit', async () => {
    renderWithShell(<SignInForm />)

    await submit()

    await waitFor(() => {
      expect(screen.getByLabelText('Электронная почта')).toHaveFocus()
    })
  })

  it('announces server failures in a live region', async () => {
    execute.mockRejectedValue(
      new ApiError({ kind: 'network', message: 'Network Error' }),
    )

    renderWithShell(<SignInForm />)

    await fill({ email: 'ivan@example.com', password: validPassword })
    await submit()

    const alert = await screen.findByRole('alert')
    expect(alert).toHaveTextContent('Нет связи с сервером')
  })

  it('labels every field', () => {
    renderWithShell(<SignInForm />)

    for (const label of ['Электронная почта', 'Пароль']) {
      expect(screen.getByLabelText(label)).toBeInTheDocument()
    }
  })
})
