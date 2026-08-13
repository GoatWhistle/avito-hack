import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { BukovkiGame } from './BukovkiGame'
import { makeGameState, renderWithProviders } from '../test-utils'

const state = vi.fn()
const startRound = vi.fn()
const guess = vi.fn()
const claimReward = vi.fn()

vi.mock('#/features/games/repository', () => ({
  gameRepository: {
    state: (slug: string) => state(slug),
    startRound: (slug: string) => startRound(slug),
    guess: (slug: string, roundId: string, move: unknown) =>
      guess(slug, roundId, move),
    claimReward: (slug: string) => claimReward(slug),
  },
}))

const prompt = { word_length: 5, max_tries: 6, history: [] }

const staleState = () =>
  makeGameState({
    slug: 'bukovki',
    active_round: {
      round_id: 'stale0000001',
      streak: 0,
      prompt,
      attempts_used: 0,
      max_attempts: 6,
    },
  })

const win = {
  correct: true,
  reveal: { feedback: [], game_over: true, win: true, secret: 'диван' },
  streak: 1,
  state: 'won',
  attempt_completed: true,
}

const playRound = async (word: string) => {
  for (const char of word) {
    await userEvent.click(await screen.findByRole('button', { name: char }))
  }
  await userEvent.click(screen.getByRole('button', { name: 'Проверить' }))
}

const replay = async () => {
  await userEvent.click(
    await screen.findByRole('button', { name: /Играть ещё раз/ }),
  )
  await waitFor(() => {
    expect(startRound).toHaveBeenCalledWith('bukovki')
  })
}

beforeEach(() => {
  vi.clearAllMocks()
  window.localStorage.removeItem('avito-hack.bukovki.help-dismissed')
  state.mockResolvedValue(staleState())
  guess.mockResolvedValue(win)
})

describe('BukovkiGame round lifecycle', () => {
  it('sends the next guess to the freshly started round', async () => {
    startRound.mockResolvedValue({
      round_id: 'fresh0000002',
      streak: 0,
      target_streak: 0,
      state: 'active',
      prompt,
    })

    renderWithProviders(<BukovkiGame />)

    await playRound('диван')
    await waitFor(() => {
      expect(guess).toHaveBeenCalledWith('bukovki', 'stale0000001', {
        guess: 'диван',
      })
    })

    await replay()
    await playRound('диван')

    await waitFor(() => {
      expect(guess).toHaveBeenLastCalledWith('bukovki', 'fresh0000002', {
        guess: 'диван',
      })
    })
  })

  it('falls back to a browse link when the round ends with no listings', async () => {
    renderWithProviders(<BukovkiGame />)

    await playRound('диван')

    expect(
      await screen.findByText('Подходящие объявления сейчас недоступны'),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('link', { name: 'Смотреть все объявления' }),
    ).toHaveAttribute('href', '/items')
  })

  it('rebuilds the board for the new round instead of reusing the old word', async () => {
    startRound.mockResolvedValue({
      round_id: 'fresh0000002',
      streak: 0,
      target_streak: 0,
      state: 'active',
      prompt: { word_length: 4, max_tries: 6, history: [] },
    })

    renderWithProviders(<BukovkiGame />)

    await playRound('диван')
    await replay()
    await playRound('дива')

    await waitFor(() => {
      expect(guess).toHaveBeenLastCalledWith('bukovki', 'fresh0000002', {
        guess: 'дива',
      })
    })
  })

  it('does not re-adopt the finished round when the state query refetches', async () => {
    startRound.mockResolvedValue({
      round_id: 'fresh0000002',
      streak: 0,
      target_streak: 0,
      state: 'active',
      prompt,
    })

    renderWithProviders(<BukovkiGame />)

    await playRound('диван')
    await replay()

    expect(
      screen.queryByRole('button', { name: /Играть ещё раз/ }),
    ).not.toBeInTheDocument()

    await playRound('диван')

    await waitFor(() => {
      expect(guess).toHaveBeenLastCalledWith('bukovki', 'fresh0000002', {
        guess: 'диван',
      })
    })
  })
})
