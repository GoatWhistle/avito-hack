import { screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { NotFoundScreen } from './NotFoundScreen'
import { renderWithShell } from '#/features/layout/test-utils'

describe('NotFoundScreen', () => {
  it('shows the not found title and hint', () => {
    renderWithShell(<NotFoundScreen />)

    expect(
      screen.getByRole('heading', { name: 'Страница не найдена' }),
    ).toBeInTheDocument()
    expect(
      screen.getByText('Проверьте адрес или вернитесь на главную'),
    ).toBeInTheDocument()
  })

  it('offers a link back to the home page', () => {
    renderWithShell(<NotFoundScreen />)

    const link = screen.getByRole('link', { name: 'На главную' })
    expect(link).toBeInTheDocument()
    expect(link).toHaveAttribute('href', '/')
  })
})
