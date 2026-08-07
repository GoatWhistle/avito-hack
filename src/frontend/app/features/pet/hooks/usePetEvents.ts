import { useEffect, useRef } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { getToken } from '#/api/token-store'
import { connectPetSocket, type PetSocketStatus } from '#/features/pet/lib'
import type { Pet, PetEvent } from '#/features/pet/types'
import { petQueryKey } from './usePetQuery'

const DEFAULT_WS_PATH = '/api/v1/ws'

export interface UsePetEventsOptions {
  onEvent?: (event: PetEvent) => void
  onStatus?: (status: PetSocketStatus) => void
  url?: string
  token?: string | null
  factory?: (url: string) => WebSocket
  enabled?: boolean
}

export const usePetEvents = ({
  onEvent,
  onStatus,
  url = DEFAULT_WS_PATH,
  token,
  factory,
  enabled = true,
}: UsePetEventsOptions = {}) => {
  const queryClient = useQueryClient()
  const eventRef = useRef(onEvent)
  const statusRef = useRef(onStatus)
  eventRef.current = onEvent
  statusRef.current = onStatus

  const resolvedToken = token === undefined ? getToken() : token

  useEffect(() => {
    if (!enabled) return
    if (resolvedToken === null || resolvedToken === '') return
    if (typeof WebSocket === 'undefined' && factory === undefined) return

    return connectPetSocket({
      url,
      token: resolvedToken,
      factory,
      onStatus: (status) => statusRef.current?.(status),
      onEvent: (event) => {
        applyEventToCache(queryClient, event)
        eventRef.current?.(event)
      },
    })
  }, [enabled, factory, queryClient, resolvedToken, url])
}

type Cache = ReturnType<typeof useQueryClient>

const applyEventToCache = (queryClient: Cache, event: PetEvent) => {
  switch (event.type) {
    case 'pet.state':
    case 'pet.updated':
    case 'pet.hatched':
      queryClient.setQueryData<Pet>(petQueryKey, event.payload)

      return
    case 'level.up':
      queryClient.setQueryData<Pet>(petQueryKey, (current) =>
        current === undefined
          ? current
          : { ...current, level: event.payload.level },
      )

      return
    case 'streak.updated':
      queryClient.setQueryData<Pet>(petQueryKey, (current) =>
        current === undefined
          ? current
          : { ...current, streak_days: event.payload.days },
      )

      return
    default:
      return
  }
}
