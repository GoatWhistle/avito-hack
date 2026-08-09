import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { renderWithProviders } from '#/features/items/components/test-utils'
import type { LeaderboardEntry } from '#/features/leaderboard/types'
import { LeaderboardScreen } from './LeaderboardScreen'

const list = vi.fn()

vi.mock('#/features/leaderboard/repository', () => ({
  leaderboardRepository: {
    list: (query: unknown, signal?: AbortSignal) => list(query, signal),
  },
}))

const makeEntry = (
  overrides: Partial<LeaderboardEntry> = {},
): LeaderboardEntry => ({
  user_id: 'user-1',
  name: 'Анна Ковалёва',
  level: 15,
  xp: 581,
  streak_days: 41,
  rank: 1,
  ...overrides,
})

beforeEach(() => {
  vi.clearAllMocks()
  list.mockResolvedValue({
    items: [
      makeEntry(),
      makeEntry({ user_id: 'user-2', name: 'Борис Гурьев', rank: 2, xp: 312 }),
    ],
    myRank: 1,
    nextCursor: '',
  })
})

describe('LeaderboardScreen', () => {
  it('invites a guest to sign in instead of showing an error', async () => {
    renderWithProviders(<LeaderboardScreen />, { userId: null })

    expect(
      await screen.findByText('Рейтинг открыт участникам'),
    ).toBeInTheDocument()
    expect(screen.getByRole('link', { name: 'Войти' })).toHaveAttribute(
      'href',
      '/sign-in',
    )
    expect(screen.getByRole('link', { name: 'Начать' })).toHaveAttribute(
      'href',
      '/sign-up',
    )
    expect(screen.queryByRole('alert')).not.toBeInTheDocument()
    expect(list).not.toHaveBeenCalled()
  })

  it('renders the ranking table for a signed-in user', async () => {
    renderWithProviders(<LeaderboardScreen />)

    const rows = await screen.findAllByTestId('leaderboard-row')

    expect(rows.length).toBeGreaterThanOrEqual(2)
    expect(screen.getAllByText('Анна Ковалёва').length).toBeGreaterThan(0)
    expect(screen.getByTestId('my-rank')).toHaveTextContent('Ваше место: 1')
  })

  it('marks the current user row', async () => {
    renderWithProviders(<LeaderboardScreen />)

    const rows = await screen.findAllByTestId('leaderboard-row')
    const mine = rows.filter((row) => row.dataset.me === 'true')

    expect(mine.length).toBeGreaterThan(0)
    expect(mine[0]).toHaveAttribute('aria-current', 'true')
  })

  it('keeps long names readable without truncating them', async () => {
    const longName = 'Александра Константинопольская-Вишневецкая'
    list.mockResolvedValue({
      items: [makeEntry({ name: longName })],
      myRank: 1,
      nextCursor: '',
    })
    renderWithProviders(<LeaderboardScreen />)

    const name = (await screen.findAllByText(longName))[0]

    expect(name.className).toContain('break-words')
    expect(name.className).not.toContain('truncate')
  })

  it('shows an empty state when nobody is ranked', async () => {
    list.mockResolvedValue({ items: [], myRank: null, nextCursor: '' })
    renderWithProviders(<LeaderboardScreen />)

    expect(await screen.findByText('Рейтинг пока пуст')).toBeInTheDocument()
  })

  it('shows an error with a retry action', async () => {
    list.mockRejectedValue(new Error('offline'))
    renderWithProviders(<LeaderboardScreen />)

    const alerts = await screen.findAllByRole('alert')

    expect(alerts[0]).toHaveTextContent('Не удалось загрузить рейтинг')

    list.mockResolvedValue({
      items: [makeEntry()],
      myRank: 1,
      nextCursor: '',
    })
    await userEvent.click(
      screen.getAllByRole('button', { name: 'Повторить' })[0],
    )

    expect(
      (await screen.findAllByTestId('leaderboard-row')).length,
    ).toBeGreaterThan(0)
  })
})
