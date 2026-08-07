import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { LoadMore } from './LoadMore'
import { renderWithProviders } from './test-utils'

type ObserverCallback = (entries: { isIntersecting: boolean }[]) => void

let callbacks: ObserverCallback[] = []
let disconnects = 0
const original = globalThis.IntersectionObserver

class FakeObserver {
  constructor(private readonly callback: ObserverCallback) {
    callbacks.push(callback)
  }

  observe() {}

  disconnect() {
    disconnects += 1
  }
}

beforeEach(() => {
  callbacks = []
  disconnects = 0
  globalThis.IntersectionObserver =
    FakeObserver as unknown as typeof IntersectionObserver
})

afterEach(() => {
  globalThis.IntersectionObserver = original
})

describe('LoadMore', () => {
  it('renders nothing on the last page', () => {
    const { container } = renderWithProviders(
      <LoadMore hasNextPage={false} isFetching={false} onLoadMore={vi.fn()} />,
    )

    expect(container).toBeEmptyDOMElement()
  })

  it('loads the next page on click', async () => {
    const onLoadMore = vi.fn()
    renderWithProviders(
      <LoadMore hasNextPage isFetching={false} onLoadMore={onLoadMore} />,
    )

    await userEvent.click(screen.getByRole('button'))

    expect(onLoadMore).toHaveBeenCalledTimes(1)
  })

  it('disables the button and shows progress while fetching', () => {
    renderWithProviders(
      <LoadMore hasNextPage isFetching onLoadMore={vi.fn()} />,
    )

    expect(screen.getByRole('button')).toBeDisabled()
  })

  it('loads more when the sentinel scrolls into view', () => {
    const onLoadMore = vi.fn()
    renderWithProviders(
      <LoadMore hasNextPage isFetching={false} onLoadMore={onLoadMore} />,
    )

    expect(callbacks).toHaveLength(1)
    callbacks[0]([{ isIntersecting: true }])

    expect(onLoadMore).toHaveBeenCalledTimes(1)
  })

  it('ignores a sentinel that stays out of view', () => {
    const onLoadMore = vi.fn()
    renderWithProviders(
      <LoadMore hasNextPage isFetching={false} onLoadMore={onLoadMore} />,
    )

    callbacks[0]([{ isIntersecting: false }])

    expect(onLoadMore).not.toHaveBeenCalled()
  })

  it('does not observe while a fetch is already running', () => {
    renderWithProviders(
      <LoadMore hasNextPage isFetching onLoadMore={vi.fn()} />,
    )

    expect(callbacks).toHaveLength(0)
  })

  it('disconnects the observer on unmount', () => {
    const view = renderWithProviders(
      <LoadMore hasNextPage isFetching={false} onLoadMore={vi.fn()} />,
    )

    view.unmount()

    expect(disconnects).toBe(1)
  })

  it('still renders a usable button without IntersectionObserver', async () => {
    Reflect.deleteProperty(
      globalThis as unknown as Record<string, unknown>,
      'IntersectionObserver',
    )
    const onLoadMore = vi.fn()
    renderWithProviders(
      <LoadMore hasNextPage isFetching={false} onLoadMore={onLoadMore} />,
    )

    await userEvent.click(screen.getByRole('button'))

    expect(onLoadMore).toHaveBeenCalledTimes(1)
  })
})
