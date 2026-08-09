import { screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { LandingScreen } from './LandingScreen'
import {
  makeEntry,
  renderWithProviders,
} from '#/features/items/components/test-utils'

const list = vi.fn()

vi.mock('#/features/items/repository', () => ({
  itemRepository: {
    list: (params: unknown) => list(params),
  },
}))

beforeEach(() => {
  vi.clearAllMocks()
  list.mockResolvedValue({ items: [makeEntry()] })
})

describe('LandingScreen', () => {
  it('shows the product, the showcase and the marketplace value', async () => {
    renderWithProviders(<LandingScreen />)

    expect(
      screen.getByRole('heading', { name: 'Мини-Авито внутри продукта' }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('heading', {
        name: 'Что это даёт метрикам площадки',
      }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('heading', { name: 'Зачем это Авито' }),
    ).toBeInTheDocument()
    expect(
      await screen.findByRole('link', { name: 'Велосипед Stels' }),
    ).toBeInTheDocument()
  })

  it('offers a guest the sign-up call to action', () => {
    renderWithProviders(<LandingScreen />, { userId: null })

    expect(
      screen.getByRole('link', { name: 'Завести питомца' }),
    ).toHaveAttribute('href', '/sign-up')
    expect(
      screen.getByRole('heading', { name: 'Заведите Ноти за минуту' }),
    ).toBeInTheDocument()
  })

  it('offers an authenticated user a path to the pet and hides the sign-up block', () => {
    renderWithProviders(<LandingScreen />)

    expect(
      screen.getByRole('link', { name: 'К моему питомцу' }),
    ).toHaveAttribute('href', '/pet')
    expect(
      screen.queryByRole('heading', { name: 'Заведите Ноти за минуту' }),
    ).not.toBeInTheDocument()
  })
})
