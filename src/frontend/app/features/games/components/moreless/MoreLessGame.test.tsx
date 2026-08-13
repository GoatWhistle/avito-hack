import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { MoreLessGame } from './MoreLessGame'
import {
  makeGameState,
  makePrompt,
  makeStreak,
  renderWithProviders,
} from '../test-utils'

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

beforeEach(() => {
  vi.clearAllMocks()
  window.localStorage.removeItem('avito-hack.moreless.help-dismissed')
  state.mockResolvedValue(makeGameState())
  startRound.mockResolvedValue({
    round_id: 'round1234567',
    streak: 0,
    target_streak: 7,
    state: 'active',
    prompt: makePrompt(),
  })
})

describe('MoreLessGame', () => {
  it('starts a round and shows both listings', async () => {
    renderWithProviders(<MoreLessGame />)


    expect(await screen.findByText('Велосипед Stels')).toBeInTheDocument()
    expect(screen.getByText('Диван угловой')).toBeInTheDocument()
  })

  it('never renders the hidden price before a guess', async () => {
    renderWithProviders(<MoreLessGame />)

    await screen.findByText('Диван угловой')

    expect(screen.getByText(/12 500/)).toBeInTheDocument()
    expect(screen.queryByText(/640/)).not.toBeInTheDocument()
  })

  it('sends the chosen direction and reports a correct guess', async () => {
    guess.mockResolvedValue({
      correct: true,
      reveal: { right_price: 2_000_000 },
      streak: 1,
      state: 'active',
      prompt: makePrompt(),
      attempt_completed: false,
    })

    renderWithProviders(<MoreLessGame />)
    await userEvent.click(await screen.findByRole('button', { name: 'Дороже' }))

    await waitFor(() => {
      expect(guess).toHaveBeenCalledWith('moreless', 'round1234567', {
        choice: 'higher',
      })
    })
    expect(await screen.findByText('Верно!')).toBeInTheDocument()
  })

  it('ends the round on a wrong guess', async () => {
    guess.mockResolvedValue({
      correct: false,
      reveal: { right_price: 100 },
      streak: 0,
      state: 'lost',
      attempt_completed: false,
    })

    renderWithProviders(<MoreLessGame />)
    await userEvent.click(
      await screen.findByRole('button', { name: 'Дешевле' }),
    )

    expect(await screen.findByText('Игра окончена')).toBeInTheDocument()
  })

  it('keeps the streak card out of the game board', async () => {
    state.mockResolvedValue(
      makeGameState({
        streak: makeStreak({ current_days: 7, reward_ready: true }),
      }),
    )

    renderWithProviders(<MoreLessGame />)

    expect(screen.queryByText('Недельный стрик')).not.toBeInTheDocument()
    expect(
      screen.queryByRole('button', { name: 'Забрать награду' }),
    ).not.toBeInTheDocument()
  })

  it('reveals both prices and links after a guess', async () => {
    guess.mockResolvedValue({
      correct: true,
      reveal: { right_price: 2_000_000 },
      streak: 1,
      state: 'active',
      prompt: makePrompt(),
      attempt_completed: false,
    })

    renderWithProviders(<MoreLessGame />)
    await userEvent.click(await screen.findByRole('button', { name: 'Дороже' }))

    const links = await screen.findAllByRole('link', {
      name: /Смотреть объявление/,
    })

    expect(links).toHaveLength(2)
    expect(links[0]).toHaveAttribute('href', '/items/aaa111222333')
    expect(links[1]).toHaveAttribute('href', '/items/bbb444555666')
    expect(
      await screen.findByRole('button', { name: 'Дальше' }),
    ).toBeInTheDocument()
  })

  it('does not repeat the rules inside the board', async () => {
    renderWithProviders(<MoreLessGame />)

    await screen.findByText('Диван угловой')

    expect(
      screen.queryByText(
        'Слева цена известна. Угадайте, дороже или дешевле объявление справа.',
      ),
    ).not.toBeInTheDocument()
    expect(screen.queryByText('Слева цена уже известна')).not.toBeInTheDocument()
  })

  it('shows the price when the round is lost', async () => {
    guess.mockResolvedValue({
      correct: false,
      reveal: { right_price: 990_000 },
      streak: 0,
      state: 'lost',
      attempt_completed: false,
    })

    renderWithProviders(<MoreLessGame />)
    await userEvent.click(
      await screen.findByRole('button', { name: 'Дешевле' }),
    )

    expect(await screen.findByText('Игра окончена')).toBeInTheDocument()
    expect(screen.getByText(/9 900/)).toBeInTheDocument()
  })
})
