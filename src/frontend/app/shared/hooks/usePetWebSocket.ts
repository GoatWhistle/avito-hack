import { TAMAGOTCHI_QUERY_KEYS } from '#/features/tamagotchi/lib/query-keys'
import { localTokenStorage } from '#/shared/storage/local-token.storage'
import type { Pet } from '#/shared/types/pet.type'
import type { WSServerMessage } from '#/shared/types/websocket.type'
import { useQueryClient } from '@tanstack/react-query'
import { useEffect, useRef } from 'react'

export function usePetWebSocket() {
  const token = localTokenStorage.get()
  const queryClient = useQueryClient()
  const wsRef = useRef<WebSocket | null>(null)
  const pingIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null)

  useEffect(() => {
    if (!token) return

    async function handleMessage(msg: WSServerMessage) {
      const { type, payload } = msg

      switch (type) {
        case 'pet.state':
        case 'pet.updated':
        case 'pet.hatched': {
          const pet = payload as Pet
          queryClient.setQueryData(TAMAGOTCHI_QUERY_KEYS.pet.my(), pet)
          queryClient.setQueryData(TAMAGOTCHI_QUERY_KEYS.raccoon.my(), {
            level: pet.level,
            xp: pet.xp,
            xp_to_next_level: pet.next_level_xp,
            stage: pet.stage,
            streak_days: pet.streak_days,
          })
          break
        }

        case 'xp.gained': {
          await queryClient.invalidateQueries({
            queryKey: TAMAGOTCHI_QUERY_KEYS.pet.my(),
          })
          await queryClient.invalidateQueries({
            queryKey: TAMAGOTCHI_QUERY_KEYS.raccoon.my(),
          })
          break
        }

        case 'level.up': {
          await queryClient.invalidateQueries({
            queryKey: TAMAGOTCHI_QUERY_KEYS.pet.my(),
          })
          await queryClient.invalidateQueries({
            queryKey: TAMAGOTCHI_QUERY_KEYS.raccoon.my(),
          })
          break
        }

        case 'reward.granted': {
          await queryClient.invalidateQueries({
            queryKey: TAMAGOTCHI_QUERY_KEYS.rewards.my(),
          })
          await queryClient.invalidateQueries({
            queryKey: TAMAGOTCHI_QUERY_KEYS.rewards.all(),
          })
          break
        }

        case 'streak.updated': {
          await queryClient.invalidateQueries({
            queryKey: TAMAGOTCHI_QUERY_KEYS.pet.my(),
          })
          await queryClient.invalidateQueries({
            queryKey: TAMAGOTCHI_QUERY_KEYS.raccoon.my(),
          })
          break
        }

        case 'error':
          console.error('WS error:', payload)
          break

        case 'pong':
        default:
          break
      }
    }

    const connect = () => {
      const wsUrl = `/api/v1/ws?token=${encodeURIComponent(token)}`
      const ws = new WebSocket(wsUrl)
      wsRef.current = ws

      ws.onopen = () => {
        if (pingIntervalRef.current !== null) {
          clearInterval(pingIntervalRef.current)
          pingIntervalRef.current = null
        }
        pingIntervalRef.current = setInterval(() => {
          if (ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({ type: 'ping' }))
          }
        }, 30_000)
      }

      ws.onmessage = async event => {
        try {
          const msg: WSServerMessage = JSON.parse(event.data)
          await handleMessage(msg)
        } catch (e) {
          console.error('Invalid WS message', e)
        }
      }

      ws.onclose = () => {
        if (pingIntervalRef.current !== null) {
          clearInterval(pingIntervalRef.current)
          pingIntervalRef.current = null
        }
        setTimeout(connect, 5000)
      }

      ws.onerror = e => {
        console.error('WebSocket error', e)
        ws.close()
      }
    }

    connect()

    return () => {
      if (pingIntervalRef.current !== null) {
        clearInterval(pingIntervalRef.current)
        pingIntervalRef.current = null
      }
      wsRef.current?.close()
    }
  }, [token, queryClient])
}
