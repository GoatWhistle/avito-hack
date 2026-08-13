import { act, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import {
  makeGameState,
  renderWithProviders,
} from '#/features/games/components/test-utils'
import { RaccoonJumpScreen } from './RaccoonJumpScreen'
import {
  SEED,
  captureFrames,
  dropPlayerBelowTheKillLine,
  liveGames,
  makeJumpReveal,
  makeJumpRound,
  nextTimestamp,
  seeds,
  setDevicePixelRatio,
  spyOnGameReset,
  stubContext,
} from './test-utils'

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
    claimReward: () => claimReward(),
  },
}))

const play = async () => {
  await userEvent.click(await screen.findByRole('button', { name: 'Играть' }))
  await waitFor(() => expect(liveGames.length).toBeGreaterThan(0))
}

const finishRun = async (frames: FrameRequestCallback[]) => {
  dropPlayerBelowTheKillLine()

  await act(async () => {
    frames.at(-1)?.(nextTimestamp())
  })
}

describe('RaccoonJumpScreen server round', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    stubContext()
    setDevicePixelRatio(1)

    state.mockResolvedValue(makeGameState({ slug: 'raccoonjump' }))
    startRound.mockResolvedValue(makeJumpRound())
    guess.mockResolvedValue(makeJumpReveal())

    spyOnGameReset()
  })

  afterEach(() => {
    vi.restoreAllMocks()
    vi.useRealTimers()
  })

  it('starts a server round before the run begins', async () => {
    renderWithProviders(<RaccoonJumpScreen />)

    await play()

    expect(startRound).toHaveBeenCalledWith('raccoonjump')
  })

  it('passes the server seed to the game', async () => {
    renderWithProviders(<RaccoonJumpScreen />)

    await play()

    expect(seeds.at(-1)).toBe(SEED)
    expect(liveGames.at(-1)?.seed).toBe(SEED)
  })

  it('does not start a run when the round request fails', async () => {
    startRound.mockRejectedValue(new Error('boom'))

    renderWithProviders(<RaccoonJumpScreen />)

    await userEvent.click(await screen.findByRole('button', { name: 'Играть' }))

    await screen.findByRole('alert')
    expect(liveGames).toHaveLength(0)
  })

  it('submits the final score and the collected indexes on game over', async () => {
    const frames = captureFrames()

    renderWithProviders(<RaccoonJumpScreen />)

    await play()

    const game = liveGames.at(-1)
    if (!game) throw new Error('no game instance was created')
    game.score = 42
    game.collected = [30]

    await finishRun(frames)

    await waitFor(() => expect(guess).toHaveBeenCalled())
    expect(guess).toHaveBeenCalledWith('raccoonjump', 'round1234567', {
      score: 42,
      collected: [30],
    })
  })

  it('shows the server best score instead of any local value', async () => {
    window.localStorage.setItem('raccoon-jump-best', '9999')
    state.mockResolvedValue(
      makeGameState({ slug: 'raccoonjump', best_score: 120 }),
    )

    renderWithProviders(<RaccoonJumpScreen />)

    expect(await screen.findByText('Рекорд: 120')).toBeInTheDocument()
    expect(screen.queryByText('Рекорд: 9999')).not.toBeInTheDocument()

    window.localStorage.removeItem('raccoon-jump-best')
  })

  it('never writes the best score to localStorage', async () => {
    const frames = captureFrames()
    const setItem = vi.spyOn(Storage.prototype, 'setItem')

    renderWithProviders(<RaccoonJumpScreen />)

    await play()
    setItem.mockClear()

    await finishRun(frames)

    expect(setItem).not.toHaveBeenCalledWith(
      'raccoon-jump-best',
      expect.anything(),
    )
  })

  it('renders the collected listings once the round is scored', async () => {
    guess.mockResolvedValue(
      makeJumpReveal({
        score: 80,
        best_score: 80,
        listings: [
          {
            display_id: 'aaa111222333',
            title: 'Велосипед Stels',
            price_kopeks: 1_250_000,
            photo_url: '/uploads/a/1.jpg',
          },
        ],
      }),
    )

    const frames = captureFrames()

    renderWithProviders(<RaccoonJumpScreen />)

    await play()
    await finishRun(frames)

    expect(await screen.findByText('Собранные объявления')).toBeInTheDocument()
    expect(await screen.findByText('Велосипед Stels')).toBeInTheDocument()
    expect(await screen.findByText('Новый рекорд!')).toBeInTheDocument()

    const link = screen.getByRole('link', { name: /Велосипед Stels/ })
    expect(link).toHaveAttribute('href', '/items/aaa111222333')
  })

  it('surfaces a failed score submission and retries it', async () => {
    guess.mockRejectedValueOnce(new Error('nope'))

    const frames = captureFrames()

    renderWithProviders(<RaccoonJumpScreen />)

    await play()
    await finishRun(frames)

    await screen.findByRole('alert')

    guess.mockResolvedValue(
      makeJumpReveal({
        score: 0,
        best_score: 7,
        new_best: false,
        counts_toward_streak: false,
      }),
    )

    await userEvent.click(screen.getByRole('button', { name: 'Повторить' }))

    await waitFor(() => expect(guess).toHaveBeenCalledTimes(2))
    expect(await screen.findByText('Рекорд: 7')).toBeInTheDocument()
  })

  it('collects a collectible once the score passes its index', async () => {
    renderWithProviders(<RaccoonJumpScreen />)

    await play()

    const game = liveGames.at(-1)
    if (!game) throw new Error('no game instance was created')

    expect(game.collectibleIndexes).toEqual([30, 70])
    expect(game.collected).toEqual([])

    game.player.y = game.player.y - 1000
    game.update()

    expect(game.score).toBeGreaterThanOrEqual(30)
    expect(game.collected).toContain(30)
  })

  it('does not submit a score after unmount', async () => {
    const frames = captureFrames()
    vi.spyOn(window, 'cancelAnimationFrame').mockImplementation(() => {})

    const { unmount } = renderWithProviders(<RaccoonJumpScreen />)

    await play()
    dropPlayerBelowTheKillLine()

    const pending = frames.at(-1)
    expect(pending).toBeTypeOf('function')

    unmount()
    guess.mockClear()

    await act(async () => {
      pending?.(nextTimestamp())
    })

    expect(guess).not.toHaveBeenCalled()
  })
})
