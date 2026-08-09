import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi } from 'vitest'
import { ApiError } from '#/api/api-error'

export const execute = vi.fn()
export const navigate = vi.fn()
export const validPassword = 'demo1234'

export const fill = async (values: {
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

export const submit = () =>
  userEvent.click(screen.getByRole('button', { name: 'Зарегистрироваться' }))

export const validationError = (message: string, field: string) =>
  new ApiError({
    kind: 'validation_error',
    message,
    status: 400,
    field,
  })
