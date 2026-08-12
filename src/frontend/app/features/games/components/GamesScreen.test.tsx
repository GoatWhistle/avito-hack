import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { GamesScreen } from './GamesScreen'
import { makeStreak, makeSummary, renderWithProviders } from './test-utils'

const list = vi.fn()
const claimReward = vi.fn()

vi.mock('#/features/games/repository', () => ({
  gameRepository: {
    list: () => list(),
    claimReward: (slug: string) => claimReward(slug),
  },
}))

beforeEach(() => {
  vi.clearAllMocks()
  list.mockResolvedValue([makeSummary()])
})

describe('GamesScreen', () => {
  it('lists a registered game with a link to play it', async () => {
    renderWithProviders(<GamesScreen />)

    const link = await screen.findByRole('link', { name: /Больше или меньше/ })

    expect(link).toHaveAttribute('href', '/play/moreless')
  })

  it('shows the current streak', async () => {
    list.mockResolvedValue([
      makeSummary({ streak: makeStreak({ current_days: 4 }) }),
    ])

    renderWithProviders(<GamesScreen />)

    expect(await screen.findByText('4')).toBeInTheDocument()
  })

  it('shows the streak card and claims the reward', async () => {
    list.mockResolvedValue([
      makeSummary({
        streak: makeStreak({ current_days: 7, reward_ready: true }),
      }),
    ])
    claimReward.mockResolvedValue({ code: 'PROMO-ABCD1234' })

    renderWithProviders(<GamesScreen />)

    expect(await screen.findByText('Недельный стрик')).toBeInTheDocument()
    await userEvent.click(
      await screen.findByRole('button', { name: 'Забрать награду' }),
    )

    expect(await screen.findByText('PROMO-ABCD1234')).toBeInTheDocument()
  })

  it('hides games that have no frontend implementation', async () => {
    list.mockResolvedValue([makeSummary({ slug: 'not-implemented-yet' })])

    renderWithProviders(<GamesScreen />)

    expect(await screen.findByText('Игр пока нет')).toBeInTheDocument()
  })
})
