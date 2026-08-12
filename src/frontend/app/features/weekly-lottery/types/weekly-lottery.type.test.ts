import { describe, expect, it } from 'vitest'
import {
  lotteryRevealSchema,
  lotteryRunSchema,
  lotteryStateSchema,
} from './weekly-lottery.type'

describe('lotteryRunSchema', () => {
  it('rejects a symbol leaked into a closed slot', () => {
    const result = lotteryRunSchema.safeParse({
      id: 'run123456789',
      state: 'active',
      slots: Array.from({ length: 9 }, (_, index) => ({
        index,
        opened: false,
        ...(index === 0 ? { symbol: 'bicycle' } : {}),
      })),
      created_at: '2026-08-12T12:00:00Z',
    })

    expect(result.success).toBe(false)
  })

  it('rejects duplicate slot indexes', () => {
    const result = lotteryRunSchema.safeParse({
      id: 'run123456789',
      state: 'active',
      slots: Array.from({ length: 9 }, () => ({ index: 0, opened: false })),
      created_at: '2026-08-12T12:00:00Z',
    })

    expect(result.success).toBe(false)
  })
})

describe('weekly lottery API consistency', () => {
  const slots = Array.from({ length: 9 }, (_, index) => ({
    index,
    opened: false,
  }))
  const run = {
    id: 'run123456789',
    state: 'active',
    slots,
    created_at: '2026-08-12T12:00:00Z',
  }

  it('rejects availability that contradicts the run state', () => {
    const result = lotteryStateSchema.safeParse({
      available: false,
      week_start: '2026-08-10T00:00:00+03:00',
      next_available_at: '2026-08-17T00:00:00+03:00',
      run,
    })

    expect(result.success).toBe(false)
  })

  it('rejects a reveal that is missing from the returned run', () => {
    const result = lotteryRevealSchema.safeParse({
      index: 4,
      symbol: 'bicycle',
      run,
    })

    expect(result.success).toBe(false)
  })
})
