import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ApiError } from '#/api/api-error'
import { PetScreen } from './PetScreen'
import { makePet, renderWithProviders } from './test-utils'

const state = vi.fn()
const stroke = vi.fn()
const checkIn = vi.fn()
const summaryToday = vi.fn()

vi.mock('#/features/pet/repository', () => ({
  petRepository: {
    state: () => state(),
    stroke: () => stroke(),
    checkIn: () => checkIn(),
    summaryToday: () => summaryToday(),
  },
}))

const noSocket = { events: { enabled: false as const } }

beforeEach(() => {
  vi.clearAllMocks()
  state.mockResolvedValue(makePet())
  stroke.mockResolvedValue(makePet({ happiness: 85 }))
  checkIn.mockResolvedValue({
    pet: makePet({ level: 4, xp: 25 }),
    xp_granted: 10,
    level: 4,
    previous_level: 3,
    next_level_xp: 35,
    unlocked_rewards: [],
    streak: {
      days: 5,
      continued: true,
      freeze_used: false,
      reset: false,
      milestone_bonus: 0,
      milestone_reached: 0,
      freezes_left: 1,
    },
  })
  summaryToday.mockResolvedValue(null)
})

describe('PetScreen rendering', () => {
  it('shows a skeleton while loading', () => {
    state.mockReturnValue(new Promise(() => {}))
    renderWithProviders(<PetScreen {...noSocket} />)

    expect(screen.getByRole('status', { busy: true })).toBeInTheDocument()
  })

  it('renders the pet, stats, level and streak', async () => {
    renderWithProviders(<PetScreen {...noSocket} />)

    expect(
      await screen.findByRole('heading', { name: 'Ноти', level: 1 }),
    ).toBeInTheDocument()

    const meters = screen.getAllByRole('meter')
    expect(meters).toHaveLength(3)
    expect(meters[0]).toHaveAttribute('aria-valuenow', '70')

    expect(screen.getByRole('progressbar')).toHaveAttribute(
      'aria-valuenow',
      '68',
    )
    expect(screen.getByText('4 дня подряд')).toBeInTheDocument()
  })

  it('renders the egg state when the pet has not hatched', async () => {
    state.mockResolvedValue(makePet({ is_hatched: false, stage: 'egg' }))
    renderWithProviders(<PetScreen {...noSocket} />)

    expect(
      await screen.findByRole('heading', { name: /яйцо вот-вот треснет/i }),
    ).toBeInTheDocument()
    expect(screen.queryAllByRole('meter')).toHaveLength(0)
  })

  it('marks the max level with a dedicated state', async () => {
    state.mockResolvedValue(makePet({ level: 15, xp: 500, next_level_xp: 0 }))
    renderWithProviders(<PetScreen {...noSocket} />)

    expect(await screen.findByText('МАКС')).toBeInTheDocument()
    expect(screen.getByText(/легендой/i)).toBeInTheDocument()
    expect(screen.queryByText('500 / 0 XP')).not.toBeInTheDocument()
  })
})

describe('PetScreen errors', () => {
  it('shows an error with a retry action', async () => {
    state.mockRejectedValue(
      new ApiError({ kind: 'network', message: 'offline' }),
    )
    renderWithProviders(<PetScreen {...noSocket} />)

    const alert = await screen.findByRole('alert')
    expect(alert).toHaveTextContent('Нет связи с сервером')

    state.mockResolvedValue(makePet())
    await userEvent.click(screen.getByRole('button', { name: /попробовать/i }))

    expect(
      await screen.findByRole('heading', { name: 'Ноти', level: 1 }),
    ).toBeInTheDocument()
  })

  it('reports a failed check-in without breaking the screen', async () => {
    checkIn.mockRejectedValue(
      new ApiError({ kind: 'conflict', message: 'already' }),
    )
    renderWithProviders(<PetScreen {...noSocket} />)

    await screen.findByRole('heading', { name: 'Ноти', level: 1 })
    await userEvent.click(screen.getByRole('button', { name: /отметиться/i }))

    expect(await screen.findByText('Конфликт состояния')).toBeInTheDocument()
    expect(screen.getAllByRole('meter')).toHaveLength(3)
  })
})

describe('PetScreen actions', () => {
  it('strokes the pet optimistically', async () => {
    let resolveStroke: (pet: ReturnType<typeof makePet>) => void = () => {}
    stroke.mockReturnValue(
      new Promise((resolve) => {
        resolveStroke = resolve
      }),
    )

    renderWithProviders(<PetScreen {...noSocket} />)
    await screen.findByRole('heading', { name: 'Ноти', level: 1 })

    await userEvent.click(screen.getByRole('button', { name: /погладить/i }))

    await waitFor(() => {
      expect(screen.getAllByRole('meter')[1]).toHaveAttribute(
        'aria-valuenow',
        '85',
      )
    })

    resolveStroke(makePet({ happiness: 85 }))
    expect(stroke).toHaveBeenCalledTimes(1)
  })

  it('rolls the optimistic stroke back on failure', async () => {
    stroke.mockRejectedValue(
      new ApiError({ kind: 'internal_error', message: 'boom' }),
    )
    renderWithProviders(<PetScreen {...noSocket} />)
    await screen.findByRole('heading', { name: 'Ноти', level: 1 })

    await userEvent.click(screen.getByRole('button', { name: /погладить/i }))

    await waitFor(() => {
      expect(screen.getAllByRole('meter')[1]).toHaveAttribute(
        'aria-valuenow',
        '80',
      )
    })
  })

  it('checks in and celebrates the level up', async () => {
    renderWithProviders(<PetScreen {...noSocket} />)
    await screen.findByRole('heading', { name: 'Ноти', level: 1 })

    await userEvent.click(screen.getByRole('button', { name: /отметиться/i }))

    expect(await screen.findByTestId('celebration-banner')).toHaveTextContent(
      'Новый уровень: 4',
    )
    expect(checkIn).toHaveBeenCalledTimes(1)
  })

  it('disables check-in when it already happened today', async () => {
    state.mockResolvedValue(
      makePet({ last_checkin_date: new Date().toISOString() }),
    )
    renderWithProviders(<PetScreen {...noSocket} />)

    const button = await screen.findByRole('button', {
      name: /уже отметились/i,
    })
    expect(button).toBeDisabled()
  })
})
