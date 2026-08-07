import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ApiError } from '#/api/api-error'
import { ProfileForm } from './ProfileForm'
import { makeUser, renderWithShell } from '#/features/layout/test-utils'

const updateProfile = vi.fn()

vi.mock('#/features/auth/repository', () => ({
  authRepository: {
    updateProfile: (fullName: string) => updateProfile(fullName),
  },
}))

const user = makeUser()

beforeEach(() => {
  vi.clearAllMocks()
  updateProfile.mockResolvedValue(makeUser({ fullName: 'Пётр Петров' }))
})

describe('ProfileForm', () => {
  it('shows the current full name with an edit action', () => {
    renderWithShell(<ProfileForm user={user} />)

    expect(screen.getByText('Иван Иванов')).toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: 'Редактировать' }),
    ).toBeInTheDocument()
  })

  it('saves a new full name', async () => {
    renderWithShell(<ProfileForm user={user} />)

    await userEvent.click(screen.getByRole('button', { name: 'Редактировать' }))

    const input = screen.getByLabelText('Имя и фамилия')
    await userEvent.clear(input)
    await userEvent.type(input, 'Пётр Петров')
    await userEvent.click(screen.getByRole('button', { name: 'Сохранить' }))

    await waitFor(() => {
      expect(updateProfile).toHaveBeenCalledWith('Пётр Петров')
    })
  })

  it('blocks saving an empty name', async () => {
    renderWithShell(<ProfileForm user={user} />)

    await userEvent.click(screen.getByRole('button', { name: 'Редактировать' }))
    await userEvent.clear(screen.getByLabelText('Имя и фамилия'))

    expect(screen.getByRole('button', { name: 'Сохранить' })).toBeDisabled()
    expect(screen.getByText('Обязательное поле')).toBeInTheDocument()
    expect(updateProfile).not.toHaveBeenCalled()
  })

  it('restores the original value on cancel', async () => {
    renderWithShell(<ProfileForm user={user} />)

    await userEvent.click(screen.getByRole('button', { name: 'Редактировать' }))
    await userEvent.type(screen.getByLabelText('Имя и фамилия'), ' лишнее')
    await userEvent.click(screen.getByRole('button', { name: 'Отмена' }))

    expect(screen.getByText('Иван Иванов')).toBeInTheDocument()
    expect(updateProfile).not.toHaveBeenCalled()
  })

  it('surfaces a server error message', async () => {
    updateProfile.mockRejectedValue(
      new ApiError({
        kind: 'validation_error',
        message: 'Имя слишком длинное',
        status: 422,
      }),
    )

    renderWithShell(<ProfileForm user={user} />)

    await userEvent.click(screen.getByRole('button', { name: 'Редактировать' }))
    await userEvent.type(screen.getByLabelText('Имя и фамилия'), 'ы')
    await userEvent.click(screen.getByRole('button', { name: 'Сохранить' }))

    expect(await screen.findByText('Имя слишком длинное')).toBeInTheDocument()
  })
})
