import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ItemEditScreen } from './ItemEditScreen'
import { makeItem, makePhoto, renderWithProviders } from './test-utils'

const getById = vi.fn()
const listPhotos = vi.fn()
const update = vi.fn()
const changeStatus = vi.fn()

vi.mock('#/features/items/repository', () => ({
  itemRepository: {
    getById: (id: string) => getById(id),
    listPhotos: (id: string) => listPhotos(id),
    update: (id: string, payload: unknown) => update(id, payload),
    changeStatus: (id: string, action: string) => changeStatus(id, action),
    uploadPhoto: vi.fn(),
    deletePhoto: vi.fn(),
  },
}))

beforeEach(() => {
  vi.clearAllMocks()
  getById.mockResolvedValue(makeItem({ status: 'draft' }))
  listPhotos.mockResolvedValue([makePhoto()])
  update.mockResolvedValue(makeItem({ title: 'Новый заголовок' }))
  changeStatus.mockResolvedValue(makeItem({ status: 'published' }))
})

describe('ItemEditScreen', () => {
  it('prefills the form from the loaded item', async () => {
    renderWithProviders(<ItemEditScreen itemId="item-1" />)

    expect(await screen.findByLabelText('Название')).toHaveValue(
      'Велосипед Stels',
    )
    expect(screen.getByLabelText('Описание')).toHaveValue('Отличное состояние')
    expect(screen.getByLabelText('Цена')).toHaveValue(12500)
  })

  it('saves changes and converts the price back to kopeks', async () => {
    renderWithProviders(<ItemEditScreen itemId="item-1" />)

    const title = await screen.findByLabelText('Название')
    await userEvent.clear(title)
    await userEvent.type(title, 'Новый заголовок')
    await userEvent.click(screen.getByRole('button', { name: 'Сохранить' }))

    await waitFor(() => {
      expect(update).toHaveBeenCalledWith('item-1', {
        title: 'Новый заголовок',
        description: 'Отличное состояние',
        price: 1_250_000,
        attributes: {
          category: 'electronics',
          condition: 'used',
        },
      })
    })
    expect(await screen.findByText('Сохранено')).toBeInTheDocument()
  })

  it('reports a failed save', async () => {
    update.mockRejectedValue(new Error('raw backend failure'))
    renderWithProviders(<ItemEditScreen itemId="item-1" />)

    await screen.findByLabelText('Название')
    await userEvent.click(screen.getByRole('button', { name: 'Сохранить' }))

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Неизвестная ошибка',
    )
  })

  it('exposes the photo manager and status actions', async () => {
    renderWithProviders(<ItemEditScreen itemId="item-1" />)

    expect(await screen.findByLabelText('Добавить фото')).toBeEnabled()
    expect(screen.getByText('Действия продавца')).toBeInTheDocument()

    await userEvent.click(
      screen.getByRole('button', { name: 'Отправить на модерацию' }),
    )

    await waitFor(() => {
      expect(changeStatus).toHaveBeenCalledWith('item-1', 'submit')
    })
  })

  it('hides the form from a non-owner', async () => {
    getById.mockResolvedValue(makeItem({ owner_id: 'user-2' }))
    renderWithProviders(<ItemEditScreen itemId="item-1" />)

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Объявление не найдено',
    )
    expect(screen.queryByLabelText('Название')).not.toBeInTheDocument()
  })
})
