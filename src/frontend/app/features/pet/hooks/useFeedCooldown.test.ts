import { act, renderHook } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { formatCountdown, useFeedCooldown } from './useFeedCooldown'

const NOW = new Date('2026-08-09T12:00:00.000Z')

const inSeconds = (seconds: number): string =>
  new Date(NOW.getTime() + seconds * 1_000).toISOString()

beforeEach(() => {
  vi.useFakeTimers()
  vi.setSystemTime(NOW)
})

afterEach(() => {
  vi.useRealTimers()
})

describe('formatCountdown', () => {
  it('renders hours, minutes and seconds', () => {
    expect(formatCountdown(5 * 3_600_000)).toBe('5:00:00')
    expect(formatCountdown(65_000)).toBe('0:01:05')
  })

  it('never goes below zero', () => {
    expect(formatCountdown(-5_000)).toBe('0:00:00')
  })
})

describe('useFeedCooldown', () => {
  it('locks feeding while the cooldown is in the future', () => {
    const { result } = renderHook(() => useFeedCooldown(inSeconds(5 * 3_600)))

    expect(result.current.isLocked).toBe(true)
    expect(result.current.remainingMs).toBe(5 * 3_600_000)
    expect(result.current.label).toBe('5:00:00')
  })

  it('counts the remaining time down every second', () => {
    const { result } = renderHook(() => useFeedCooldown(inSeconds(10)))

    expect(result.current.label).toBe('0:00:10')

    act(() => {
      vi.advanceTimersByTime(3_000)
    })

    expect(result.current.remainingMs).toBe(7_000)
    expect(result.current.label).toBe('0:00:07')
  })

  it('unlocks once the cooldown expires', () => {
    const { result } = renderHook(() => useFeedCooldown(inSeconds(2)))

    expect(result.current.isLocked).toBe(true)

    act(() => {
      vi.advanceTimersByTime(2_000)
    })

    expect(result.current.isLocked).toBe(false)
    expect(result.current.remainingMs).toBe(0)
    expect(result.current.label).toBe('0:00:00')
  })

  it('stops ticking after the cooldown expires', () => {
    const clearIntervalSpy = vi.spyOn(globalThis, 'clearInterval')
    renderHook(() => useFeedCooldown(inSeconds(1)))

    act(() => {
      vi.advanceTimersByTime(1_000)
    })

    expect(clearIntervalSpy).toHaveBeenCalled()
    clearIntervalSpy.mockRestore()
  })

  it('clears the interval on unmount', () => {
    const clearIntervalSpy = vi.spyOn(globalThis, 'clearInterval')
    const { unmount } = renderHook(() => useFeedCooldown(inSeconds(3_600)))

    unmount()

    expect(clearIntervalSpy).toHaveBeenCalled()
    clearIntervalSpy.mockRestore()
  })

  it('re-locks when a new cooldown arrives', () => {
    const { result, rerender } = renderHook(
      ({ at }: { at: string | null }) => useFeedCooldown(at),
      { initialProps: { at: null as string | null } },
    )

    expect(result.current.isLocked).toBe(false)

    rerender({ at: inSeconds(60) })

    expect(result.current.isLocked).toBe(true)
    expect(result.current.label).toBe('0:01:00')
  })

  it.each([null, undefined, '', 'not-a-date'])(
    'never locks for the value %s',
    (availableAt) => {
      const { result } = renderHook(() => useFeedCooldown(availableAt))

      expect(result.current.isLocked).toBe(false)
      expect(result.current.remainingMs).toBe(0)
      expect(result.current.label).toBe('0:00:00')
    },
  )

  it('never locks for a cooldown already in the past', () => {
    const { result } = renderHook(() => useFeedCooldown(inSeconds(-60)))

    expect(result.current.isLocked).toBe(false)
  })
})
