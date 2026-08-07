import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { PetActions } from './PetActions'
import {
  PetScreenEmpty,
  PetScreenError,
  PetScreenSkeleton,
} from './PetScreenStates'
import { renderWithProviders } from './test-utils'

describe('PetScreenSkeleton', () => {
  it('announces a busy loading state', () => {
    renderWithProviders(<PetScreenSkeleton />)

    const status = screen.getByRole('status', { name: 'Загружаем питомца' })
    expect(status).toHaveAttribute('aria-busy', 'true')
  })
})

describe('PetScreenError', () => {
  it('shows the failure message and allows a retry', async () => {
    const onRetry = vi.fn()
    renderWithProviders(
      <PetScreenError
        message="Сервис недоступен"
        onRetry={onRetry}
        isRetrying={false}
      />,
    )

    expect(screen.getByRole('alert')).toHaveTextContent(
      'Не удалось загрузить питомца',
    )
    expect(screen.getByText('Сервис недоступен')).toBeInTheDocument()

    await userEvent.click(
      screen.getByRole('button', { name: 'Попробовать снова' }),
    )
    expect(onRetry).toHaveBeenCalledTimes(1)
  })

  it('disables the retry button while retrying', async () => {
    const onRetry = vi.fn()
    renderWithProviders(
      <PetScreenError message="Ошибка" onRetry={onRetry} isRetrying />,
    )

    const button = screen.getByRole('button', { name: 'Загружаем…' })
    expect(button).toBeDisabled()

    await userEvent.click(button)
    expect(onRetry).not.toHaveBeenCalled()
  })
})

describe('PetScreenEmpty', () => {
  it('invites the first check in', async () => {
    const onCheckIn = vi.fn()
    renderWithProviders(
      <PetScreenEmpty onCheckIn={onCheckIn} isCheckingIn={false} />,
    )

    expect(
      screen.getByRole('heading', { name: 'Яйцо вот-вот треснет' }),
    ).toBeInTheDocument()

    await userEvent.click(screen.getByRole('button', { name: 'Отметиться' }))
    expect(onCheckIn).toHaveBeenCalledTimes(1)
  })

  it('blocks a second check in while one is running', async () => {
    const onCheckIn = vi.fn()
    renderWithProviders(<PetScreenEmpty onCheckIn={onCheckIn} isCheckingIn />)

    const button = screen.getByRole('button', { name: 'Отмечаемся…' })
    expect(button).toBeDisabled()

    await userEvent.click(button)
    expect(onCheckIn).not.toHaveBeenCalled()
  })
})

const actionProps = {
  canCheckIn: true,
  isStroking: false,
  isCheckingIn: false,
  strokeError: null,
  checkInError: null,
  onStroke: () => {},
  onCheckIn: () => {},
}

describe('PetActions', () => {
  it('runs both actions', async () => {
    const onStroke = vi.fn()
    const onCheckIn = vi.fn()
    renderWithProviders(
      <PetActions {...actionProps} onStroke={onStroke} onCheckIn={onCheckIn} />,
    )

    await userEvent.click(
      screen.getByRole('button', { name: 'Погладить питомца' }),
    )
    await userEvent.click(screen.getByRole('button', { name: 'Отметиться' }))

    expect(onStroke).toHaveBeenCalledTimes(1)
    expect(onCheckIn).toHaveBeenCalledTimes(1)
  })

  it('marks the day as done when a check in is no longer possible', () => {
    renderWithProviders(<PetActions {...actionProps} canCheckIn={false} />)

    const button = screen.getByRole('button', {
      name: 'Вы уже отметились сегодня',
    })
    expect(button).toBeDisabled()
    expect(button).toHaveTextContent('Отмечено')
  })

  it('disables stroking while a stroke is in flight', () => {
    renderWithProviders(<PetActions {...actionProps} isStroking />)

    expect(
      screen.getByRole('button', { name: 'Погладить питомца' }),
    ).toBeDisabled()
  })

  it('shows the checking in label while the request runs', () => {
    renderWithProviders(<PetActions {...actionProps} isCheckingIn />)

    expect(
      screen.getByRole('button', { name: 'Отметиться' }),
    ).toHaveTextContent('Отмечаемся…')
  })

  it('prefers the check in error over the stroke error', () => {
    renderWithProviders(
      <PetActions
        {...actionProps}
        strokeError="Не погладили"
        checkInError="Уже отмечались"
      />,
    )

    expect(screen.getByText('Уже отмечались')).toBeInTheDocument()
    expect(screen.queryByText('Не погладили')).not.toBeInTheDocument()
  })

  it('falls back to the stroke error', () => {
    renderWithProviders(
      <PetActions {...actionProps} strokeError="Не погладили" />,
    )

    expect(screen.getByText('Не погладили')).toBeInTheDocument()
  })
})
