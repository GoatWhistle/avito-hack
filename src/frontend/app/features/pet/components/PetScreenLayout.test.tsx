import { screen } from '@testing-library/react'
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

const noSocket = { events: { enabled: false as const } }

beforeEach(() => {
  vi.clearAllMocks()
  state.mockResolvedValue(makePet())
  summaryToday.mockResolvedValue(null)
})

describe('PetScreen layout', () => {
  it('puts the hud, stats, stage and rewards on one screen', async () => {
    renderWithProviders(<PetScreen {...noSocket} />)

    expect(
      await screen.findByRole('heading', { name: 'Ноти · уровень 3' }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('region', { name: 'Состояние питомца' }),
    ).toBeInTheDocument()
    expect(screen.getByTestId('pet-speech')).toBeInTheDocument()
    expect(screen.getByTestId('pet-rewards-panel')).toBeInTheDocument()
  })

  it('shows the pet stage before the side rails on small screens', async () => {
    renderWithProviders(<PetScreen {...noSocket} />)
    await screen.findByRole('heading', { name: 'Ноти', level: 1 })

    const stage = screen.getByTestId('pet-speech').closest('div.order-1')
    const stats = screen
      .getByRole('region', { name: 'Состояние питомца' })
      .closest('div.order-2')

    expect(stage).not.toBeNull()
    expect(stats).not.toBeNull()
    expect(stage).toHaveClass('lg:order-2')
    expect(stats).toHaveClass('lg:order-1')
  })

  it('keeps the level progress out of the stage to avoid duplication', async () => {
    renderWithProviders(<PetScreen {...noSocket} />)
    await screen.findByRole('heading', { name: 'Ноти', level: 1 })

    expect(screen.getAllByRole('progressbar')).toHaveLength(1)
  })

  it('renders the dashboard right away for a freshly created baby', async () => {
    state.mockResolvedValue(makePet({ stage: 'baby' }))
    renderWithProviders(<PetScreen {...noSocket} />)

    await screen.findByRole('heading', { name: 'Ноти', level: 1 })

    expect(screen.getByTestId('pet-rewards-panel')).toBeInTheDocument()
    expect(screen.getByTestId('pet-speech')).toBeInTheDocument()
  })
})
