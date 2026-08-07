import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ItemDetailScreen } from './ItemDetailScreen'
import { makeItem, makePhoto, renderWithProviders } from './test-utils'

const getById = vi.fn()
const listPhotos = vi.fn()
const changeStatus = vi.fn()
const addFavorite = vi.fn()
const removeFavorite = vi.fn()
const favoritesList = vi.fn()

vi.mock('#/features/items/repository', () => ({
  itemRepository: {
    getById: (id: string) => getById(id),
    listPhotos: (id: string) => listPhotos(id),
    changeStatus: (id: string, action: string) => changeStatus(id, action),
    addFavorite: (id: string) => addFavorite(id),
    removeFavorite: (id: string) => removeFavorite(id),
  },
}))

vi.mock('#/features/favorites/repository', () => ({
  favoriteRepository: { list: () => favoritesList() },
}))

vi.mock('#/api', async () => {
  const actual = await vi.importActual<Record<string, unknown>>('#/api')

  return { ...actual, getToken: () => 'token' }
})

beforeEach(() => {
  vi.clearAllMocks()
  getById.mockResolvedValue(makeItem())
  listPhotos.mockResolvedValue([makePhoto(), makePhoto({ id: 'photo-2' })])
  changeStatus.mockResolvedValue(makeItem({ status: 'sold' }))
  favoritesList.mockResolvedValue({ items: [] })
  addFavorite.mockResolvedValue(undefined)
  removeFavorite.mockResolvedValue(undefined)
})

describe('ItemDetailScreen', () => {
  it('shows a skeleton while loading', () => {
    getById.mockReturnValue(new Promise(() => {}))
    renderWithProviders(<ItemDetailScreen itemId="item-1" />)

    expect(screen.getByRole('status', { busy: true })).toBeInTheDocument()
  })

  it('renders title, price, status, description and the gallery', async () => {
    renderWithProviders(<ItemDetailScreen itemId="item-1" />)

    expect(
      await screen.findByRole('heading', { name: 'Велосипед Stels', level: 1 }),
    ).toBeInTheDocument()
    expect(screen.getByText(/12\s?500/)).toBeInTheDocument()
    expect(screen.getByText('Отличное состояние')).toBeInTheDocument()
    expect(screen.getByRole('list', { name: 'Фотографии' })).toBeInTheDocument()
  })

  it('shows a not-found error with a retry action', async () => {
    getById.mockRejectedValue(new Error('404'))
    renderWithProviders(<ItemDetailScreen itemId="item-1" />)

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Объявление не найдено',
    )
    expect(
      screen.getByRole('button', { name: 'Повторить' }),
    ).toBeInTheDocument()
  })

  it('offers owner actions and marks the item as sold', async () => {
    renderWithProviders(<ItemDetailScreen itemId="item-1" />)

    await screen.findByRole('heading', { name: 'Велосипед Stels', level: 1 })
    expect(screen.getByText('Действия продавца')).toBeInTheDocument()

    await userEvent.click(
      screen.getByRole('button', { name: 'Отметить проданным' }),
    )

    await waitFor(() => {
      expect(changeStatus).toHaveBeenCalledWith('item-1', 'sell')
    })
    expect(
      await screen.findByText('Поздравляем с продажей!'),
    ).toBeInTheDocument()
    expect(
      screen.getByText('Питомец получил опыт за вашу сделку'),
    ).toBeInTheDocument()
  })

  it('reports a failed status change', async () => {
    changeStatus.mockRejectedValue(new Error('Нельзя продать'))
    renderWithProviders(<ItemDetailScreen itemId="item-1" />)

    await screen.findByRole('heading', { name: 'Велосипед Stels', level: 1 })
    await userEvent.click(
      screen.getByRole('button', { name: 'Отметить проданным' }),
    )

    expect(await screen.findByRole('alert')).toHaveTextContent('Нельзя продать')
  })

  it('shows the favorite toggle for a foreign item instead of owner actions', async () => {
    getById.mockResolvedValue(makeItem({ owner_id: 'user-2' }))
    renderWithProviders(<ItemDetailScreen itemId="item-1" />)

    await screen.findByRole('heading', { name: 'Велосипед Stels', level: 1 })
    expect(screen.queryByText('Действия продавца')).not.toBeInTheDocument()

    await userEvent.click(screen.getByRole('button', { name: 'В избранное' }))

    await waitFor(() => {
      expect(addFavorite).toHaveBeenCalledWith('item-1')
    })
  })
})
