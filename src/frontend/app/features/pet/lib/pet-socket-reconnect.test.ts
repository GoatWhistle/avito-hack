import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import {
  buildSocketUrl,
  connectPetSocket,
  type PetSocketStatus,
} from './pet-socket'

class FakeSocket {
  onopen: (() => void) | null = null
  onmessage: ((event: MessageEvent<string>) => void) | null = null
  onerror: (() => void) | null = null
  onclose: (() => void) | null = null
  readonly sent: string[] = []
  closed = false

  send(data: string) {
    this.sent.push(data)
  }

  close() {
    this.closed = true
  }
}

beforeEach(() => {
  vi.useFakeTimers()
})

afterEach(() => {
  vi.useRealTimers()
})

const setup = () => {
  const sockets: FakeSocket[] = []
  const statuses: PetSocketStatus[] = []
  const dispose = connectPetSocket({
    url: 'http://api.test/ws',
    token: 'jwt',
    onEvent: vi.fn(),
    onStatus: (status) => statuses.push(status),
    factory: () => {
      const socket = new FakeSocket()
      sockets.push(socket)

      return socket as unknown as WebSocket
    },
  })

  return { sockets, statuses, dispose }
}

describe('buildSocketUrl', () => {
  it('resolves a relative path against the page origin', () => {
    expect(buildSocketUrl('/api/v1/ws', 'tok')).toBe(
      `ws://${window.location.host}/api/v1/ws?token=tok`,
    )
  })

  it('replaces an existing token instead of duplicating it', () => {
    expect(buildSocketUrl('http://api.test/ws?token=old', 'new')).toBe(
      'ws://api.test/ws?token=new',
    )
  })
})

describe('connectPetSocket lifecycle', () => {
  it('reports connecting then open', () => {
    const { sockets, statuses, dispose } = setup()
    sockets[0].onopen?.()

    expect(statuses).toEqual(['connecting', 'open'])
    dispose()
  })

  it('sends periodic pings while open', () => {
    const { sockets, dispose } = setup()
    sockets[0].onopen?.()

    vi.advanceTimersByTime(25_000)
    expect(sockets[0].sent).toContain(JSON.stringify({ type: 'ping' }))

    vi.advanceTimersByTime(25_000)
    const pings = sockets[0].sent.filter((raw) => raw.includes('ping'))
    expect(pings).toHaveLength(2)
    dispose()
  })

  it('closes the socket when it reports an error', () => {
    const { sockets, dispose } = setup()
    sockets[0].onopen?.()

    sockets[0].onerror?.()

    expect(sockets[0].closed).toBe(true)
    dispose()
  })

  it('reconnects with backoff after a close', () => {
    const { sockets, statuses, dispose } = setup()
    sockets[0].onopen?.()
    sockets[0].onclose?.()

    expect(statuses).toContain('closed')
    expect(sockets).toHaveLength(1)

    vi.advanceTimersByTime(1_000)
    expect(sockets).toHaveLength(2)

    sockets[1].onclose?.()
    vi.advanceTimersByTime(1_000)
    expect(sockets).toHaveLength(2)

    vi.advanceTimersByTime(1_000)
    expect(sockets).toHaveLength(3)
    dispose()
  })

  it('resets the backoff after a successful reconnect', () => {
    const { sockets, dispose } = setup()
    sockets[0].onopen?.()
    sockets[0].onclose?.()
    vi.advanceTimersByTime(1_000)

    sockets[1].onopen?.()
    sockets[1].onclose?.()

    vi.advanceTimersByTime(1_000)
    expect(sockets).toHaveLength(3)
    dispose()
  })

  it('stops pinging once the socket closes', () => {
    const { sockets, dispose } = setup()
    sockets[0].onopen?.()
    sockets[0].onclose?.()

    const before = sockets[0].sent.length
    vi.advanceTimersByTime(60_000)

    expect(sockets[0].sent).toHaveLength(before)
    dispose()
  })

  it('retries when the socket factory throws', () => {
    const statuses: PetSocketStatus[] = []
    let attempts = 0
    const dispose = connectPetSocket({
      url: 'http://api.test/ws',
      token: 'jwt',
      onEvent: vi.fn(),
      onStatus: (status) => statuses.push(status),
      factory: () => {
        attempts += 1
        throw new Error('blocked')
      },
    })

    expect(attempts).toBe(1)
    expect(statuses).toEqual(['connecting'])

    vi.advanceTimersByTime(1_000)
    expect(attempts).toBe(2)

    dispose()
    vi.advanceTimersByTime(60_000)
    expect(attempts).toBe(2)
  })

  it('stops reconnecting and reporting status after dispose', () => {
    const { sockets, statuses, dispose } = setup()
    sockets[0].onopen?.()
    dispose()

    const seen = statuses.length
    sockets[0].onclose?.()
    vi.advanceTimersByTime(60_000)

    expect(sockets).toHaveLength(1)
    expect(statuses).toHaveLength(seen)
  })

  it('detaches handlers on dispose', () => {
    const { sockets, dispose } = setup()
    sockets[0].onopen?.()
    dispose()

    expect(sockets[0].onclose).toBeNull()
    expect(sockets[0].onerror).toBeNull()
    expect(sockets[0].onmessage).toBeNull()
    expect(sockets[0].closed).toBe(true)
  })
})
