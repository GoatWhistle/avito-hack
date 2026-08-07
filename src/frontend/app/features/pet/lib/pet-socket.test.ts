import { describe, expect, it, vi } from 'vitest'
import { buildSocketUrl, connectPetSocket } from './pet-socket'
import type { PetEvent } from '#/features/pet/types'

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

describe('buildSocketUrl', () => {
  it('rewrites http to ws and appends the token', () => {
    const url = buildSocketUrl('http://api.test/api/v1/ws', 'jwt-1')

    expect(url).toBe('ws://api.test/api/v1/ws?token=jwt-1')
  })

  it('rewrites https to wss', () => {
    expect(buildSocketUrl('https://api.test/ws', 't')).toBe(
      'wss://api.test/ws?token=t',
    )
  })
})

describe('connectPetSocket', () => {
  const setup = () => {
    const socket = new FakeSocket()
    const onEvent = vi.fn<(event: PetEvent) => void>()
    const dispose = connectPetSocket({
      url: 'http://api.test/ws',
      token: 'jwt',
      onEvent,
      factory: () => socket as unknown as WebSocket,
    })

    return { socket, onEvent, dispose }
  }

  it('requests the state on open', () => {
    const { socket, dispose } = setup()
    socket.onopen?.()

    expect(socket.sent).toContain(JSON.stringify({ type: 'pet.get' }))
    dispose()
  })

  it('forwards parsed events', () => {
    const { socket, onEvent, dispose } = setup()

    socket.onmessage?.({
      data: JSON.stringify({ type: 'xp.gained', payload: { amount: 12 } }),
    } as MessageEvent<string>)

    expect(onEvent).toHaveBeenCalledWith({
      type: 'xp.gained',
      payload: { amount: 12 },
    })
    dispose()
  })

  it('ignores malformed payloads', () => {
    const { socket, onEvent, dispose } = setup()

    socket.onmessage?.({ data: 'not-json' } as MessageEvent<string>)
    socket.onmessage?.({
      data: JSON.stringify({ type: 'unknown.thing' }),
    } as MessageEvent<string>)

    expect(onEvent).not.toHaveBeenCalled()
    dispose()
  })

  it('closes the socket on dispose', () => {
    const { socket, dispose } = setup()
    dispose()

    expect(socket.closed).toBe(true)
  })
})
