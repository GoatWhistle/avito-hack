import type { ReactElement } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { I18nextProvider } from 'react-i18next'
import { createMemoryRouter, RouterProvider } from 'react-router'
import { render, type RenderResult } from '@testing-library/react'
import { initI18n } from '#/i18n'
import { SessionContext, type SessionValue } from '#/features/auth/session'
import type {
  GameList,
  GameState,
  GameStreak,
  GameSummary,
  MoreLessPrompt,
} from '#/features/games/types'

if (typeof window !== 'undefined') {
  window.localStorage.setItem('avito-hack.locale', 'ru')
}

export const makeStreak = (overrides: Partial<GameStreak> = {}): GameStreak => ({
  current_days: 0,
  best_days: 0,
  reward_ready: false,
  ...overrides,
})

export const makeSummary = (
  overrides: Partial<GameSummary> = {},
): GameSummary => ({
  slug: 'moreless',
  target_streak: 7,
  daily_done: false,
  ...overrides,
})

export const makeGameList = (overrides: Partial<GameList> = {}): GameList => ({
  games: [makeSummary()],
  streak: makeStreak(),
  daily_done: false,
  ...overrides,
})

export const makeGameState = (
  overrides: Partial<GameState> = {},
): GameState => ({
  slug: 'moreless',
  target_streak: 7,
  streak: makeStreak(),
  daily: { attempts: 0, best_streak: 0 },
  ...overrides,
})

export const makePrompt = (
  overrides: Partial<MoreLessPrompt> = {},
): MoreLessPrompt => ({
  left: {
    item_id: 'aaa111222333',
    title: 'Велосипед Stels',
    photo_url: '/uploads/aaa111222333/demo-1.jpg',
    price: 1_250_000,
  },
  right: {
    item_id: 'bbb444555666',
    title: 'Диван угловой',
    photo_url: '/uploads/bbb444555666/demo-1.jpg',
  },
  ...overrides,
})

const makeSession = (userId: string | null): SessionValue => ({
  user: userId
    ? {
        id: userId,
        email: 'me@test.dev',
        fullName: 'Я',
        role: 'user',
        createdAt: '2026-08-01T10:00:00Z',
      }
    : null,
  isAuthenticated: Boolean(userId),
  isLoading: false,
  sessionExpired: false,
  signOut: () => {},
  setUser: () => {},
  acknowledgeExpiry: () => {},
})

interface RenderOptions {
  userId?: string | null
  route?: string
  path?: string
}

export const renderWithProviders = (
  ui: ReactElement,
  { userId = 'user-1', route = '/', path = '/' }: RenderOptions = {},
): RenderResult => {
  const instance = initI18n()
  if (instance.language !== 'ru') void instance.changeLanguage('ru')

  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  })
  const router = createMemoryRouter(
    [
      { path, element: ui },
      { path: '*', element: <div /> },
    ],
    { initialEntries: [route] },
  )

  return render(
    <I18nextProvider i18n={initI18n()}>
      <QueryClientProvider client={queryClient}>
        <SessionContext value={makeSession(userId)}>
          <RouterProvider router={router} />
        </SessionContext>
      </QueryClientProvider>
    </I18nextProvider>,
  )
}
