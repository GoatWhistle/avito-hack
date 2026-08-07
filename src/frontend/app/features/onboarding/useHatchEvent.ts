import { useEffect, useRef, useState } from 'react'
import { apiBaseUrl, getToken } from '#/api'

const hatchTypes = new Set(['pet.hatched', 'pet.hatch'])

const toWsUrl = () => {
  const url = new URL(apiBaseUrl, window.location.origin)
  url.protocol = url.protocol === 'https:' ? 'wss:' : 'ws:'
  url.pathname = `${url.pathname.replace(/\/+$/, '')}/ws`

  return url
}

const isHatchPayload = (data: unknown): boolean => {
  if (typeof data !== 'string') return false

  try {
    const message: unknown = JSON.parse(data)
    if (typeof message !== 'object' || message === null) return false

    const type = (message as { type?: unknown }).type
    if (typeof type !== 'string') return false
    if (hatchTypes.has(type)) return true

    if (type === 'pet.updated' || type === 'pet.state') {
      const payload = (message as { payload?: { stage?: unknown } }).payload
      return typeof payload?.stage === 'string' && payload.stage === 'baby'
    }

    return false
  } catch {
    return false
  }
}

export const useHatchEvent = (enabled: boolean) => {
  const [hatched, setHatched] = useState(false)
  const socketRef = useRef<WebSocket | null>(null)

  useEffect(() => {
    if (!enabled) return
    if (typeof window === 'undefined' || typeof WebSocket === 'undefined')
      return

    const token = getToken()
    if (!token) return

    let cancelled = false
    const url = toWsUrl()
    url.searchParams.set('token', token)

    let socket: WebSocket
    try {
      socket = new WebSocket(url.toString())
    } catch {
      return
    }

    socketRef.current = socket

    const onMessage = (event: MessageEvent<unknown>) => {
      if (cancelled) return
      if (isHatchPayload(event.data)) setHatched(true)
    }

    socket.addEventListener('message', onMessage)

    return () => {
      cancelled = true
      socket.removeEventListener('message', onMessage)
      socket.close()
      socketRef.current = null
    }
  }, [enabled])

  return hatched
}
