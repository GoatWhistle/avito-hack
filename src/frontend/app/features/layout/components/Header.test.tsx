import { screen, waitFor } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { Header } from './Header'
import { makeUser, renderWithShell } from '#/features/layout/test-utils'

const profile = vi.fn()

vi.mock('#/features/progress/progress.repository', () => ({
  progressRepository: {
    profile: () => profile(),
  },
}))

const makeProgress = () => ({
  id: 'raccoon-1',
  name: 'Ноти',
  level: 4,
  xp: 30,
  xpToNextLevel: 35,
  currentStreak: 3,
  badges: [],
  earnedBadgeCount: 0,
})

beforeEach(() => {
  vi.clearAllMocks()
  profile.mockResolvedValue(makeProgress())
})

describe('Header', () => {
  it('shows primary navigation links', () => {
    renderWithShell(<Header />, {
      session: { user: makeUser(), isAuthenticated: true },
    })

    const navs = screen.getAllByRole('navigation', {
      name: 'Основная навигация',
    })
    expect(navs.length).toBeGreaterThan(0)
    expect(
      screen.getAllByRole('link', { name: /Питомец/ }).length,
    ).toBeGreaterThan(0)
  })

  it('keeps the level and xp indicator out of the header for a signed in user', async () => {
    renderWithShell(<Header />, {
      session: { user: makeUser(), isAuthenticated: true },
    })

    await waitFor(() => {
      expect(
        screen.getAllByRole('link', { name: /Питомец/ }).length,
      ).toBeGreaterThan(0)
    })

    expect(screen.queryByTestId('level-indicator')).not.toBeInTheDocument()
    expect(
      screen.queryByRole('progressbar', { name: 'Прогресс опыта' }),
    ).not.toBeInTheDocument()
  })

  it('offers a sign in action to anonymous visitors and hides the indicator', () => {
    renderWithShell(<Header />, { session: { isAuthenticated: false } })

    expect(screen.getByRole('link', { name: 'Войти' })).toBeInTheDocument()
    expect(screen.queryByTestId('level-indicator')).not.toBeInTheDocument()
  })
})
