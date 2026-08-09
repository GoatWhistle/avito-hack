import {
  useCallback,
  useEffect,
  useMemo,
  useState,
  type PropsWithChildren,
} from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import {
  clearToken,
  getToken,
  subscribeToSessionEnd,
  subscribeToToken,
} from '#/api'
import { authRepository } from '#/features/auth/repository'
import type { User } from '#/types'
import { SessionContext, type SessionValue } from './session-context'

export const sessionQueryKey = ['session', 'me'] as const

export function SessionProvider({ children }: PropsWithChildren) {
  const queryClient = useQueryClient()
  const [token, setTokenState] = useState<string | null>(null)
  const [hydrated, setHydrated] = useState(false)
  const [expired, setExpired] = useState(false)

  useEffect(() => {
    setTokenState(getToken())
    setHydrated(true)

    const unsubscribeToken = subscribeToToken((next) => {
      setTokenState(next)
      if (next !== null) setExpired(false)
    })

    const unsubscribeEnd = subscribeToSessionEnd((reason) => {
      setExpired(reason === 'expired')
    })

    return () => {
      unsubscribeToken()
      unsubscribeEnd()
    }
  }, [])

  const { data, isLoading } = useQuery({
    queryKey: sessionQueryKey,
    queryFn: () => authRepository.me(),
    enabled: hydrated && token !== null,
    retry: false,
    staleTime: 60_000,
  })

  const signOut = useCallback(() => {
    clearToken('signed-out')
    queryClient.clear()
  }, [queryClient])

  const acknowledgeExpiry = useCallback(() => {
    setExpired(false)
  }, [])

  const setUser = useCallback(
    (user: User) => {
      queryClient.setQueryData(sessionQueryKey, user)
    },
    [queryClient],
  )

  const value = useMemo<SessionValue>(
    () => ({
      user: data ?? null,
      isAuthenticated: Boolean(token && data),
      isLoading: !hydrated || (token !== null && isLoading),
      sessionExpired: expired,
      signOut,
      setUser,
      acknowledgeExpiry,
    }),
    [
      data,
      token,
      hydrated,
      isLoading,
      expired,
      signOut,
      setUser,
      acknowledgeExpiry,
    ],
  )

  return <SessionContext value={value}>{children}</SessionContext>
}
