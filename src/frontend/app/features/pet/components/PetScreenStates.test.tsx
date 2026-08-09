import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { PetActions } from './PetActions'
import { PetScreenError, PetScreenSkeleton } from './PetScreenStates'
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

describe('PetActions', () => {
  it('offers no check-in button because check-in is automatic', () => {
    renderWithProviders(<PetActions strokeError={null} />)

    expect(screen.queryByRole('button')).not.toBeInTheDocument()
  })

  it('offers no separate stroke button', () => {
    renderWithProviders(<PetActions strokeError={null} />)

    expect(
      screen.queryByRole('button', { name: 'Погладить питомца' }),
    ).not.toBeInTheDocument()
  })

  it('keeps an empty live region when nothing failed', () => {
    renderWithProviders(<PetActions strokeError={null} />)

    const status = screen.getByRole('status')
    expect(status).toHaveAttribute('aria-live', 'polite')
    expect(status).toHaveTextContent('')
  })

  it('announces a stroke failure in a live region', () => {
    renderWithProviders(<PetActions strokeError="Не погладили" />)

    const status = screen.getByRole('status')
    expect(status).toHaveAttribute('aria-live', 'polite')
    expect(status).toHaveTextContent('Не погладили')
  })
})
