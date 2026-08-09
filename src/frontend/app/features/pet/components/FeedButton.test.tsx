import { act, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { FeedButton } from './FeedButton'
import { renderWithProviders } from './test-utils'

const renderButton = (overrides: Partial<Parameters<typeof FeedButton>[0]>) => {
  const onFeed = overrides.onFeed ?? vi.fn()

  renderWithProviders(
    <FeedButton
      value={overrides.value ?? 50}
      isPending={overrides.isPending ?? false}
      error={overrides.error ?? null}
      availableAt={overrides.availableAt ?? null}
      onFeed={onFeed}
    />,
  )

  return { onFeed }
}

const NOW = new Date('2026-08-09T12:00:00.000Z')

const inSeconds = (seconds: number): string =>
  new Date(NOW.getTime() + seconds * 1_000).toISOString()

describe('FeedButton', () => {
  it('exposes an accessible name and triggers feeding on click', async () => {
    const user = userEvent.setup()
    const { onFeed } = renderButton({ value: 40 })

    const button = screen.getByRole('button', {
      name: /покормить питомца/i,
    })
    await user.click(button)

    expect(onFeed).toHaveBeenCalledTimes(1)
  })

  it('is reachable and activatable from the keyboard', async () => {
    const user = userEvent.setup()
    const { onFeed } = renderButton({ value: 40 })

    await user.tab()
    expect(screen.getByTestId('feed-button')).toHaveFocus()

    await user.keyboard('{Enter}')
    expect(onFeed).toHaveBeenCalledTimes(1)

    await user.keyboard(' ')
    expect(onFeed).toHaveBeenCalledTimes(2)
  })

  it('marks itself busy and blocks repeat clicks while pending', async () => {
    const user = userEvent.setup()
    const { onFeed } = renderButton({ value: 40, isPending: true })

    const button = screen.getByTestId('feed-button')
    expect(button).toBeDisabled()
    expect(button).toHaveAttribute('aria-busy', 'true')
    expect(screen.getByRole('status')).toHaveTextContent(/кормим/i)

    await user.click(button)
    expect(onFeed).not.toHaveBeenCalled()
  })

  it('disables feeding and explains why when the pet is full', () => {
    renderButton({ value: 100 })

    expect(screen.getByTestId('feed-button')).toBeDisabled()
    expect(screen.getByRole('status')).toHaveTextContent(/сыт/i)
  })

  it('announces errors politely', () => {
    renderButton({ value: 40, error: 'Сеть недоступна' })

    const status = screen.getByRole('status')
    expect(status).toHaveAttribute('aria-live', 'polite')
    expect(status).toHaveTextContent('Сеть недоступна')
  })
})

describe('FeedButton cooldown', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    vi.setSystemTime(NOW)
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('blocks feeding and shows the countdown while the cooldown runs', () => {
    const { onFeed } = renderButton({
      value: 40,
      availableAt: inSeconds(5 * 3_600),
    })

    const button = screen.getByTestId('feed-button')
    expect(button).toBeDisabled()
    expect(screen.getByTestId('feed-countdown')).toHaveTextContent('5:00:00')
    expect(screen.getByRole('status')).toHaveTextContent(
      'Снова покормить можно через 5:00:00',
    )
    expect(onFeed).not.toHaveBeenCalled()
  })

  it('ticks the countdown down while mounted', () => {
    renderButton({ value: 40, availableAt: inSeconds(10) })

    expect(screen.getByTestId('feed-countdown')).toHaveTextContent('0:00:10')

    act(() => {
      vi.advanceTimersByTime(4_000)
    })

    expect(screen.getByTestId('feed-countdown')).toHaveTextContent('0:00:06')
  })

  it('enables feeding once the cooldown expires', () => {
    renderButton({ value: 40, availableAt: inSeconds(3) })

    expect(screen.getByTestId('feed-button')).toBeDisabled()

    act(() => {
      vi.advanceTimersByTime(3_000)
    })

    const button = screen.getByTestId('feed-button')
    expect(button).toBeEnabled()
    expect(button).toHaveAccessibleName(/Покормить/)
    expect(screen.queryByTestId('feed-countdown')).not.toBeInTheDocument()
    expect(screen.getByRole('status')).toHaveTextContent('')
  })

  it('enables feeding for a cooldown already in the past', () => {
    renderButton({ value: 40, availableAt: inSeconds(-60) })

    expect(screen.getByTestId('feed-button')).toBeEnabled()
  })

  it('stays blocked at full satiety even without a cooldown', () => {
    renderButton({ value: 100, availableAt: null })

    expect(screen.getByTestId('feed-button')).toBeDisabled()
    expect(screen.getByRole('status')).toHaveTextContent('Питомец сыт')
  })

  it('prefers the error over the cooldown hint', () => {
    renderButton({
      value: 40,
      error: 'Сеть недоступна',
      availableAt: inSeconds(600),
    })

    expect(screen.getByRole('status')).toHaveTextContent('Сеть недоступна')
    expect(screen.getByTestId('feed-button')).toBeDisabled()
  })
})
