import type { PropsWithChildren } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { act, renderHook, waitFor } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useCheckInMutation, useStrokeMutation } from './usePetActions'
import { petQueryKey } from './usePetQuery'
import { makePet } from '#/features/pet/components/test-utils'
import type { CheckInResult, Pet } from '#/features/pet/types'

const stroke = vi.fn()
const checkIn = vi.fn()
const state = vi.fn()

vi.mock('#/features/pet/repository', () => ({
  petRepository: {
    stroke: () => stroke(),
    checkIn: () => checkIn(),
    state: () => state(),
  },
}))

const makeClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  })

const wrapperFor = (client: QueryClient) =>
  function Wrapper({ children }: PropsWithChildren) {
    return <QueryClientProvider client={client}>{children}</QueryClientProvider>
  }

beforeEach(() => {
  vi.clearAllMocks()
})

describe('useStrokeMutation', () => {
  it('optimistically raises happiness and keeps the server value', async () => {
    const client = makeClient()
    client.setQueryData(petQueryKey, makePet({ happiness: 80 }))
    let release = (_pet: Pet) => {}
    stroke.mockReturnValue(
      new Promise<Pet>((resolve) => {
        release = resolve
      }),
    )

    const { result } = renderHook(() => useStrokeMutation(), {
      wrapper: wrapperFor(client),
    })

    act(() => result.current.mutate())

    await waitFor(() => {
      expect(client.getQueryData<Pet>(petQueryKey)?.happiness).toBe(85)
    })

    await act(async () => {
      release(makePet({ happiness: 91 }))
    })

    await waitFor(() => {
      expect(client.getQueryData<Pet>(petQueryKey)?.happiness).toBe(91)
    })
  })

  it('clamps the optimistic happiness at 100', async () => {
    const client = makeClient()
    client.setQueryData(petQueryKey, makePet({ happiness: 98 }))
    stroke.mockImplementation(() => new Promise(() => {}))

    const { result } = renderHook(() => useStrokeMutation(), {
      wrapper: wrapperFor(client),
    })

    act(() => result.current.mutate())

    await waitFor(() => {
      expect(client.getQueryData<Pet>(petQueryKey)?.happiness).toBe(100)
    })
  })

  it('rolls back the optimistic update when the request fails', async () => {
    const client = makeClient()
    client.setQueryData(petQueryKey, makePet({ happiness: 40 }))
    stroke.mockRejectedValue(new Error('offline'))

    const { result } = renderHook(() => useStrokeMutation(), {
      wrapper: wrapperFor(client),
    })

    act(() => result.current.mutate())

    await waitFor(() => expect(result.current.isError).toBe(true))
    expect(client.getQueryData<Pet>(petQueryKey)?.happiness).toBe(40)
  })

  it('skips the optimistic write when no pet is cached', async () => {
    const client = makeClient()
    stroke.mockRejectedValue(new Error('offline'))

    const { result } = renderHook(() => useStrokeMutation(), {
      wrapper: wrapperFor(client),
    })

    act(() => result.current.mutate())

    await waitFor(() => expect(result.current.isError).toBe(true))
    expect(client.getQueryData<Pet>(petQueryKey)).toBeUndefined()
  })
})

describe('useCheckInMutation', () => {
  const result: CheckInResult = {
    pet: makePet({ streak_days: 5 }),
    xp_gained: 10,
    streak: { days: 5, milestone: true, freezes: 1 },
  } as unknown as CheckInResult

  it('stores the returned pet and notifies the caller', async () => {
    const client = makeClient()
    const onResult = vi.fn()
    checkIn.mockResolvedValue(result)

    const { result: hook } = renderHook(() => useCheckInMutation(onResult), {
      wrapper: wrapperFor(client),
    })

    act(() => hook.current.mutate())

    await waitFor(() => expect(hook.current.isSuccess).toBe(true))
    expect(client.getQueryData<Pet>(petQueryKey)?.streak_days).toBe(5)
    expect(onResult).toHaveBeenCalledWith(result)
  })

  it('works without a result callback', async () => {
    const client = makeClient()
    checkIn.mockResolvedValue(result)

    const { result: hook } = renderHook(() => useCheckInMutation(), {
      wrapper: wrapperFor(client),
    })

    act(() => hook.current.mutate())

    await waitFor(() => expect(hook.current.isSuccess).toBe(true))
    expect(client.getQueryData<Pet>(petQueryKey)?.streak_days).toBe(5)
  })

  it('surfaces a failed check in without touching the cache', async () => {
    const client = makeClient()
    client.setQueryData(petQueryKey, makePet({ streak_days: 2 }))
    checkIn.mockRejectedValue(new Error('already checked in'))

    const { result: hook } = renderHook(() => useCheckInMutation(), {
      wrapper: wrapperFor(client),
    })

    act(() => hook.current.mutate())

    await waitFor(() => expect(hook.current.isError).toBe(true))
    expect(client.getQueryData<Pet>(petQueryKey)?.streak_days).toBe(2)
  })
})
