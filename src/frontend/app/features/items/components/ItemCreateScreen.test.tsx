import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ItemCreateScreen } from './ItemCreateScreen'
import { makeItem, renderWithProviders } from './test-utils'

const create = vi.fn()

vi.mock('#/features/items/repository', () => ({
  itemRepository: {
    create: (payload: unknown) => create(payload),
  },
}))

beforeEach(() => {
  vi.clearAllMocks()
  create.mockResolvedValue(makeItem())
})

const fill = async (title: string, description: string, price: string) => {
  await userEvent.type(screen.getByLabelText('Название'), title)
  await userEvent.type(screen.getByLabelText('Описание'), description)

  const priceField = screen.getByLabelText('Цена')
  await userEvent.clear(priceField)
  await userEvent.type(priceField, price)
}

describe('ItemCreateScreen', () => {
  it('renders the form fields and the quality hint', () => {
    renderWithProviders(<ItemCreateScreen />)

    expect(screen.getByLabelText('Название')).toBeInTheDocument()
    expect(screen.getByLabelText('Описание')).toBeInTheDocument()
    expect(screen.getByLabelText('Цена')).toBeInTheDocument()
    expect(screen.getByText('Как улучшить объявление')).toBeInTheDocument()
  })

  it('creates an item and converts the price to kopeks', async () => {
    renderWithProviders(<ItemCreateScreen />)

    await fill('Велосипед Stels', 'Отличное состояние', '12500')
    await userEvent.click(screen.getByRole('button', { name: 'Создать' }))

    await waitFor(() => {
      expect(create).toHaveBeenCalledWith({
        title: 'Велосипед Stels',
        description: 'Отличное состояние',
        price: 1_250_000,
        attributes: {
          category: 'electronics',
          condition: 'used',
        },
      })
    })
  })

  it('blocks submission when the title is too short', async () => {
    renderWithProviders(<ItemCreateScreen />)

    await userEvent.type(screen.getByLabelText('Название'), 'ab')

    await waitFor(() => {
      expect(screen.getByLabelText('Название')).toHaveAttribute(
        'aria-invalid',
        'true',
      )
    })

    await userEvent.click(screen.getByRole('button', { name: 'Создать' }))
    expect(create).not.toHaveBeenCalled()
  })

  it('reports a failed creation without losing input', async () => {
    create.mockRejectedValue(new Error('raw backend failure'))
    renderWithProviders(<ItemCreateScreen />)

    await fill('Велосипед Stels', 'Отличное состояние', '100')
    await userEvent.click(screen.getByRole('button', { name: 'Создать' }))

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Неизвестная ошибка',
    )
    expect(screen.getByLabelText('Название')).toHaveValue('Велосипед Stels')
  })

  it('marks the quality checklist as complete for a rich listing', async () => {
    renderWithProviders(<ItemCreateScreen />)

    await userEvent.type(screen.getByLabelText('Описание'), 'а'.repeat(210))
    const priceField = screen.getByLabelText('Цена')
    await userEvent.clear(priceField)
    await userEvent.type(priceField, '100')

    expect(await screen.findByText('Указана цена')).toBeInTheDocument()
    expect(screen.getByText('Есть подробное описание')).toBeInTheDocument()
  }, 20_000)
})
