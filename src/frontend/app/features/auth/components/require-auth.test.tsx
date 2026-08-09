import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen } from '@testing-library/react'
import { I18nextProvider } from 'react-i18next'
import { createMemoryRouter, RouterProvider, useLocation } from 'react-router'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { resetTokenCache, setToken } from '#/api'
import { SessionContext, type SessionValue } from '#/features/auth/session'
import { initI18n } from '#/i18n'
import { RequireAuth } from './RequireAuth'

const makeSession = (overrides: Partial<SessionValue> = {}): SessionValue => ({
  user: null,
  isAuthenticated: false,
  isLoading: false,
  sessionExpired: false,
  signOut: () => {},
  setUser: () => {},
  acknowledgeExpiry: () => {},
  ...overrides,
})

function SignInProbe() {
  const location = useLocation()
  const state = location.state as { from?: string } | null

  return <p data-testid="sign-in">from:{state?.from ?? 'none'}</p>
}

const renderAt = (path: string, session: Partial<SessionValue>) => {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  })

  const router = createMemoryRouter(
    [
      {
        path: '/pet',
        element: (
          <RequireAuth>
            <p data-testid="pet">Питомец</p>
          </RequireAuth>
        ),
      },
      { path: '/sign-in', element: <SignInProbe /> },
    ],
    { initialEntries: [path] },
  )

  return render(
    <I18nextProvider i18n={initI18n()}>
      <QueryClientProvider client={queryClient}>
        <SessionContext value={makeSession(session)}>
          <RouterProvider router={router} />
        </SessionContext>
      </QueryClientProvider>
    </I18nextProvider>,
  )
}

beforeEach(() => {
  window.localStorage.clear()
  resetTokenCache()
})

afterEach(() => {
  vi.clearAllMocks()
  window.localStorage.clear()
  resetTokenCache()
})

describe('RequireAuth on a private route', () => {
  it('sends a guest to the sign in screen instead of rendering the page', () => {
    renderAt('/pet', { isAuthenticated: false })

    expect(screen.getByTestId('sign-in')).toBeInTheDocument()
    expect(screen.queryByTestId('pet')).not.toBeInTheDocument()
  })

  it('leaves any stored token alone while redirecting', () => {
    setToken('guest-visit-token')

    renderAt('/pet', { isAuthenticated: false })

    expect(window.localStorage.getItem('avito-hack.token')).toBe(
      'guest-visit-token',
    )
  })

  it('remembers where the guest was heading', () => {
    renderAt('/pet', { isAuthenticated: false })

    expect(screen.getByTestId('sign-in')).toHaveTextContent('from:/pet')
  })

  it('waits instead of redirecting while the session is still loading', () => {
    renderAt('/pet', { isLoading: true })

    expect(screen.getByRole('status')).toBeInTheDocument()
    expect(screen.queryByTestId('sign-in')).not.toBeInTheDocument()
    expect(screen.queryByTestId('pet')).not.toBeInTheDocument()
  })

  it('renders the page for an authenticated visitor', () => {
    renderAt('/pet', { isAuthenticated: true })

    expect(screen.getByTestId('pet')).toBeInTheDocument()
  })
})
