import type { AxiosInstance } from 'axios'
import { getToken, setToken } from './token-store'

const REFRESH_PATH = '/auth/refresh'

const AUTH_PATHS = ['/auth/login', '/auth/register', REFRESH_PATH]

type SessionPayload = { token?: unknown }

export const isAuthPath = (url: string | undefined): boolean => {
  if (!url) return false
  const path = url.split('?')[0].replace(/\/+$/, '')
  return AUTH_PATHS.some((candidate) => path.endsWith(candidate))
}

export class SessionRefresher {
  private pending: Promise<string | null> | null = null

  constructor(private readonly client: AxiosInstance) {}

  refresh(): Promise<string | null> {
    this.pending ??= this.run().finally(() => {
      this.pending = null
    })

    return this.pending
  }

  private async run(): Promise<string | null> {
    const current = getToken()
    if (!current) return null

    try {
      const response = await this.client.post<SessionPayload>(
        REFRESH_PATH,
        null,
        { headers: { Authorization: `Bearer ${current}` } },
      )

      const token = response.data?.token
      if (typeof token !== 'string' || token.length === 0) return null

      setToken(token)

      return token
    } catch {
      return null
    }
  }
}
