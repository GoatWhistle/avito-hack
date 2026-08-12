import type { AxiosInstance } from 'axios'
import { describe, expect, it, vi } from 'vitest'
import { WeeklyLotteryRepository } from './weekly-lottery.repository'

const closedSlots = Array.from({ length: 9 }, (_, index) => ({
  index,
  opened: false,
}))

const run = {
  id: 'run123456789',
  state: 'active',
  slots: closedSlots,
  created_at: '2026-08-12T12:00:00Z',
}

const state = {
  available: true,
  week_start: '2026-08-10T00:00:00+03:00',
  next_available_at: '2026-08-17T00:00:00+03:00',
  run,
}

const revealedRun = {
  ...run,
  slots: closedSlots.map((slot) =>
    slot.index === 4
      ? { index: 4, opened: true as const, symbol: 'bicycle' }
      : slot,
  ),
}

const makeClient = () => {
  const get = vi.fn((url: string) =>
    Promise.resolve({ data: url.endsWith('/state') ? state : [] }),
  )
  const post = vi.fn((url: string) =>
    Promise.resolve({
      data: url.endsWith('/runs')
        ? run
        : { index: 4, symbol: 'bicycle', run: revealedRun },
    }),
  )

  return {
    client: { get, post } as unknown as AxiosInstance,
    get,
    post,
  }
}

describe('WeeklyLotteryRepository contract', () => {
  it('uses the backend state and prize endpoints', async () => {
    const { client, get } = makeClient()
    const repository = new WeeklyLotteryRepository(client)

    await repository.state()
    await repository.prizes()

    expect(get).toHaveBeenNthCalledWith(1, 'weekly-lottery/state', {
      signal: undefined,
    })
    expect(get).toHaveBeenNthCalledWith(2, 'weekly-lottery/prizes', {
      signal: undefined,
    })
  })

  it('uses the backend start and zero-based reveal endpoints', async () => {
    const { client, post } = makeClient()
    const repository = new WeeklyLotteryRepository(client)

    await repository.start()
    const revealed = await repository.reveal('run123456789', 4)

    expect(post).toHaveBeenNthCalledWith(1, 'weekly-lottery/runs')
    expect(post).toHaveBeenNthCalledWith(
      2,
      'weekly-lottery/runs/run123456789/slots/4/reveal',
    )
    expect(revealed.index).toBe(4)
    expect(revealed.symbol).toBe('bicycle')
  })

  it('accepts the state shape returned by the backend', async () => {
    const { client } = makeClient()

    await expect(
      new WeeklyLotteryRepository(client).state(),
    ).resolves.toEqual(state)
  })
})
