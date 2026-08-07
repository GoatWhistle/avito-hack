import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import {
  defaultOnboardingState,
  onboardingStorageKey,
  readOnboardingState,
  writeOnboardingState,
} from './onboarding-state'

beforeEach(() => {
  window.localStorage.clear()
})

afterEach(() => {
  vi.restoreAllMocks()
  window.localStorage.clear()
})

describe('readOnboardingState', () => {
  it('returns the defaults with nothing stored', () => {
    expect(readOnboardingState()).toEqual(defaultOnboardingState)
  })

  it('returns the defaults for an empty string', () => {
    window.localStorage.setItem(onboardingStorageKey, '')

    expect(readOnboardingState()).toEqual(defaultOnboardingState)
  })

  it('merges a partial stored state over the defaults', () => {
    window.localStorage.setItem(
      onboardingStorageKey,
      JSON.stringify({ dismissed: true }),
    )

    expect(readOnboardingState()).toEqual({
      dismissed: true,
      hatched: false,
      celebrated: false,
    })
  })

  it('reads a fully stored state', () => {
    const state = { dismissed: true, hatched: true, celebrated: true }
    window.localStorage.setItem(onboardingStorageKey, JSON.stringify(state))

    expect(readOnboardingState()).toEqual(state)
  })

  it('ignores a non object payload', () => {
    window.localStorage.setItem(onboardingStorageKey, JSON.stringify(42))

    expect(readOnboardingState()).toEqual(defaultOnboardingState)
  })

  it('ignores a stored null', () => {
    window.localStorage.setItem(onboardingStorageKey, 'null')

    expect(readOnboardingState()).toEqual(defaultOnboardingState)
  })

  it('recovers from malformed json', () => {
    window.localStorage.setItem(onboardingStorageKey, '{oops')

    expect(readOnboardingState()).toEqual(defaultOnboardingState)
  })

  it('recovers when storage access throws', () => {
    vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => {
      throw new Error('blocked')
    })

    expect(readOnboardingState()).toEqual(defaultOnboardingState)
  })
})

describe('writeOnboardingState', () => {
  it('round trips a state through storage', () => {
    const state = { dismissed: true, hatched: true, celebrated: false }

    writeOnboardingState(state)

    expect(readOnboardingState()).toEqual(state)
  })

  it('swallows quota errors', () => {
    vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {
      throw new Error('quota')
    })

    expect(() => writeOnboardingState(defaultOnboardingState)).not.toThrow()
  })
})
