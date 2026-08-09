import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, beforeEach, describe, expect, it } from 'vitest'
import { PlatformTourGate } from './PlatformTourGate'
import {
  onboardingStorageKey,
  readOnboardingState,
  requestPlatformTour,
} from './onboarding-state'
import { makeUser, renderWithShell } from '#/features/layout/test-utils'

const session = { user: makeUser(), isAuthenticated: true }

const renderGate = () => renderWithShell(<PlatformTourGate />, { session })

beforeEach(() => {
  window.localStorage.clear()
})

afterEach(() => {
  window.localStorage.clear()
})

describe('PlatformTourGate', () => {
  it('shows the guide to a freshly registered user', async () => {
    requestPlatformTour()

    renderGate()

    const dialog = await screen.findByTestId('platform-tour')
    expect(dialog).toHaveAttribute('aria-modal', 'true')
    expect(dialog).toHaveAccessibleName('Объявления')
    expect(screen.getByTestId('platform-tour-counter')).toHaveTextContent(
      'Шаг 1 из 4',
    )
  })

  it('stays hidden when no registration requested it', async () => {
    renderGate()

    await waitFor(() => {
      expect(screen.queryByTestId('platform-tour')).not.toBeInTheDocument()
    })
  })

  it('moves forward and back through the steps', async () => {
    requestPlatformTour()

    renderGate()
    await screen.findByTestId('platform-tour')

    await userEvent.click(screen.getByRole('button', { name: 'Далее' }))
    expect(screen.getByTestId('platform-tour-counter')).toHaveTextContent(
      'Шаг 2 из 4',
    )
    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent(
      'Питомец',
    )

    await userEvent.click(screen.getByRole('button', { name: 'Назад' }))
    expect(screen.getByTestId('platform-tour-counter')).toHaveTextContent(
      'Шаг 1 из 4',
    )
  })

  it('disables the back button on the first step', async () => {
    requestPlatformTour()

    renderGate()
    await screen.findByTestId('platform-tour')

    expect(screen.getByRole('button', { name: 'Назад' })).toBeDisabled()
  })

  it('hides itself for good once finished', async () => {
    requestPlatformTour()

    renderGate()
    await screen.findByTestId('platform-tour')

    await userEvent.click(screen.getByRole('button', { name: 'Далее' }))
    await userEvent.click(screen.getByRole('button', { name: 'Далее' }))
    await userEvent.click(screen.getByRole('button', { name: 'Далее' }))
    await userEvent.click(screen.getByRole('button', { name: 'Начать' }))

    await waitFor(() => {
      expect(screen.queryByTestId('platform-tour')).not.toBeInTheDocument()
    })
    expect(readOnboardingState().tourSeen).toBe(true)
    expect(readOnboardingState().tourRequested).toBe(false)
  })

  it('hides itself for good once skipped', async () => {
    requestPlatformTour()

    renderGate()
    await screen.findByTestId('platform-tour')

    await userEvent.click(screen.getByRole('button', { name: 'Пропустить' }))

    await waitFor(() => {
      expect(screen.queryByTestId('platform-tour')).not.toBeInTheDocument()
    })
    expect(readOnboardingState().tourSeen).toBe(true)
  })

  it('closes on Escape and counts as seen', async () => {
    requestPlatformTour()

    renderGate()
    await screen.findByTestId('platform-tour')

    await userEvent.keyboard('{Escape}')

    await waitFor(() => {
      expect(screen.queryByTestId('platform-tour')).not.toBeInTheDocument()
    })
    expect(readOnboardingState().tourSeen).toBe(true)
  })

  it('moves focus into the dialog when it opens', async () => {
    requestPlatformTour()

    const { container } = renderGate()
    const dialog = await screen.findByTestId('platform-tour')

    await waitFor(() => {
      expect(dialog.contains(document.activeElement)).toBe(true)
    })
    expect(container.contains(document.activeElement)).toBe(false)
  })

  it('keeps Tab focus inside the dialog', async () => {
    requestPlatformTour()

    renderGate()
    const dialog = await screen.findByTestId('platform-tour')

    for (let step = 0; step < 6; step += 1) {
      await userEvent.tab()
      expect(dialog.contains(document.activeElement)).toBe(true)
    }

    await userEvent.tab({ shift: true })
    expect(dialog.contains(document.activeElement)).toBe(true)
  })

  it('does not return on a later visit', async () => {
    requestPlatformTour()

    const first = renderGate()
    await screen.findByTestId('platform-tour')
    await userEvent.click(screen.getByRole('button', { name: 'Пропустить' }))
    first.unmount()

    renderGate()

    await waitFor(() => {
      expect(screen.queryByTestId('platform-tour')).not.toBeInTheDocument()
    })
  })

  it('is not re-requested after it has been seen', async () => {
    requestPlatformTour()
    const seen = { ...readOnboardingState(), tourSeen: true }
    window.localStorage.setItem(onboardingStorageKey, JSON.stringify(seen))

    requestPlatformTour()

    expect(readOnboardingState().tourRequested).toBe(false)

    renderGate()

    await waitFor(() => {
      expect(screen.queryByTestId('platform-tour')).not.toBeInTheDocument()
    })
  })
})
