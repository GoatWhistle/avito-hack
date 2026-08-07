import { screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ItemsScreen } from './ItemsScreen'
import { makeEntry, renderWithProviders } from './test-utils'

const list = vi.fn()
const addFavorite = vi.fn()
const removeFavorite = vi.fn()
const favoritesList = vi.fn()

vi.mock('#/features/items/repository', () => ({
  itemRepository: {
    list: (params: unknown) => list(params),
    addFavorite: (id: string) => addFavorite(id),
    removeFavorite: (id: string) => removeFavorite(id),
  },
}))

vi.mock('#/features/favorites/repository', () => ({
  favoriteRepository: {
    list: (params: unknown) => favoritesList(params),
  },
}))

vi.mock('#/api', async () => {
  const actual = await vi.importActual<Record<string, unknown>>('#/api')

  return { ...actual, getToken: () => 'token' }
})

beforeEach(() => {
  vi.clearAllMocks()
  list.mockResolvedValue({
    items: [
      makeEntry(),
      makeEntry({ id: 'item-2', title: 'Диван', price: 500000 }),
    ],
  })
  favoritesList.mockResolvedValue({ items: [] })
  addFavorite.mockResolvedValue(undefined)
  removeFavorite.mockResolvedValue(undefined)
})

describe('ItemsScreen', () => {
  it('shows a skeleton while the first page loads', () => {
    list.mockReturnValue(new Promise(() => {}))
    renderWithProviders(<ItemsScreen />)

    expect(screen.getByRole('status', { busy: true })).toBeInTheDocument()
  })

  it('renders cards with title, price and status', async () => {
    renderWithProviders(<ItemsScreen />)

    expect(
      await screen.findByRole('link', { name: 'Велосипед Stels' }),
    ).toHaveAttribute('href', '/items/item-1')
    expect(screen.getByRole('link', { name: 'Диван' })).toBeInTheDocument()
    expect(screen.getByText(/12\s?500/)).toBeInTheDocument()

    const grid = within(screen.getByRole('list', { name: 'Объявления' }))
    expect(grid.getAllByText('Опубликовано')).toHaveLength(2)
  })

  it('shows an empty state with a call to action', async () => {
    list.mockResolvedValue({ items: [] })
    renderWithProviders(<ItemsScreen />)

    expect(await screen.findByText('Объявлений пока нет')).toBeInTheDocument()
    expect(
      screen.getAllByRole('link', { name: 'Разместить объявление' }).length,
    ).toBeGreaterThan(0)
  })

  it('shows an error with a retry action', async () => {
    list.mockRejectedValue(new Error('offline'))
    renderWithProviders(<ItemsScreen />)

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Что-то пошло не так',
    )

    list.mockResolvedValue({ items: [makeEntry()] })
    await userEvent.click(screen.getByRole('button', { name: 'Повторить' }))

    expect(
      await screen.findByRole('link', { name: 'Велосипед Stels' }),
    ).toBeInTheDocument()
  })

  it('filters by status', async () => {
    renderWithProviders(<ItemsScreen />)
    await screen.findByRole('link', { name: 'Велосипед Stels' })

    await userEvent.click(screen.getByRole('button', { name: 'Продано' }))

    await waitFor(() => {
      expect(list).toHaveBeenCalledWith(
        expect.objectContaining({ status: 'sold' }),
      )
    })
  })

  it('searches with a debounced query', async () => {
    renderWithProviders(<ItemsScreen />)
    await screen.findByRole('link', { name: 'Велосипед Stels' })

    await userEvent.type(screen.getByRole('searchbox'), 'диван')

    await waitFor(
      () => {
        expect(list).toHaveBeenCalledWith(
          expect.objectContaining({ search: 'диван' }),
        )
      },
      { timeout: 2000 },
    )
  })

  it('adds an item to favorites optimistically', async () => {
    renderWithProviders(<ItemsScreen />)
    await screen.findByRole('link', { name: 'Велосипед Stels' })

    const buttons = await screen.findAllByRole('button', {
      name: 'В избранное',
    })
    await userEvent.click(buttons[0])

    await waitFor(() => {
      expect(addFavorite).toHaveBeenCalledWith('item-1')
    })
  })
})
