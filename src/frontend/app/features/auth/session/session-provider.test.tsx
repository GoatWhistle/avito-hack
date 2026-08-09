import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen, waitFor } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import {
  afterAll,
  afterEach,
  beforeAll,
  beforeEach,
  describe,
  expect,
  it,
} from 'vitest'
import { apiBaseUrl, httpClient, resetTokenCache, setToken } from '#/api'
import { SessionProvider } from './SessionProvider'
import { useSession } from './session-context'

const url = (path: string) => `${apiBaseUrl}${path}`

const unauthorized = () =>
  HttpResponse.json(
    { error: { code: 'unauthorized', message: 'token expired' } },
    { status: 401 },
  )

const server = setupServer()

function Probe() {
  const { sessionExpired, isAuthenticated } = useSession()

  return (
    <p data-testid="probe">
      {sessionExpired ? 'expired' : 'active'}:
      {isAuthenticated ? 'authenticated' : 'guest'}
    </p>
  )
}

const renderProvider = () => {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  })

  return render(
    <QueryClientProvider client={queryClient}>
      <SessionProvider>
        <Probe />
      </SessionProvider>
    </QueryClientProvider>,
  )
}

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))
afterAll(() => server.close())

beforeEach(() => {
  window.localStorage.clear()
  resetTokenCache()
})

afterEach(() => {
  server.resetHandlers()
  window.localStorage.clear()
  resetTokenCache()
})

describe('SessionProvider', () => {
  it('flags an expiry once the refresh is refused', async () => {
    setToken('stale-token')
    server.use(
      http.get(url('/users/me'), () => unauthorized()),
      http.post(url('/auth/refresh'), () => unauthorized()),
    )

    renderProvider()

    await waitFor(() => {
      expect(screen.getByTestId('probe')).toHaveTextContent('expired:guest')
    })
  })

  it('keeps the session alive when the refresh succeeds', async () => {
    setToken('stale-token')
    let meCalls = 0
    server.use(
      http.get(url('/users/me'), () => {
        meCalls += 1
        return meCalls === 1
          ? unauthorized()
          : HttpResponse.json({
              id: 'e6f0d1c2-0000-4000-a000-000000000001',
              email: 'anna@demo.avito',
              full_name: 'Анна',
              role: 'user',
              created_at: '2026-08-01T10:00:00Z',
            })
      }),
      http.post(url('/auth/refresh'), () =>
        HttpResponse.json({ token: 'fresh-token' }),
      ),
    )

    renderProvider()

    await waitFor(() => {
      expect(screen.getByTestId('probe')).toHaveTextContent(
        'active:authenticated',
      )
    })
  })

  it('does not flag an expiry for a visitor who never signed in', async () => {
    renderProvider()

    await waitFor(() => {
      expect(screen.getByTestId('probe')).toHaveTextContent('active:guest')
    })
  })

  it('clears the expiry flag once a new token arrives', async () => {
    setToken('stale-token')
    server.use(
      http.get(url('/users/me'), () => unauthorized()),
      http.post(url('/auth/refresh'), () => unauthorized()),
    )

    renderProvider()

    await waitFor(() => {
      expect(screen.getByTestId('probe')).toHaveTextContent('expired:guest')
    })

    setToken('signed-in-again')

    await waitFor(() => {
      expect(screen.getByTestId('probe')).toHaveTextContent('active:')
    })
  })

  it('leaves the token in place when an unrelated call is forbidden', async () => {
    setToken('valid-token')
    server.use(
      http.get(url('/users/me'), () =>
        HttpResponse.json({
          id: 'e6f0d1c2-0000-4000-a000-000000000001',
          email: 'anna@demo.avito',
          full_name: 'Анна',
          role: 'user',
          created_at: '2026-08-01T10:00:00Z',
        }),
      ),
      http.get(url('/rewards'), () =>
        HttpResponse.json(
          { error: { code: 'forbidden', message: 'nope' } },
          { status: 403 },
        ),
      ),
    )

    renderProvider()

    await waitFor(() => {
      expect(screen.getByTestId('probe')).toHaveTextContent(
        'active:authenticated',
      )
    })

    await expect(httpClient.get('/rewards')).rejects.toMatchObject({
      kind: 'forbidden',
    })

    expect(screen.getByTestId('probe')).toHaveTextContent(
      'active:authenticated',
    )
  })
})
