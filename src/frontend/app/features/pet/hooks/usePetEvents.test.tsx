import type { PropsWithChildren } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { act, renderHook } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { usePetEvents } from './usePetEvents'
import { petQueryKey } from './usePetQuery'
import { makePet } from '#/features/pet/components/test-utils'
import type { Pet } from '#/features/pet/types'

const getToken = vi.fn()

vi.mock('#/api/token-store', () => ({
  getToken: () => getToken(),
}))

class FakeSocket {
  static last: FakeSocket | null = null
  static created = 0
  onopen: (() => void) | null = null
  onmessage: ((event: MessageEvent<string>) => void) | null = null
  onerror: (() => void) | null = null
  onclose: (() => void) | null = null

  constructor() {
    FakeSocket.last = this
    FakeSocket.created += 1
  }

  send() {}
  close() {}
}

const factory = () => new FakeSocket() as unknown as WebSocket

const makeClient = () =>
  new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } })

const wrapperFor = (client: QueryClient) =>
  function Wrapper({ children }: PropsWithChildren) {
    return <QueryClientProvider client={client}>{children}</QueryClientProvider>
  }

const emit = (type: string, payload: unknown) => {
  act(() => {
    FakeSocket.last?.onmessage?.({
      data: JSON.stringify({ type, payload }),
    } as MessageEvent<string>)
  })
}

beforeEach(() => {
  vi.clearAllMocks()
  FakeSocket.last = null
  FakeSocket.created = 0
  getToken.mockReturnValue('jwt')
})

const setup = (client: QueryClient, options = {}) => {
  const onEvent = vi.fn()
  const onStatus = vi.fn()
  const view = renderHook(
    () => usePetEvents({ factory, onEvent, onStatus, ...options }),
    { wrapper: wrapperFor(client) },
  )

  return { onEvent, onStatus, view }
}

describe('usePetEvents connection guards', () => {
  it('connects with the token from the store', () => {
    setup(makeClient())

    expect(FakeSocket.created).toBe(1)
  })

  it('does not connect when disabled', () => {
    setup(makeClient(), { enabled: false })

    expect(FakeSocket.created).toBe(0)
  })

  it('does not connect without a token', () => {
    getToken.mockReturnValue(null)
    setup(makeClient())

    expect(FakeSocket.created).toBe(0)
  })

  it('does not connect for an empty token', () => {
    setup(makeClient(), { token: '' })

    expect(FakeSocket.created).toBe(0)
  })

  it('uses an explicitly provided token over the store', () => {
    setup(makeClient(), { token: 'explicit' })

    expect(getToken).not.toHaveBeenCalled()
    expect(FakeSocket.created).toBe(1)
  })

  it('reconnects when the token changes', () => {
    const client = makeClient()
    const { rerender } = renderHook(
      ({ token }: { token: string }) => usePetEvents({ factory, token }),
      {
        wrapper: wrapperFor(client),
        initialProps: { token: 'a' },
      },
    )

    expect(FakeSocket.created).toBe(1)
    rerender({ token: 'b' })
    expect(FakeSocket.created).toBe(2)
  })
})

describe('usePetEvents cache updates', () => {
  it('stores the pet from pet.state', () => {
    const client = makeClient()
    const { onEvent } = setup(client)

    const pet = makePet({ satiety: 33 })
    emit('pet.state', pet)

    expect(client.getQueryData<Pet>(petQueryKey)?.satiety).toBe(33)
    expect(onEvent).toHaveBeenCalledTimes(1)
  })

  it('stores the pet from pet.updated', () => {
    const client = makeClient()
    setup(client)

    emit('pet.updated', makePet({ energy: 12 }))
    expect(client.getQueryData<Pet>(petQueryKey)?.energy).toBe(12)

    emit('pet.updated', makePet({ stage: 'teen', energy: 50 }))
    expect(client.getQueryData<Pet>(petQueryKey)?.stage).toBe('teen')
  })

  it('patches the level on level.up', () => {
    const client = makeClient()
    client.setQueryData(petQueryKey, makePet({ level: 3 }))
    setup(client)

    emit('level.up', { level: 9 })

    expect(client.getQueryData<Pet>(petQueryKey)?.level).toBe(9)
  })

  it('ignores level.up without a cached pet', () => {
    const client = makeClient()
    setup(client)

    emit('level.up', { level: 9 })

    expect(client.getQueryData<Pet>(petQueryKey)).toBeUndefined()
  })

  it('patches the streak on streak.updated', () => {
    const client = makeClient()
    client.setQueryData(petQueryKey, makePet({ streak_days: 1 }))
    setup(client)

    emit('streak.updated', { days: 6 })

    expect(client.getQueryData<Pet>(petQueryKey)?.streak_days).toBe(6)
  })

  it('ignores streak.updated without a cached pet', () => {
    const client = makeClient()
    setup(client)

    emit('streak.updated', { days: 6 })

    expect(client.getQueryData<Pet>(petQueryKey)).toBeUndefined()
  })

  it('forwards non cache events untouched', () => {
    const client = makeClient()
    const { onEvent } = setup(client)

    emit('xp.gained', { amount: 4 })

    expect(client.getQueryData<Pet>(petQueryKey)).toBeUndefined()
    expect(onEvent).toHaveBeenCalledWith({
      type: 'xp.gained',
      payload: { amount: 4 },
    })
  })

  it('drops malformed frames', () => {
    const client = makeClient()
    const { onEvent } = setup(client)

    emit('nope.unknown', {})

    expect(onEvent).not.toHaveBeenCalled()
  })

  it('reports socket status changes', () => {
    const { onStatus } = setup(makeClient())

    expect(onStatus).toHaveBeenCalledWith('connecting')
    act(() => FakeSocket.last?.onopen?.())
    expect(onStatus).toHaveBeenCalledWith('open')
  })

  it('works without optional callbacks', () => {
    const client = makeClient()
    renderHook(() => usePetEvents({ factory }), { wrapper: wrapperFor(client) })

    emit('pet.state', makePet({ satiety: 51 }))

    expect(client.getQueryData<Pet>(petQueryKey)?.satiety).toBe(51)
  })
})
