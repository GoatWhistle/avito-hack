import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { OnboardingCard } from './OnboardingCard'
import { onboardingStorageKey } from './onboarding-state'
import { makeUser, renderWithShell } from '#/features/layout/test-utils'

const profile = vi.fn()

vi.mock('#/features/progress/progress.repository', () => ({
  progressRepository: {
    profile: () => profile(),
  },
}))

const makeProgress = (overrides: Record<string, unknown> = {}) => ({
  id: 'raccoon-1',
  name: 'Ноти',
  level: 1,
  xp: 0,
  xpToNextLevel: 5,
  currentStreak: 0,
  badges: [],
  earnedBadgeCount: 0,
  ...overrides,
})

const session = { user: makeUser(), isAuthenticated: true }

beforeEach(() => {
  vi.clearAllMocks()
  window.localStorage.clear()
  profile.mockResolvedValue(makeProgress())
})

afterEach(() => {
  window.localStorage.clear()
})

describe('OnboardingCard', () => {
  it('lists the first steps for a fresh account', async () => {
    renderWithShell(<OnboardingCard />, { session })

    await waitFor(() => {
      expect(screen.getByTestId('onboarding-checklist')).toBeInTheDocument()
    })

    const items = screen.getAllByRole('listitem')
    expect(items).toHaveLength(3)
    expect(screen.getByText('Разместите объявление')).toBeInTheDocument()
    expect(screen.getByText('Возвращайтесь каждый день')).toBeInTheDocument()
    expect(screen.getAllByText('Шаг не выполнен')).toHaveLength(3)
  })

  it('marks steps as done once the pet has progress', async () => {
    profile.mockResolvedValue(
      makeProgress({ xp: 12, level: 2, currentStreak: 3 }),
    )

    renderWithShell(<OnboardingCard />, { session })

    await waitFor(() => {
      expect(screen.getAllByText('Шаг выполнен')).toHaveLength(3)
    })
  })

  it('never opens a hatching celebration dialog', async () => {
    renderWithShell(<OnboardingCard />, { session })

    await waitFor(() => {
      expect(screen.getByTestId('onboarding-card')).toBeInTheDocument()
    })

    expect(screen.queryByRole('dialog')).not.toBeInTheDocument()
  })

  it('can be dismissed and stays dismissed', async () => {
    renderWithShell(<OnboardingCard />, { session })

    await waitFor(() => {
      expect(screen.getByTestId('onboarding-card')).toBeInTheDocument()
    })

    await userEvent.click(screen.getByRole('button', { name: 'Закрыть' }))

    await waitFor(() => {
      expect(screen.queryByTestId('onboarding-card')).not.toBeInTheDocument()
    })
    expect(window.localStorage.getItem(onboardingStorageKey)).toContain(
      '"dismissed":true',
    )
  })

  it('keeps rendering in persistent mode even when dismissed', async () => {
    window.localStorage.setItem(
      onboardingStorageKey,
      JSON.stringify({ dismissed: true, tourSeen: true, tourRequested: false }),
    )

    renderWithShell(<OnboardingCard persistent />, { session })

    await waitFor(() => {
      expect(screen.getByTestId('onboarding-card')).toBeInTheDocument()
    })
  })
})
