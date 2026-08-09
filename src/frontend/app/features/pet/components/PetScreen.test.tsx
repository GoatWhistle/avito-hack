import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ApiError } from '#/api/api-error'
import { PetScreen } from './PetScreen'
import {
  makeCheckedInPet,
  makePet,
  renderWithProviders,
} from './test-utils'

const state = vi.fn()
const stroke = vi.fn()
const feed = vi.fn()
const checkIn = vi.fn()
const summaryToday = vi.fn()

vi.mock('#/features/pet/repository', () => ({
  petRepository: {
    state: () => state(),
    stroke: () => stroke(),
    feed: () => feed(),
    checkIn: () => checkIn(),
    summaryToday: () => summaryToday(),
  },
}))

vi.mock('@lottiefiles/dotlottie-react', () => ({
  DotLottieReact: () => <canvas data-testid="dotlottie-canvas" />,
  setWasmUrl: vi.fn(),
}))

const noSocket = { events: { enabled: false as const } }

const strokeTarget = () => screen.getByRole('button', { name: /погладить/i })

beforeEach(() => {
  vi.clearAllMocks()
  state.mockResolvedValue(makePet())
  stroke.mockResolvedValue(makePet({ happiness: 85 }))
  feed.mockResolvedValue(makePet({ satiety: 85 }))
  checkIn.mockRejectedValue(new Error('check-in is applied by the backend'))
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

  it('reports a failed stroke without breaking the screen', async () => {
    stroke.mockRejectedValue(
      new ApiError({ kind: 'conflict', message: 'already' }),
    )
    renderWithProviders(<PetScreen {...noSocket} />)

    await screen.findByRole('heading', { name: 'Ноти', level: 1 })
    await userEvent.click(strokeTarget())

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

    await userEvent.click(strokeTarget())

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

    await userEvent.click(strokeTarget())

    await waitFor(() => {
      expect(screen.getAllByRole('meter')[1]).toHaveAttribute(
        'aria-valuenow',
        '80',
      )
    })
  })

  it('strokes the lottie raccoon on the teen stage', async () => {
    state.mockResolvedValue(makePet({ stage: 'teen' }))
    renderWithProviders(<PetScreen {...noSocket} />)
    await screen.findByRole('heading', { name: 'Ноти', level: 1 })

    expect(screen.getByTestId('pet-lottie')).toBeInTheDocument()
    await userEvent.click(strokeTarget())

    await waitFor(() => {
      expect(stroke).toHaveBeenCalledTimes(1)
    })
  })

  it('strokes the pet with the keyboard', async () => {
    renderWithProviders(<PetScreen {...noSocket} />)
    await screen.findByRole('heading', { name: 'Ноти', level: 1 })

    strokeTarget().focus()
    await userEvent.keyboard('{Enter}')

    await waitFor(() => {
      expect(stroke).toHaveBeenCalledTimes(1)
    })
  })

  it('offers no manual check-in button', async () => {
    renderWithProviders(<PetScreen {...noSocket} />)
    await screen.findByRole('heading', { name: 'Ноти', level: 1 })

    expect(
      screen.queryByRole('button', { name: /отметиться/i }),
    ).not.toBeInTheDocument()
    expect(checkIn).not.toHaveBeenCalled()
  })

  it('banners the streak when the backend applied a check-in', async () => {
    state.mockResolvedValue(makeCheckedInPet())
    renderWithProviders(<PetScreen {...noSocket} />)

    expect(await screen.findByTestId('celebration-banner')).toHaveTextContent(
      'Серия 5 дней · +10 XP за сегодня',
    )
  })

  it('keeps the banner away when no check-in was applied', async () => {
    renderWithProviders(<PetScreen {...noSocket} />)
    await screen.findByRole('heading', { name: 'Ноти', level: 1 })

    expect(screen.queryByTestId('celebration-banner')).not.toBeInTheDocument()
  })
})
