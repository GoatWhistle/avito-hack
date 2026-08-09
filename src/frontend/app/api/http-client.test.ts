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
import { ApiError } from './api-error'
import { apiBaseUrl, httpClient } from './http-client.api'
import { getToken, resetTokenCache, setToken } from './token-store'

const url = (path: string) => `${apiBaseUrl}${path}`

const unauthorized = () =>
  HttpResponse.json(
    { error: { code: 'unauthorized', message: 'token expired' } },
    { status: 401 },
  )

const server = setupServer()

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))
afterAll(() => server.close())

beforeEach(() => {
  window.localStorage.clear()
  resetTokenCache()
  setToken('stale-token')
})

afterEach(() => {
  server.resetHandlers()
  window.localStorage.clear()
  resetTokenCache()
})

describe('401 from a non-auth endpoint', () => {
  it('refreshes the session and keeps the user signed in', async () => {
    let petCalls = 0
    server.use(
      http.get(url('/pet'), () => {
        petCalls += 1
        return petCalls === 1
          ? unauthorized()
          : HttpResponse.json({ stage: 'egg' })
      }),
      http.post(url('/auth/refresh'), () =>
        HttpResponse.json({ token: 'fresh-token' }),
      ),
    )

    const response = await httpClient.get('/pet')

    expect(response.data).toEqual({ stage: 'egg' })
    expect(getToken()).toBe('fresh-token')
  })

  it('replays the original request with the refreshed token', async () => {
    const seen: (string | null)[] = []
    let calls = 0
    server.use(
      http.get(url('/rewards'), ({ request }) => {
        seen.push(request.headers.get('authorization'))
        calls += 1
        return calls === 1 ? unauthorized() : HttpResponse.json([])
      }),
      http.post(url('/auth/refresh'), () =>
        HttpResponse.json({ token: 'fresh-token' }),
      ),
    )

    await httpClient.get('/rewards')

    expect(seen).toEqual(['Bearer stale-token', 'Bearer fresh-token'])
  })

  it('keeps the token when the replayed request fails for another reason', async () => {
    let calls = 0
    server.use(
      http.get(url('/pet'), () => {
        calls += 1
        return calls === 1
          ? unauthorized()
          : HttpResponse.json(
              { error: { code: 'forbidden', message: 'nope' } },
              { status: 403 },
            )
      }),
      http.post(url('/auth/refresh'), () =>
        HttpResponse.json({ token: 'fresh-token' }),
      ),
    )

    await expect(httpClient.get('/pet')).rejects.toMatchObject({
      kind: 'forbidden',
    })
    expect(getToken()).toBe('fresh-token')
  })

  it('clears the token only once the refresh itself is rejected', async () => {
    server.use(
      http.get(url('/pet'), () => unauthorized()),
      http.post(url('/auth/refresh'), () => unauthorized()),
    )

    await expect(httpClient.get('/pet')).rejects.toMatchObject({
      kind: 'unauthorized',
    })
    expect(getToken()).toBeNull()
  })

  it('refreshes once for concurrent failures', async () => {
    let refreshes = 0
    let petCalls = 0
    let rewardCalls = 0
    server.use(
      http.get(url('/pet'), () => {
        petCalls += 1
        return petCalls === 1 ? unauthorized() : HttpResponse.json({})
      }),
      http.get(url('/rewards'), () => {
        rewardCalls += 1
        return rewardCalls === 1 ? unauthorized() : HttpResponse.json([])
      }),
      http.post(url('/auth/refresh'), () => {
        refreshes += 1
        return HttpResponse.json({ token: 'fresh-token' })
      }),
    )

    await Promise.all([httpClient.get('/pet'), httpClient.get('/rewards')])

    expect(refreshes).toBe(1)
  })

  it('retries at most once instead of looping', async () => {
    let refreshes = 0
    let petCalls = 0
    server.use(
      http.get(url('/pet'), () => {
        petCalls += 1
        return unauthorized()
      }),
      http.post(url('/auth/refresh'), () => {
        refreshes += 1
        return HttpResponse.json({ token: 'fresh-token' })
      }),
    )

    await expect(httpClient.get('/pet')).rejects.toBeInstanceOf(ApiError)

    expect(refreshes).toBe(1)
    expect(petCalls).toBe(2)
  })
})

describe('401 from an auth endpoint', () => {
  it('never refreshes when the credentials are rejected', async () => {
    let refreshes = 0
    server.use(
      http.post(url('/auth/login'), () => unauthorized()),
      http.post(url('/auth/refresh'), () => {
        refreshes += 1
        return HttpResponse.json({ token: 'fresh-token' })
      }),
    )

    await expect(
      httpClient.post('/auth/login', { email: 'a@b.c', password: 'nope' }),
    ).rejects.toMatchObject({ kind: 'unauthorized' })

    expect(refreshes).toBe(0)
  })

  it('leaves the existing token untouched when the login fails', async () => {
    server.use(http.post(url('/auth/login'), () => unauthorized()))

    await expect(
      httpClient.post('/auth/login', { email: 'a@b.c', password: 'nope' }),
    ).rejects.toBeInstanceOf(ApiError)

    expect(getToken()).toBe('stale-token')
  })
})

describe('requests made without a token', () => {
  it('rejects a 401 without touching the refresh endpoint', async () => {
    setToken(null)
    let refreshes = 0
    server.use(
      http.get(url('/pet'), () => unauthorized()),
      http.post(url('/auth/refresh'), () => {
        refreshes += 1
        return HttpResponse.json({ token: 'fresh-token' })
      }),
    )

    await expect(httpClient.get('/pet')).rejects.toMatchObject({
      kind: 'unauthorized',
    })

    expect(refreshes).toBe(0)
  })
})
