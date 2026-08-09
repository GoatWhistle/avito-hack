const STORAGE_KEY = 'avito-hack.token'

export type SessionEndReason = 'expired' | 'signed-out'

type Listener = (token: string | null) => void
type ReasonListener = (reason: SessionEndReason) => void

let cached: string | null = null
let hydrated = false
const listeners = new Set<Listener>()
const reasonListeners = new Set<ReasonListener>()

const readStorage = (): string | null => {
  if (typeof window === 'undefined') return null
  try {
    return window.localStorage.getItem(STORAGE_KEY)
  } catch {
    return null
  }
}

const writeStorage = (token: string | null) => {
  if (typeof window === 'undefined') return
  try {
    if (token === null) {
      window.localStorage.removeItem(STORAGE_KEY)
    } else {
      window.localStorage.setItem(STORAGE_KEY, token)
    }
  } catch {
    return
  }
}

export const getToken = (): string | null => {
  if (!hydrated) {
    cached = readStorage()
    hydrated = true
  }
  return cached
}

export const setToken = (token: string | null) => {
  cached = token
  hydrated = true
  writeStorage(token)
  for (const listener of listeners) listener(token)
}

export const clearToken = (reason: SessionEndReason = 'signed-out') => {
  const had = getToken() !== null
  setToken(null)
  if (!had) return
  for (const listener of reasonListeners) listener(reason)
}

export const subscribeToToken = (listener: Listener) => {
  listeners.add(listener)
  return () => {
    listeners.delete(listener)
  }
}

export const subscribeToSessionEnd = (listener: ReasonListener) => {
  reasonListeners.add(listener)
  return () => {
    reasonListeners.delete(listener)
  }
}

export const resetTokenCache = () => {
  cached = null
  hydrated = false
}
