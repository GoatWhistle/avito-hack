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

beforeEach(() => {
  vi.clearAllMocks()
  window.localStorage.removeItem('avito-hack.bukovki.help-dismissed')
  state.mockResolvedValue(makeGameState({ slug: 'bukovki' }))
  startRound.mockResolvedValue({
    round_id: 'round1234567',
    streak: 0,
    target_streak: 0,
    state: 'active',
    prompt,
  })
})

describe('BukovkiGame', () => {
  it('renders an on-screen alphabet instead of a text field', async () => {
    renderWithProviders(<BukovkiGame />)

    await screen.findByRole('button', { name: 'д' })

    expect(screen.queryByRole('textbox')).not.toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Удалить' })).toBeInTheDocument()
  })

  it('types into the cells and submits the word', async () => {
    guess.mockResolvedValue({
      correct: false,
      reveal: { feedback: [], game_over: false, win: false },
      streak: 0,
      state: 'active',
      prompt: { ...prompt, history: [[{ char: 'д', status: 'absent' }]] },
      attempt_completed: false,
    })

    renderWithProviders(<BukovkiGame />)

    for (const char of 'диван') {
      await userEvent.click(await screen.findByRole('button', { name: char }))
    }
    await userEvent.click(screen.getByRole('button', { name: 'Проверить' }))

    await waitFor(() => {
      expect(guess).toHaveBeenCalledWith('bukovki', 'round1234567', {
        guess: 'диван',
      })
    })
  })

  it('marks keyboard letters with their status', async () => {
    startRound.mockResolvedValue({
      round_id: 'round1234567',
      streak: 0,
      target_streak: 0,
      state: 'active',
      prompt: {
        ...prompt,
        history: [
          [
            { char: 'д', status: 'correct' },
            { char: 'о', status: 'absent' },
          ],
        ],
      },
    })

    renderWithProviders(<BukovkiGame />)

    const correct = await screen.findByRole('button', { name: 'д' })
    const absent = screen.getByRole('button', { name: 'о' })

    expect(correct.className).toContain('bg-success')
    expect(absent.className).toContain('bg-destructive')
  })

  it('shows the secret and matching listings once the round ends', async () => {
    guess.mockResolvedValue({
      correct: true,
      reveal: {
        feedback: [],
        game_over: true,
        win: true,
        secret: 'диван',
        listings: [
          {
            display_id: 'aaa111222333',
            title: 'Диван угловой, велюр',
            price_kopeks: 4_500_000,
            photo_url: '/uploads/aaa111222333/demo-1.jpg',
          },
        ],
      },
      streak: 0,
      state: 'won',
      attempt_completed: true,
    })

    renderWithProviders(<BukovkiGame />)
    for (const char of 'диван') {
      await userEvent.click(await screen.findByRole('button', { name: char }))
    }
    await userEvent.click(screen.getByRole('button', { name: 'Проверить' }))

    expect(await screen.findByText(/ДИВАН/)).toBeInTheDocument()

    const link = await screen.findByRole('link', {
      name: /Диван угловой/,
    })
    expect(link).toHaveAttribute('href', '/items/aaa111222333')
  })

  it('does not repeat the rules inside the board', async () => {
    renderWithProviders(<BukovkiGame />)

    await screen.findByRole('button', { name: 'д' })

    expect(screen.queryByText('Буква на своём месте')).not.toBeInTheDocument()
  })

  it('supports the physical keyboard', async () => {
    guess.mockResolvedValue({
      correct: false,
      reveal: { feedback: [], game_over: false, win: false },
      streak: 0,
      state: 'active',
      prompt,
      attempt_completed: false,
    })

    renderWithProviders(<BukovkiGame />)
    await screen.findByRole('button', { name: 'д' })

    await userEvent.keyboard('sofa')
    await userEvent.keyboard('диванн')
    await userEvent.keyboard('{Backspace}л')
    await userEvent.keyboard('{Enter}')

    await waitFor(() => {
      expect(guess).toHaveBeenCalledWith('bukovki', 'round1234567', {
        guess: 'дивал',
      })
    })
  })
})
