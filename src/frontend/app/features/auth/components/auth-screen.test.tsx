import { cleanup, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import type * as routerModule from 'react-router'
import { setToken } from '#/api/token-store'
import { AuthScreen } from './AuthScreen'
import { renderWithShell } from '#/features/layout/test-utils'

type RouterModule = typeof routerModule

const navigate = vi.fn()

vi.mock('#/features/auth/use-cases', () => ({
  signInUseCase: { execute: vi.fn() },
  signUpUseCase: { execute: vi.fn() },
}))

vi.mock('react-router', async () => {
  const actual = await vi.importActual<RouterModule>('react-router')

  return { ...actual, useNavigate: () => navigate }
})

let matchMediaResult = false

beforeEach(() => {
  vi.clearAllMocks()
  setToken(null)
  matchMediaResult = false
  window.matchMedia = vi.fn().mockImplementation((query: string) => ({
    matches: query.includes('prefers-reduced-motion') && matchMediaResult,
    media: query,
    onchange: null,
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    addListener: vi.fn(),
    removeListener: vi.fn(),
    dispatchEvent: vi.fn(),
  }))
})

describe('AuthScreen', () => {
  it('opens on the sign in tab and shows only the sign in fields', () => {
    renderWithShell(<AuthScreen mode="signIn" />)

    expect(screen.getByRole('tab', { name: 'Войти' })).toHaveAttribute(
      'data-active',
    )
    expect(screen.getByLabelText('Электронная почта')).toBeInTheDocument()
    expect(screen.getByLabelText('Пароль')).toBeInTheDocument()
    expect(screen.queryByLabelText('Имя и фамилия')).not.toBeInTheDocument()
  })

  it('opens on the sign up tab and adds the full name field', () => {
    renderWithShell(<AuthScreen mode="signUp" />)

    expect(screen.getByRole('tab', { name: 'Регистрация' })).toHaveAttribute(
      'data-active',
    )
    expect(screen.getByLabelText('Имя и фамилия')).toBeInTheDocument()
  })

  it('navigates to the sign up route when the register tab is picked', async () => {
    renderWithShell(<AuthScreen mode="signIn" />)

    await userEvent.click(screen.getByRole('tab', { name: 'Регистрация' }))

    await waitFor(() => {
      expect(navigate).toHaveBeenCalledWith('/sign-up', { replace: true })
    })
  })

  it('navigates back to the sign in route when the login tab is picked', async () => {
    renderWithShell(<AuthScreen mode="signUp" />)

    await userEvent.click(screen.getByRole('tab', { name: 'Войти' }))

    await waitFor(() => {
      expect(navigate).toHaveBeenCalledWith('/sign-in', { replace: true })
    })
  })

  it('titles the card for the active tab', () => {
    renderWithShell(<AuthScreen mode="signUp" />)

    expect(
      screen.getByText('Создайте аккаунт, чтобы начать'),
    ).toBeInTheDocument()
  })

  it('exposes both tabs in a single tablist', () => {
    renderWithShell(<AuthScreen mode="signIn" />)

    expect(screen.getAllByRole('tab')).toHaveLength(2)
    expect(screen.getByRole('tablist')).toBeInTheDocument()
  })

  it('explains an expired session instead of leaving the screen bare', () => {
    renderWithShell(<AuthScreen mode="signIn" />, {
      session: { sessionExpired: true },
    })

    expect(screen.getByTestId('session-expired-notice')).toHaveTextContent(
      'Сессия истекла, войдите снова',
    )
  })

  it('stays quiet for a visitor who simply opened the sign in page', () => {
    renderWithShell(<AuthScreen mode="signIn" />)

    expect(
      screen.queryByTestId('session-expired-notice'),
    ).not.toBeInTheDocument()
  })

  it('renders a sliding indicator inside the tablist', () => {
    const { container } = renderWithShell(<AuthScreen mode="signIn" />)

    const indicator = container.querySelector('[data-slot="tabs-indicator"]')

    expect(indicator).not.toBeNull()
    expect(screen.getByRole('tablist')).toContainElement(
      indicator as HTMLElement,
    )
  })

  it('marks only the tab matching the current mode as active', () => {
    renderWithShell(<AuthScreen mode="signIn" />)

    expect(screen.getByRole('tab', { name: 'Войти' })).toHaveAttribute(
      'data-active',
    )
    expect(screen.getByRole('tab', { name: 'Регистрация' })).not.toHaveAttribute(
      'data-active',
    )

    cleanup()
    renderWithShell(<AuthScreen mode="signUp" />)

    expect(screen.getByRole('tab', { name: 'Регистрация' })).toHaveAttribute(
      'data-active',
    )
    expect(screen.getByRole('tab', { name: 'Войти' })).not.toHaveAttribute(
      'data-active',
    )
  })

  it('keeps the panel height unconstrained when motion is reduced', () => {
    matchMediaResult = true

    const { container } = renderWithShell(<AuthScreen mode="signIn" />)

    const panels = container.querySelector<HTMLElement>('.auth-card__panels')

    expect(panels).not.toBeNull()
    expect(panels?.style.height).toBe('')
  })
})
