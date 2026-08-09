import type { ReactElement } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { I18nextProvider } from 'react-i18next'
import { createMemoryRouter, RouterProvider } from 'react-router'
import { render, type RenderResult } from '@testing-library/react'
import { initI18n } from '#/i18n'
import { SessionContext, type SessionValue } from '#/features/auth/session'
import { ThemeProvider } from '#/features/layout/theme'
import type { User } from '#/types'

export const makeUser = (overrides: Partial<User> = {}): User => ({
  id: 'user-1',
  email: 'demo@example.com',
  fullName: 'Иван Иванов',
  role: 'user',
  createdAt: '2026-01-15T09:00:00Z',
  ...overrides,
})

interface RenderOptions {
  session?: Partial<SessionValue>
  initialEntries?: string[]
}

export const makeSession = (
  overrides: Partial<SessionValue> = {},
): SessionValue => ({
  user: null,
  isAuthenticated: false,
  isLoading: false,
  sessionExpired: false,
  signOut: () => {},
  setUser: () => {},
  acknowledgeExpiry: () => {},
  ...overrides,
})

export const renderWithShell = (
  ui: ReactElement,
  { session, initialEntries = ['/'] }: RenderOptions = {},
): RenderResult => {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  })
  const router = createMemoryRouter([{ path: '*', element: ui }], {
    initialEntries,
  })

  return render(
    <I18nextProvider i18n={initI18n()}>
      <ThemeProvider>
        <QueryClientProvider client={queryClient}>
          <SessionContext value={makeSession(session)}>
            <RouterProvider router={router} />
          </SessionContext>
        </QueryClientProvider>
      </ThemeProvider>
    </I18nextProvider>,
  )
}
