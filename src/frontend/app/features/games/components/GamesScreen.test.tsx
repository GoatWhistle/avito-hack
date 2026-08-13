import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { GamesScreen } from './GamesScreen'
import {
  makeGameList,
  makeStreak,
  makeSummary,
  renderWithProviders,
} from './test-utils'

const list = vi.fn()
const claimReward = vi.fn()

vi.mock('#/features/games/repository', () => ({
  gameRepository: {
    list: () => list(),
    claimReward: () => claimReward(),
  },
}))

beforeEach(() => {
  vi.clearAllMocks()
  list.mockResolvedValue(makeGameList())
})

describe('GamesScreen', () => {
  it('lists a registered game with a link to play it', async () => {
    renderWithProviders(<GamesScreen />)

    const link = await screen.findByRole('link', { name: /Больше или меньше/ })

    expect(link).toHaveAttribute('href', '/play/higher-lower')
  })

  it('shows one streak card for every mini-game', async () => {
    list.mockResolvedValue(
      makeGameList({
        games: [makeSummary(), makeSummary({ slug: 'bukovki' })],
        streak: makeStreak({ current_days: 4 }),
      }),
    )

    renderWithProviders(<GamesScreen />)

    expect(await screen.findAllByText('Недельный стрик')).toHaveLength(1)
    expect(screen.getByText('4 из 7 дней')).toBeInTheDocument()
  })

  it('claims the global reward once the streak is full', async () => {
    list.mockResolvedValue(
      makeGameList({
        streak: makeStreak({ current_days: 7, reward_ready: true }),
      }),
    )
    claimReward.mockResolvedValue({ code: 'PROMO-ABCD1234' })

    renderWithProviders(<GamesScreen />)

    await userEvent.click(
      await screen.findByRole('button', { name: 'Забрать награду' }),
    )

    expect(await screen.findByText('PROMO-ABCD1234')).toBeInTheDocument()
    expect(claimReward).toHaveBeenCalledWith()
  })

  it('marks only the games that were played today', async () => {
    list.mockResolvedValue(
      makeGameList({
        games: [
          makeSummary({ daily_done: true }),
          makeSummary({ slug: 'bukovki' }),
        ],
        daily_done: true,
      }),
    )

    renderWithProviders(<GamesScreen />)

    expect(await screen.findByText(/Сегодня засчитано/)).toBeInTheDocument()
  })

  it('hides games that have no frontend implementation', async () => {
    list.mockResolvedValue(
      makeGameList({ games: [makeSummary({ slug: 'not-implemented-yet' })] }),
    )

    renderWithProviders(<GamesScreen />)

    expect(await screen.findByText('Игр пока нет')).toBeInTheDocument()
  })
})
