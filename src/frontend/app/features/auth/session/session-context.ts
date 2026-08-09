import { createContext, use } from 'react'
import type { User } from '#/types'

export interface SessionValue {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  sessionExpired: boolean
  signOut: () => void
  setUser: (user: User) => void
  acknowledgeExpiry: () => void
}

export const SessionContext = createContext<SessionValue | null>(null)

export const useSession = (): SessionValue => {
  const value = use(SessionContext)

  if (!value) {
    throw new Error('useSession must be used inside SessionProvider')
  }

  return value
}
