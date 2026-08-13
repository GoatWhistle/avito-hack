import { act } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { renderWithProviders } from '#/features/games/components/test-utils'
import { DoodleJumpScreen } from './DoodleJumpScreen'
import { Game, GAME_HEIGHT, GAME_WIDTH, PLAYER_H } from './game'

const LS_KEY = 'doodle-jump-best'

const liveGames: Game[] = []

let clock = 0
const nextTimestamp = () => {
  clock += 32
  return performance.now() + clock
}

const stubContext = () => {
  const ctx = new Proxy(
    {},
    {
      get: (_target, prop) => {
        if (prop === 'canvas') return undefined
        if (prop === 'measureText') return () => ({ width: 10 })
        if (prop === 'createLinearGradient') {
          return () => ({ addColorStop: () => {} })
        }
        return () => {}
      },
      set: () => true,
    },
  ) as unknown as CanvasRenderingContext2D

  return vi
    .spyOn(HTMLCanvasElement.prototype, 'getContext')
    .mockReturnValue(ctx as never)
}

const dropPlayerBelowTheKillLine = () => {
  const game = liveGames.at(-1)
  if (!game) throw new Error('no game instance was created')

  game.player.y = game.cameraY + GAME_HEIGHT + PLAYER_H * 4
  game.player.vy = 1
}

const setDevicePixelRatio = (value: number) => {
  Object.defineProperty(window, 'devicePixelRatio', {
    configurable: true,
    value,
  })
}

describe('DoodleJumpScreen', () => {
  beforeEach(() => {
    stubContext()
    window.localStorage.removeItem(LS_KEY)
    setDevicePixelRatio(1)

    liveGames.length = 0
    const reset = Game.prototype.reset
    vi.spyOn(Game.prototype, 'reset').mockImplementation(function (
      this: Game,
      seed?: number,
    ) {
      liveGames.push(this)
      reset.call(this, seed)
    })
  })

  afterEach(() => {
    vi.restoreAllMocks()
    vi.useRealTimers()
  })

  it('does not write the best score to localStorage after unmount', () => {
    const frames: FrameRequestCallback[] = []
    vi.spyOn(window, 'requestAnimationFrame').mockImplementation((cb) => {
      frames.push(cb)
      return frames.length
    })
    vi.spyOn(window, 'cancelAnimationFrame').mockImplementation(() => {})

    const setItem = vi.spyOn(Storage.prototype, 'setItem')

    const { unmount } = renderWithProviders(<DoodleJumpScreen />)

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: ' ' }))
    })

    act(() => {
      frames.at(-1)?.(nextTimestamp())
    })

    dropPlayerBelowTheKillLine()

    const pending = frames.at(-1)
    expect(pending).toBeTypeOf('function')

    unmount()
    setItem.mockClear()

    act(() => {
      pending?.(nextTimestamp())
    })

    expect(setItem).not.toHaveBeenCalledWith(LS_KEY, expect.anything())
  })

  it('writes the best score when the run ends while still mounted', () => {
    const frames: FrameRequestCallback[] = []
    vi.spyOn(window, 'requestAnimationFrame').mockImplementation((cb) => {
      frames.push(cb)
      return frames.length
    })

    const setItem = vi.spyOn(Storage.prototype, 'setItem')

    renderWithProviders(<DoodleJumpScreen />)

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: ' ' }))
    })

    act(() => {
      frames.at(-1)?.(nextTimestamp())
    })

    dropPlayerBelowTheKillLine()
    setItem.mockClear()

    act(() => {
      frames.at(-1)?.(nextTimestamp())
    })

    expect(setItem).toHaveBeenCalledWith(LS_KEY, expect.anything())
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

    const { container } = renderWithProviders(<DoodleJumpScreen />)
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

    const { container } = renderWithProviders(<DoodleJumpScreen />)
    const canvas = container.querySelector('canvas')

    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: ' ' }))
    })

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
    const { container } = renderWithProviders(<DoodleJumpScreen />)
    const canvas = container.querySelector('canvas')

    expect(canvas).toHaveAttribute('aria-label', 'Прыжки Ноти')
    expect(canvas).toHaveTextContent('Прыгайте по платформам как можно выше')
  })

  it('exposes a polite live region for the score', () => {
    const { container } = renderWithProviders(<DoodleJumpScreen />)
    const status = container.querySelector('[role="status"]')

    expect(status).toHaveAttribute('aria-live', 'polite')
    expect(status).toHaveTextContent('')
  })

})
