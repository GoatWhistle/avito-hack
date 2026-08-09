import type { PropsWithChildren } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { act, renderHook, waitFor } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useFeedMutation } from './usePetActions'
import { petQueryKey } from './usePetQuery'
import { makePet } from '#/features/pet/components/test-utils'
import type { Pet } from '#/features/pet/types'

const feed = vi.fn()

vi.mock('#/features/pet/repository', () => ({
  petRepository: {
    feed: () => feed(),
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

describe('useFeedMutation', () => {
  it('optimistically raises satiety and keeps the server value', async () => {
    const client = makeClient()
    client.setQueryData(petQueryKey, makePet({ satiety: 50 }))
    let release = (_pet: Pet) => {}
    feed.mockReturnValue(
      new Promise<Pet>((resolve) => {
        release = resolve
      }),
    )

    const { result } = renderHook(() => useFeedMutation(), {
      wrapper: wrapperFor(client),
    })

    act(() => result.current.mutate())

    await waitFor(() => {
      expect(client.getQueryData<Pet>(petQueryKey)?.satiety).toBe(65)
    })

    await act(async () => {
      release(makePet({ satiety: 68 }))
    })

    await waitFor(() => {
      expect(client.getQueryData<Pet>(petQueryKey)?.satiety).toBe(68)
    })
  })

  it('clamps the optimistic satiety at 100', async () => {
    const client = makeClient()
    client.setQueryData(petQueryKey, makePet({ satiety: 95 }))
    feed.mockImplementation(() => new Promise(() => {}))

    const { result } = renderHook(() => useFeedMutation(), {
      wrapper: wrapperFor(client),
    })

    act(() => result.current.mutate())

    await waitFor(() => {
      expect(client.getQueryData<Pet>(petQueryKey)?.satiety).toBe(100)
    })
  })

  it('rolls back the optimistic update when the request fails', async () => {
    const client = makeClient()
    client.setQueryData(petQueryKey, makePet({ satiety: 30 }))
    feed.mockRejectedValue(new Error('offline'))

    const { result } = renderHook(() => useFeedMutation(), {
      wrapper: wrapperFor(client),
    })

    act(() => result.current.mutate())

    await waitFor(() => expect(result.current.isError).toBe(true))
    expect(client.getQueryData<Pet>(petQueryKey)?.satiety).toBe(30)
  })

  it('leaves happiness untouched while feeding', async () => {
    const client = makeClient()
    client.setQueryData(petQueryKey, makePet({ satiety: 40, happiness: 55 }))
    feed.mockImplementation(() => new Promise(() => {}))

    const { result } = renderHook(() => useFeedMutation(), {
      wrapper: wrapperFor(client),
    })

    act(() => result.current.mutate())

    await waitFor(() => {
      expect(client.getQueryData<Pet>(petQueryKey)?.satiety).toBe(55)
    })
    expect(client.getQueryData<Pet>(petQueryKey)?.happiness).toBe(55)
  })

  it('skips the optimistic write when no pet is cached', async () => {
    const client = makeClient()
    feed.mockRejectedValue(new Error('offline'))

    const { result } = renderHook(() => useFeedMutation(), {
      wrapper: wrapperFor(client),
    })

    act(() => result.current.mutate())

    await waitFor(() => expect(result.current.isError).toBe(true))
    expect(client.getQueryData<Pet>(petQueryKey)).toBeUndefined()
  })
})
