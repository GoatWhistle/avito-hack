import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { MoreLessGame } from './MoreLessGame'
import { makeGameState, makePrompt, renderWithProviders } from '../test-utils'

const state = vi.fn()
const startRound = vi.fn()
const guess = vi.fn()
const claimReward = vi.fn()
vi.mock('#/features/games/repository', () => ({
  gameRepository: {
    state: (s: string) => state(s),
    startRound: (s: string) => startRound(s),
    guess: (s: string, r: string, m: unknown) => guess(s, r, m),
    claimReward: (s: string) => claimReward(s),
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

describe('MoreLessGame round lifecycle', () => {
  it('shows the new round items after replay, not the frozen ones', async () => {
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

    startRound.mockResolvedValue({
      round_id: 'round7654321',
      streak: 0,
      target_streak: 7,
      state: 'active',
      prompt: makePrompt({
        left: {
          item_id: 'ccc1',
          title: 'Ноутбук Lenovo',
          photo_url: '/x.jpg',
          price: 5000,
        },
        right: {
          item_id: 'ddd2',
          title: 'Телефон Xiaomi',
          photo_url: '/y.jpg',
        },
      }),
    })

    await userEvent.click(
      await screen.findByRole('button', { name: /Играть ещё раз/ }),
    )
    await waitFor(() => expect(startRound).toHaveBeenCalledTimes(2))

    expect(await screen.findByText('Ноутбук Lenovo')).toBeInTheDocument()
    expect(screen.queryByText('Велосипед Stels')).not.toBeInTheDocument()
  })
})
