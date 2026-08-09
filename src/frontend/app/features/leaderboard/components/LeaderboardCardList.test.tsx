import { screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { renderWithProviders } from '#/features/items/components/test-utils'
import { initialsOf, medalOf } from '#/features/leaderboard/lib'
import type { LeaderboardEntry } from '#/features/leaderboard/types'
import { LeaderboardCardList } from './LeaderboardCardList'
import { MyRankCard } from './MyRankCard'

const makeEntry = (
  overrides: Partial<LeaderboardEntry> = {},
): LeaderboardEntry => ({
  user_id: 'user-1',
  name: 'Анна Ковалёва',
  level: 15,
  xp: 455,
  streak_days: 41,
  rank: 1,
  ...overrides,
})

const entries = [
  makeEntry(),
  makeEntry({ user_id: 'user-2', name: 'Борис Гурьев', rank: 2, xp: 312 }),
  makeEntry({ user_id: 'user-3', name: 'Вера Лис', rank: 3, xp: 280 }),
  makeEntry({ user_id: 'user-4', name: 'Глеб Орлов', rank: 12, xp: 90 }),
]

describe('rank helpers', () => {
  it('maps the top three ranks to medals', () => {
    expect(medalOf(1)).toBe('gold')
    expect(medalOf(2)).toBe('silver')
    expect(medalOf(3)).toBe('bronze')
    expect(medalOf(4)).toBeNull()
  })

  it('builds two-letter initials from cyrillic names', () => {
    expect(initialsOf('Анна Ковалёва')).toBe('АК')
    expect(initialsOf('Глеб')).toBe('Г')
    expect(initialsOf('  Вера   Лис ')).toBe('ВЛ')
  })
})

describe('LeaderboardCardList', () => {
  it('renders one card per entry', async () => {
    renderWithProviders(
      <LeaderboardCardList entries={entries} myRank={12} label="Рейтинг" />,
    )

    expect(await screen.findAllByTestId('leaderboard-card')).toHaveLength(4)
  })

  it('shows medal icons for the podium and a number elsewhere', async () => {
    renderWithProviders(
      <LeaderboardCardList entries={entries} myRank={12} label="Рейтинг" />,
    )

    const cards = await screen.findAllByTestId('leaderboard-card')

    expect(cards[0].querySelector('svg')).toBeInTheDocument()
    expect(cards[1].querySelector('svg')).toBeInTheDocument()
    expect(cards[2].querySelector('svg')).toBeInTheDocument()
    expect(cards[3].querySelector('svg')).not.toBeInTheDocument()
    expect(cards[3]).toHaveTextContent('12')
  })

  it('renders initials and the stat line for an entry', async () => {
    renderWithProviders(
      <LeaderboardCardList entries={entries} myRank={12} label="Рейтинг" />,
    )

    const cards = await screen.findAllByTestId('leaderboard-card')

    expect(cards[0]).toHaveTextContent('АК')
    expect(cards[0]).toHaveTextContent('15 ур. · 455 XP · 41 дн.')
  })

  it('highlights the current user row', async () => {
    renderWithProviders(
      <LeaderboardCardList entries={entries} myRank={12} label="Рейтинг" />,
    )

    const cards = await screen.findAllByTestId('leaderboard-card')
    const mine = cards.filter((card) => card.dataset.me === 'true')

    expect(mine).toHaveLength(1)
    expect(mine[0]).toHaveAttribute('aria-current', 'true')
    expect(mine[0]).toHaveTextContent('Глеб Орлов')
  })
})

describe('MyRankCard', () => {
  it('shows the rank with a hash prefix', async () => {
    renderWithProviders(<MyRankCard rank={12} />)

    const card = await screen.findByTestId('my-rank-card')

    expect(card).toHaveTextContent('Ваше место')
    expect(card).toHaveTextContent('#12')
  })

  it('falls back to the not-ranked hint', async () => {
    renderWithProviders(<MyRankCard rank={null} />)

    expect(await screen.findByTestId('my-rank-card')).toHaveTextContent(
      'Вы ещё не в рейтинге',
    )
  })
})
