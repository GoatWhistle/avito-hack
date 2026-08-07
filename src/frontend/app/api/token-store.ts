const STORAGE_KEY = 'avito-hack.token'

type Listener = (token: string | null) => void

let cached: string | null = null
let hydrated = false
const listeners = new Set<Listener>()

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

export const clearToken = () => setToken(null)

export const subscribeToToken = (listener: Listener) => {
  listeners.add(listener)
  return () => {
    listeners.delete(listener)
  }
}

export const resetTokenCache = () => {
  cached = null
  hydrated = false
}
