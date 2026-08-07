import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import {
  makeFavorite,
  renderWithProviders,
} from '#/features/items/components/test-utils'
import { FavoritesScreen } from './FavoritesScreen'

const favoritesList = vi.fn()
const removeFavorite = vi.fn()

vi.mock('#/features/favorites/repository', () => ({
  favoriteRepository: { list: (params: unknown) => favoritesList(params) },
}))

vi.mock('#/features/items/repository', () => ({
  itemRepository: {
    addFavorite: vi.fn(),
    removeFavorite: (id: string) => removeFavorite(id),
  },
}))

beforeEach(() => {
  vi.clearAllMocks()
  favoritesList.mockResolvedValue({
    items: [
      makeFavorite(),
      makeFavorite({ item_id: 'item-2', title: 'Диван', price: 500000 }),
    ],
  })
  removeFavorite.mockResolvedValue(undefined)
})

describe('FavoritesScreen', () => {
  it('shows a skeleton while loading', () => {
    favoritesList.mockReturnValue(new Promise(() => {}))
    renderWithProviders(<FavoritesScreen />)

    expect(screen.getByRole('status', { busy: true })).toBeInTheDocument()
  })

  it('renders saved items with price and status', async () => {
    renderWithProviders(<FavoritesScreen />)

    expect(
      await screen.findByRole('link', { name: 'Велосипед Stels' }),
    ).toHaveAttribute('href', '/items/item-1')
    expect(screen.getByRole('link', { name: 'Диван' })).toBeInTheDocument()
    expect(screen.getByText(/12\s?500/)).toBeInTheDocument()
  })

  it('shows an empty state pointing back to the catalog', async () => {
    favoritesList.mockResolvedValue({ items: [] })
    renderWithProviders(<FavoritesScreen />)

    expect(
      await screen.findByText('В избранном пока пусто'),
    ).toBeInTheDocument()
    expect(screen.getByRole('link', { name: 'Объявления' })).toHaveAttribute(
      'href',
      '/items',
    )
  })

  it('shows an error with a retry action', async () => {
    favoritesList.mockRejectedValue(new Error('offline'))
    renderWithProviders(<FavoritesScreen />)

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Что-то пошло не так',
    )

    favoritesList.mockResolvedValue({ items: [makeFavorite()] })
    await userEvent.click(screen.getByRole('button', { name: 'Повторить' }))

    expect(
      await screen.findByRole('link', { name: 'Велосипед Stels' }),
    ).toBeInTheDocument()
  })

  it('removes an item from favorites optimistically', async () => {
    renderWithProviders(<FavoritesScreen />)

    const buttons = await screen.findAllByRole('button', {
      name: 'Убрать из избранного',
    })

    favoritesList.mockResolvedValue({
      items: [makeFavorite({ item_id: 'item-2', title: 'Диван' })],
    })
    await userEvent.click(buttons[0])

    await waitFor(() => {
      expect(removeFavorite).toHaveBeenCalledWith('item-1')
    })
    await waitFor(() => {
      expect(
        screen.queryByRole('link', { name: 'Велосипед Stels' }),
      ).not.toBeInTheDocument()
    })
  })
})
