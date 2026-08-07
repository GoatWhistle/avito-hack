import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { MyItemsScreen } from './MyItemsScreen'
import { makeEntry, makeItem, renderWithProviders } from './test-utils'

const listMine = vi.fn()
const changeStatus = vi.fn()

vi.mock('#/features/items/repository', () => ({
  itemRepository: {
    listMine: (params: unknown) => listMine(params),
    changeStatus: (id: string, action: string) => changeStatus(id, action),
  },
}))

beforeEach(() => {
  vi.clearAllMocks()
  listMine.mockResolvedValue({
    items: [makeEntry({ status: 'draft' })],
    next_cursor: 'cursor-2',
  })
  changeStatus.mockResolvedValue(makeItem({ status: 'published' }))
})

describe('MyItemsScreen', () => {
  it('shows a skeleton while loading', () => {
    listMine.mockReturnValue(new Promise(() => {}))
    renderWithProviders(<MyItemsScreen />)

    expect(screen.getByRole('status', { busy: true })).toBeInTheDocument()
  })

  it('lists own items with an edit link and quick actions', async () => {
    renderWithProviders(<MyItemsScreen />)

    expect(
      await screen.findByRole('link', { name: 'Велосипед Stels' }),
    ).toBeInTheDocument()
    expect(screen.getByRole('link', { name: 'Редактировать' })).toHaveAttribute(
      'href',
      '/items/item-1/edit',
    )
    expect(
      screen.getByRole('button', { name: 'Опубликовать' }),
    ).toBeInTheDocument()
  })

  it('publishes a draft from the list', async () => {
    renderWithProviders(<MyItemsScreen />)
    await screen.findByRole('link', { name: 'Велосипед Stels' })

    await userEvent.click(screen.getByRole('button', { name: 'Опубликовать' }))

    await waitFor(() => {
      expect(changeStatus).toHaveBeenCalledWith('item-1', 'publish')
    })
  })

  it('offers a call to action when there is nothing yet', async () => {
    listMine.mockResolvedValue({ items: [] })
    renderWithProviders(<MyItemsScreen />)

    expect(
      await screen.findByText('Вы ещё не разместили объявлений'),
    ).toBeInTheDocument()
    expect(
      screen.getByText(
        'Разместите первое объявление — питомец получит за это опыт',
      ),
    ).toBeInTheDocument()
  })

  it('loads the next page on demand', async () => {
    renderWithProviders(<MyItemsScreen />)
    await screen.findByRole('link', { name: 'Велосипед Stels' })

    await userEvent.click(screen.getByRole('button', { name: 'Показать ещё' }))

    await waitFor(() => {
      expect(listMine).toHaveBeenCalledWith(
        expect.objectContaining({ cursor: 'cursor-2' }),
      )
    })
  })

  it('shows an error with a retry action', async () => {
    listMine.mockRejectedValue(new Error('offline'))
    renderWithProviders(<MyItemsScreen />)

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Что-то пошло не так',
    )
    expect(
      screen.getByRole('button', { name: 'Повторить' }),
    ).toBeInTheDocument()
  })
})
