import { act, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import {
  makeGameState,
  renderWithProviders,
} from '#/features/games/components/test-utils'
import { RaccoonJumpScreen } from './RaccoonJumpScreen'
import { GAME_HEIGHT, GAME_WIDTH } from './game'
import {
  liveGames,
  makeJumpReveal,
  makeJumpRound,
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

describe('RaccoonJumpScreen', () => {
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

  it('resizes the canvas backing buffer when the device pixel ratio changes', async () => {
    const listeners: (() => void)[] = []
    vi.spyOn(window, 'matchMedia').mockImplementation(
      (media: string) =>
        ({
          media,
          matches: false,
          addEventListener: (_e: string, fn: () => void) => listeners.push(fn),
          removeEventListener: () => {},
        }) as unknown as MediaQueryList,
    )

    setDevicePixelRatio(2)

    const { container } = renderWithProviders(<RaccoonJumpScreen />)
    const canvas = container.querySelector('canvas')

    expect(canvas?.width).toBe(GAME_WIDTH * 2)
    expect(canvas?.height).toBe(GAME_HEIGHT * 2)

    setDevicePixelRatio(1)
    await act(async () => {
      listeners.forEach((listener) => listener())
    })

    expect(canvas?.width).toBe(GAME_WIDTH)
    expect(canvas?.height).toBe(GAME_HEIGHT)
  })

  it('keeps the run alive when the device pixel ratio changes', async () => {
    setDevicePixelRatio(2)

    const { container } = renderWithProviders(<RaccoonJumpScreen />)
    const canvas = container.querySelector('canvas')

    await play()

    const game = liveGames.at(-1)
    if (!game) throw new Error('no game instance was created')
    game.player.y = 123
    game.cameraY = -45

    setDevicePixelRatio(1)
    await act(async () => {
      window.dispatchEvent(new Event('resize'))
    })

    expect(canvas?.width).toBe(GAME_WIDTH)
    expect(game.status).toBe('playing')
    expect(game.player.y).toBe(123)
    expect(game.cameraY).toBe(-45)
  })

  it('gives the canvas an accessible name and fallback text', () => {
    const { container } = renderWithProviders(<RaccoonJumpScreen />)
    const canvas = container.querySelector('canvas')

    expect(canvas).toHaveAttribute('aria-label', 'Прыжки Ноти')
    expect(canvas).toHaveTextContent('Прыгайте по платформам как можно выше')
  })

  it('exposes a polite live region for the score', () => {
    const { container } = renderWithProviders(<RaccoonJumpScreen />)
    const status = container.querySelector('[role="status"]')

    expect(status).toHaveAttribute('aria-live', 'polite')
    expect(status).toHaveTextContent('')
  })
})
