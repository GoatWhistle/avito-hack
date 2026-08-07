import { screen, waitFor } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { LevelIndicator } from './LevelIndicator'
import { levelFloor, xpProgressRatio, type Progress } from './progress.types'
import { makeUser, renderWithShell } from '#/features/layout/test-utils'

const profile = vi.fn()

vi.mock('./progress.repository', () => ({
  progressRepository: {
    profile: () => profile(),
  },
}))

const makeProgress = (overrides: Partial<Progress> = {}): Progress => ({
  id: 'raccoon-1',
  name: 'Ноти',
  level: 4,
  xp: 30,
  xpToNextLevel: 35,
  currentStreak: 3,
  badges: [],
  earnedBadgeCount: 0,
  ...overrides,
})

const session = { user: makeUser(), isAuthenticated: true }

beforeEach(() => {
  vi.clearAllMocks()
  profile.mockResolvedValue(makeProgress())
})

describe('levelFloor', () => {
  it('maps a level to its xp threshold', () => {
    expect(levelFloor(1)).toBe(0)
    expect(levelFloor(4)).toBe(22)
  })

  it('clamps out of range levels', () => {
    expect(levelFloor(0)).toBe(0)
    expect(levelFloor(-3)).toBe(0)
    expect(levelFloor(999)).toBe(420)
  })
})

describe('xpProgressRatio', () => {
  it('reports a partial ratio between thresholds', () => {
    expect(
      xpProgressRatio(makeProgress({ level: 4, xp: 28, xpToNextLevel: 35 })),
    ).toBeCloseTo(6 / 13)
  })

  it('treats a non positive next level as complete', () => {
    expect(xpProgressRatio(makeProgress({ xpToNextLevel: 0 }))).toBe(1)
    expect(xpProgressRatio(makeProgress({ xpToNextLevel: -5 }))).toBe(1)
  })

  it('treats a collapsed span as complete', () => {
    expect(
      xpProgressRatio(makeProgress({ level: 4, xp: 25, xpToNextLevel: 22 })),
    ).toBe(1)
  })

  it('clamps below zero and above one', () => {
    expect(
      xpProgressRatio(makeProgress({ level: 4, xp: 0, xpToNextLevel: 35 })),
    ).toBe(0)
    expect(
      xpProgressRatio(makeProgress({ level: 4, xp: 400, xpToNextLevel: 35 })),
    ).toBe(1)
  })
})

describe('LevelIndicator', () => {
  it('shows a skeleton while the profile loads', () => {
    profile.mockReturnValue(new Promise(() => {}))

    renderWithShell(<LevelIndicator />, { session })

    expect(screen.getByTestId('level-indicator-skeleton')).toBeInTheDocument()
  })

  it('renders the level, progress bar and streak', async () => {
    renderWithShell(<LevelIndicator />, { session })

    const indicator = await screen.findByTestId('level-indicator')
    expect(indicator).toHaveAccessibleName(/Уровень 4/)
    expect(screen.getByRole('progressbar')).toHaveAttribute(
      'aria-valuenow',
      '62',
    )
    expect(screen.getByLabelText('Серия из 3 дней')).toBeInTheDocument()
  })

  it('hides the streak badge at zero days', async () => {
    profile.mockResolvedValue(makeProgress({ currentStreak: 0 }))

    renderWithShell(<LevelIndicator />, { session })

    await screen.findByTestId('level-indicator')
    expect(screen.queryByLabelText(/Серия из/)).not.toBeInTheDocument()
  })

  it('announces the max level when there is no next level', async () => {
    profile.mockResolvedValue(makeProgress({ xpToNextLevel: 0 }))

    renderWithShell(<LevelIndicator />, { session })

    expect(await screen.findByText('Максимальный уровень')).toBeInTheDocument()
    expect(screen.getByRole('progressbar')).toHaveAttribute(
      'aria-valuenow',
      '100',
    )
  })

  it('hides the xp caption in compact mode', async () => {
    renderWithShell(<LevelIndicator compact />, { session })

    await screen.findByTestId('level-indicator')
    expect(screen.queryByText(/XP, до уровня/)).not.toBeInTheDocument()
  })

  it('renders nothing when the profile request fails', async () => {
    profile.mockRejectedValue(new Error('boom'))

    renderWithShell(<LevelIndicator />, { session })

    await waitFor(() => {
      expect(
        screen.queryByTestId('level-indicator-skeleton'),
      ).not.toBeInTheDocument()
    })
    expect(screen.queryByTestId('level-indicator')).not.toBeInTheDocument()
  })
})
