import { act, renderHook } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { usePetCelebration } from './usePetCelebration'
import { makePet } from '#/features/pet/components/test-utils'
import type { PetEvent } from '#/features/pet/types'

beforeEach(() => {
  vi.useFakeTimers()
})

afterEach(() => {
  vi.useRealTimers()
})

const send = (
  result: { current: ReturnType<typeof usePetCelebration> },
  event: PetEvent,
) => act(() => result.current.handleEvent(event))

describe('usePetCelebration', () => {
  it('starts with nothing to show', () => {
    const { result } = renderHook(() => usePetCelebration())

    expect(result.current.emotion).toBeNull()
    expect(result.current.banner).toBeNull()
    expect(result.current.xpToasts).toHaveLength(0)
  })

  it('queues an xp toast and celebrates', () => {
    const { result } = renderHook(() => usePetCelebration())

    send(result, { type: 'xp.gained', payload: { amount: 12 } })

    expect(result.current.emotion).toBe('celebrate')
    expect(result.current.xpToasts).toEqual([{ id: 1, amount: 12 }])
  })

  it('expires an xp toast after the timeout', () => {
    const { result } = renderHook(() => usePetCelebration())

    send(result, { type: 'xp.gained', payload: { amount: 5 } })
    act(() => {
      vi.advanceTimersByTime(1_800)
    })

    expect(result.current.xpToasts).toHaveLength(0)
  })

  it('stacks concurrent toasts with unique ids', () => {
    const { result } = renderHook(() => usePetCelebration())

    send(result, { type: 'xp.gained', payload: { amount: 1 } })
    send(result, { type: 'xp.gained', payload: { amount: 2 } })

    expect(result.current.xpToasts.map((toast) => toast.id)).toEqual([1, 2])
  })

  it('raises a level up banner', () => {
    const { result } = renderHook(() => usePetCelebration())

    send(result, { type: 'level.up', payload: { level: 6 } })

    expect(result.current.emotion).toBe('levelup')
    expect(result.current.banner).toEqual({ kind: 'levelUp', level: 6 })
  })

  it('raises a hatched banner', () => {
    const { result } = renderHook(() => usePetCelebration())

    send(result, { type: 'pet.hatched', payload: makePet() })

    expect(result.current.emotion).toBe('hatching')
    expect(result.current.banner).toEqual({ kind: 'hatched' })
  })

  it('raises a reward banner', () => {
    const { result } = renderHook(() => usePetCelebration())

    send(result, {
      type: 'reward.granted',
      payload: { reward_id: 'r1', title: 'Промокод' },
    })

    expect(result.current.banner).toEqual({ kind: 'reward', title: 'Промокод' })
  })

  it('only banners a streak milestone', () => {
    const { result } = renderHook(() => usePetCelebration())

    send(result, { type: 'streak.updated', payload: { days: 3 } })
    expect(result.current.banner).toBeNull()

    send(result, {
      type: 'streak.updated',
      payload: { days: 7, milestone: true },
    })
    expect(result.current.banner).toEqual({ kind: 'streak', days: 7 })
  })

  it('ignores events without a celebration', () => {
    const { result } = renderHook(() => usePetCelebration())

    send(result, { type: 'pet.updated', payload: makePet() })

    expect(result.current.emotion).toBeNull()
    expect(result.current.banner).toBeNull()
  })

  it('clears the emotion on demand', () => {
    const { result } = renderHook(() => usePetCelebration())

    send(result, { type: 'level.up', payload: { level: 2 } })
    act(() => result.current.clearEmotion())

    expect(result.current.emotion).toBeNull()
  })

  it('dismisses the banner on demand', () => {
    const { result } = renderHook(() => usePetCelebration())

    send(result, { type: 'level.up', payload: { level: 2 } })
    act(() => result.current.dismissBanner())

    expect(result.current.banner).toBeNull()
  })

  it('celebrates an emotion directly', () => {
    const { result } = renderHook(() => usePetCelebration())

    act(() => result.current.celebrate('eating'))

    expect(result.current.emotion).toBe('eating')
  })

  it('clears pending toast timers on unmount', () => {
    const clearTimeoutSpy = vi.spyOn(globalThis, 'clearTimeout')
    const { result, unmount } = renderHook(() => usePetCelebration())

    send(result, { type: 'xp.gained', payload: { amount: 3 } })
    unmount()

    expect(clearTimeoutSpy).toHaveBeenCalled()
    clearTimeoutSpy.mockRestore()
  })
})
