import { act, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { PetScreen } from './PetScreen'
import { makePet, renderWithProviders } from './test-utils'

const state = vi.fn()
const summaryToday = vi.fn()

vi.mock('#/features/pet/repository', () => ({
  petRepository: {
    state: () => state(),
    stroke: () => Promise.resolve(makePet()),
    feed: () => Promise.resolve(makePet()),
    checkIn: () => Promise.reject(new Error('unused')),
    summaryToday: () => summaryToday(),
  },
}))

class FakeSocket {
  static last: FakeSocket | null = null
  onopen: (() => void) | null = null
  onmessage: ((event: MessageEvent<string>) => void) | null = null
  onerror: (() => void) | null = null
  onclose: (() => void) | null = null

  constructor() {
    FakeSocket.last = this
  }

  send() {}
  close() {}

  emit(type: string, payload: unknown) {
    this.onmessage?.({
      data: JSON.stringify({ type, payload }),
    } as MessageEvent<string>)
  }
}

const socketOptions = {
  events: {
    token: 'jwt',
    url: '/api/v1/ws',
    factory: () => new FakeSocket() as unknown as WebSocket,
  },
}

beforeEach(() => {
  vi.clearAllMocks()
  FakeSocket.last = null
  state.mockResolvedValue(makePet())
  summaryToday.mockResolvedValue(null)
})

const emit = async (type: string, payload: unknown) => {
  await act(async () => {
    FakeSocket.last?.emit(type, payload)
  })
}

describe('PetScreen realtime', () => {
  it('opens a socket and applies pet.updated', async () => {
    renderWithProviders(<PetScreen {...socketOptions} />)
    await screen.findByRole('heading', { name: 'Ноти', level: 1 })

    expect(FakeSocket.last).not.toBeNull()

    await emit('pet.updated', makePet({ satiety: 20, happiness: 25 }))

    await waitFor(() => {
      expect(screen.getAllByRole('meter')[0]).toHaveAttribute(
        'aria-valuenow',
        '20',
      )
    })
    expect(screen.getByText(/питомец голоден/i)).toBeInTheDocument()
  })

  it('shows a floating toast on xp.gained', async () => {
    renderWithProviders(<PetScreen {...socketOptions} />)
    await screen.findByRole('heading', { name: 'Ноти', level: 1 })

    await emit('xp.gained', { amount: 25 })

    expect(await screen.findByTestId('xp-toast')).toHaveTextContent('+25 XP')
  })

  it('celebrates level.up and updates the level', async () => {
    renderWithProviders(<PetScreen {...socketOptions} />)
    await screen.findByRole('heading', { name: 'Ноти', level: 1 })

    await emit('level.up', { level: 7 })

    expect(await screen.findByTestId('celebration-banner')).toHaveTextContent(
      'Новый уровень: 7',
    )
    expect(
      screen.getByRole('heading', { name: 'Ноти · уровень 7' }),
    ).toBeInTheDocument()
  })

  it('applies pet.state from the socket', async () => {
    renderWithProviders(<PetScreen {...socketOptions} />)
    await screen.findByRole('heading', { name: 'Ноти', level: 1 })

    await emit('pet.state', makePet({ stage: 'teen', energy: 44 }))

    await waitFor(() => {
      expect(screen.getAllByRole('meter')[2]).toHaveAttribute(
        'aria-valuenow',
        '44',
      )
    })
    expect(screen.getAllByRole('meter')).toHaveLength(3)
  })

  it('updates the streak from streak.updated and shows milestones', async () => {
    renderWithProviders(<PetScreen {...socketOptions} />)
    await screen.findByRole('heading', { name: 'Ноти', level: 1 })

    await emit('streak.updated', { days: 7, milestone: true })

    expect(await screen.findByText('7 дней подряд')).toBeInTheDocument()
    expect(screen.getByTestId('celebration-banner')).toHaveTextContent(
      'Рубеж 7 дней достигнут',
    )
  })

  it('announces granted rewards and lets them be dismissed', async () => {
    renderWithProviders(<PetScreen {...socketOptions} />)
    await screen.findByRole('heading', { name: 'Ноти', level: 1 })

    await emit('reward.granted', {
      reward_id: 'r1',
      title: 'Бесплатная доставка',
    })

    const banner = await screen.findByTestId('celebration-banner')
    expect(banner).toHaveTextContent('Бесплатная доставка')

    await userEvent.click(screen.getByRole('button', { name: /закрыть/i }))
    expect(screen.queryByTestId('celebration-banner')).not.toBeInTheDocument()
  })
})
