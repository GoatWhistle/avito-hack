import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import {
  clearToken,
  getToken,
  resetTokenCache,
  setToken,
  subscribeToSessionEnd,
} from './token-store'

beforeEach(() => {
  window.localStorage.clear()
  resetTokenCache()
})

afterEach(() => {
  window.localStorage.clear()
  resetTokenCache()
})

describe('session end notifications', () => {
  it('reports an expiry so the interface can explain it', () => {
    setToken('token')
    const listener = vi.fn()
    const unsubscribe = subscribeToSessionEnd(listener)

    clearToken('expired')

    expect(listener).toHaveBeenCalledWith('expired')
    expect(getToken()).toBeNull()
    unsubscribe()
  })

  it('distinguishes a deliberate sign out from an expiry', () => {
    setToken('token')
    const listener = vi.fn()
    const unsubscribe = subscribeToSessionEnd(listener)

    clearToken('signed-out')

    expect(listener).toHaveBeenCalledWith('signed-out')
    unsubscribe()
  })

  it('defaults to a deliberate sign out', () => {
    setToken('token')
    const listener = vi.fn()
    const unsubscribe = subscribeToSessionEnd(listener)

    clearToken()

    expect(listener).toHaveBeenCalledWith('signed-out')
    unsubscribe()
  })

  it('stays quiet when there was no session to end', () => {
    const listener = vi.fn()
    const unsubscribe = subscribeToSessionEnd(listener)

    clearToken('expired')

    expect(listener).not.toHaveBeenCalled()
    unsubscribe()
  })

  it('stops notifying once unsubscribed', () => {
    setToken('token')
    const listener = vi.fn()
    subscribeToSessionEnd(listener)()

    clearToken('expired')

    expect(listener).not.toHaveBeenCalled()
  })
})
