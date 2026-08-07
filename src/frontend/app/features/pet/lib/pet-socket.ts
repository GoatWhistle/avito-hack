import { parsePetEvent, type PetEvent } from '#/features/pet/types'

const RECONNECT_BASE_MS = 1_000
const RECONNECT_MAX_MS = 15_000
const PING_INTERVAL_MS = 25_000

export type PetSocketStatus = 'idle' | 'connecting' | 'open' | 'closed'

export interface PetSocketHandlers {
  onEvent: (event: PetEvent) => void
  onStatus?: (status: PetSocketStatus) => void
}

export interface PetSocketOptions extends PetSocketHandlers {
  url: string
  token: string
  factory?: (url: string) => WebSocket
}

const backoff = (attempt: number): number =>
  Math.min(RECONNECT_BASE_MS * 2 ** attempt, RECONNECT_MAX_MS)

export const buildSocketUrl = (base: string, token: string): string => {
  const url = new URL(base, resolveOrigin())
  url.protocol = url.protocol === 'https:' ? 'wss:' : 'ws:'
  url.searchParams.set('token', token)

  return url.toString()
}

const resolveOrigin = (): string =>
  typeof window === 'undefined' ? 'http://localhost' : window.location.origin

export function connectPetSocket({
  url,
  token,
  onEvent,
  onStatus,
  factory,
}: PetSocketOptions): () => void {
  const create = factory ?? ((target: string) => new WebSocket(target))

  let socket: WebSocket | null = null
  let attempt = 0
  let disposed = false
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null
  let pingTimer: ReturnType<typeof setInterval> | null = null

  const setStatus = (status: PetSocketStatus) => {
    if (!disposed) onStatus?.(status)
  }

  const clearTimers = () => {
    if (reconnectTimer !== null) clearTimeout(reconnectTimer)
    if (pingTimer !== null) clearInterval(pingTimer)
    reconnectTimer = null
    pingTimer = null
  }

  const scheduleReconnect = () => {
    if (disposed) return
    reconnectTimer = setTimeout(open, backoff(attempt))
    attempt += 1
  }

  function open() {
    if (disposed) return
    setStatus('connecting')

    try {
      socket = create(buildSocketUrl(url, token))
    } catch {
      scheduleReconnect()

      return
    }

    socket.onopen = () => {
      attempt = 0
      setStatus('open')
      pingTimer = setInterval(() => {
        socket?.send(JSON.stringify({ type: 'ping' }))
      }, PING_INTERVAL_MS)
      socket?.send(JSON.stringify({ type: 'pet.get' }))
    }

    socket.onmessage = (message: MessageEvent<string>) => {
      const event = safeParse(message.data)
      if (event !== null) onEvent(event)
    }

    socket.onerror = () => {
      socket?.close()
    }

    socket.onclose = () => {
      if (pingTimer !== null) clearInterval(pingTimer)
      pingTimer = null
      setStatus('closed')
      scheduleReconnect()
    }
  }

  open()

  return () => {
    disposed = true
    clearTimers()
    if (socket !== null) {
      socket.onclose = null
      socket.onerror = null
      socket.onmessage = null
      socket.close()
    }
  }
}

const safeParse = (raw: string): PetEvent | null => {
  try {
    return parsePetEvent(JSON.parse(raw))
  } catch {
    return null
  }
}
