import { screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { LandingShowcase } from './LandingShowcase'
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
  list.mockResolvedValue({
    items: [
      makeEntry(),
      makeEntry({ id: 'item-2', title: 'Диван', price: 500000 }),
    ],
  })
})

describe('LandingShowcase', () => {
  it('renders live listings with title, price and status', async () => {
    renderWithProviders(<LandingShowcase />)

    expect(
      await screen.findByRole('link', { name: 'Велосипед Stels' }),
    ).toHaveAttribute('href', '/items/item-1')
    expect(screen.getByRole('link', { name: 'Диван' })).toBeInTheDocument()
    expect(screen.getByText(/12\s?500/)).toBeInTheDocument()

    const grid = within(
      screen.getByRole('list', { name: 'Мини-Авито внутри продукта' }),
    )
    expect(grid.getAllByText('Опубликовано')).toHaveLength(2)
  })

  it('requests only published listings', async () => {
    renderWithProviders(<LandingShowcase />)
    await screen.findByRole('link', { name: 'Велосипед Stels' })

    expect(list).toHaveBeenCalledWith(
      expect.objectContaining({ status: 'published' }),
    )
  })

  it('caps the showcase at four listings', async () => {
    list.mockResolvedValue({
      items: Array.from({ length: 9 }, (_, index) =>
        makeEntry({ id: `item-${index}`, title: `Товар ${index}` }),
      ),
    })
    renderWithProviders(<LandingShowcase />)

    const grid = await screen.findByRole('list', {
      name: 'Мини-Авито внутри продукта',
    })
    expect(within(grid).getAllByRole('listitem')).toHaveLength(4)
  })

  it('shows a skeleton while the listings load', () => {
    list.mockReturnValue(new Promise(() => {}))
    renderWithProviders(<LandingShowcase />)

    expect(screen.getByRole('status', { busy: true })).toBeInTheDocument()
  })

  it('shows an error with a retry action', async () => {
    list.mockRejectedValue(new Error('offline'))
    renderWithProviders(<LandingShowcase />)

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Что-то пошло не так',
    )

    list.mockResolvedValue({ items: [makeEntry()] })
    await userEvent.click(screen.getByRole('button', { name: 'Повторить' }))

    expect(
      await screen.findByRole('link', { name: 'Велосипед Stels' }),
    ).toBeInTheDocument()
  })

  it('shows an empty state when there are no listings', async () => {
    list.mockResolvedValue({ items: [] })
    renderWithProviders(<LandingShowcase />)

    await waitFor(() => {
      expect(screen.getByRole('status')).toHaveTextContent(
        'Объявлений пока нет',
      )
    })
  })

  it('always links to the full catalogue', async () => {
    renderWithProviders(<LandingShowcase />)

    expect(
      screen.getByRole('link', { name: 'Все объявления' }),
    ).toHaveAttribute('href', '/items')
  })
})
