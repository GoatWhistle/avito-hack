import { vi } from 'vitest'
import { Game, GAME_HEIGHT, PLAYER_H } from './game'

export const SEED = 987654321

export const makeJumpPrompt = (overrides: Record<string, unknown> = {}) => ({
  seed: SEED,
  max_score: 100000,
  min_streak_score: 25,
  max_score_per_second: 189,
  collectibles: [
    { index: 30, title: 'Велосипед Stels', photo_url: '/uploads/a/1.jpg' },
    { index: 70, title: 'Диван угловой', photo_url: '/uploads/b/1.jpg' },
  ],
  best_score: 0,
  started_at_unix_milli: 1_700_000_000_000,
  ...overrides,
})

export const makeJumpRound = (prompt: unknown = makeJumpPrompt()) => ({
  round_id: 'round1234567',
  streak: 0,
  target_streak: 0,
  state: 'active',
  prompt,
})

export const makeJumpReveal = (overrides: Record<string, unknown> = {}) => ({
  correct: true,
  progress: 'win',
  reveal: {
    score: 42,
    best_score: 42,
    new_best: true,
    counts_toward_streak: true,
    listings: [],
    ...overrides,
  },
  streak: 0,
  state: 'won',
  attempt_completed: true,
})

export const liveGames: Game[] = []
export const seeds: (number | undefined)[] = []

let clock = 0
export const nextTimestamp = () => {
  clock += 32

  return performance.now() + clock
}

export const stubContext = () => {
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

export const spyOnGameReset = () => {
  liveGames.length = 0
  seeds.length = 0

  const reset = Game.prototype.reset
  vi.spyOn(Game.prototype, 'reset').mockImplementation(function (
    this: Game,
    seed?: number,
  ) {
    liveGames.push(this)
    seeds.push(seed)
    reset.call(this, seed)
  })
}

export const dropPlayerBelowTheKillLine = () => {
  const game = liveGames.at(-1)
  if (!game) throw new Error('no game instance was created')

  game.player.y = game.cameraY + GAME_HEIGHT + PLAYER_H * 4
  game.player.vy = 1
}

export const setDevicePixelRatio = (value: number) => {
  Object.defineProperty(window, 'devicePixelRatio', {
    configurable: true,
    value,
  })
}

export const captureFrames = () => {
  const frames: FrameRequestCallback[] = []
  vi.spyOn(window, 'requestAnimationFrame').mockImplementation((cb) => {
    frames.push(cb)

    return frames.length
  })

  return frames
}
