import type { ReactElement } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { I18nextProvider } from 'react-i18next'
import { createMemoryRouter, RouterProvider } from 'react-router'
import { render, type RenderResult } from '@testing-library/react'
import { initI18n } from '#/i18n'
import { SessionContext, type SessionValue } from '#/features/auth/session'
import type {
  FavoriteEntry,
  Item,
  ItemListEntry,
  ItemPhoto,
} from '#/features/items/types'

if (typeof window !== 'undefined') {
  window.localStorage.setItem('avito-hack.locale', 'ru')
}

export const makeItem = (overrides: Partial<Item> = {}): Item => ({
  id: 'item-1',
  owner_id: 'user-1',
  title: 'Велосипед Stels',
  description: 'Отличное состояние',
  price: 1_250_000,
  status: 'published',
  attributes: null,
  created_at: '2026-08-01T10:00:00Z',
  updated_at: '2026-08-01T10:00:00Z',
  ...overrides,
})

export const makeEntry = (
  overrides: Partial<ItemListEntry> = {},
): ItemListEntry => ({
  id: 'item-1',
  owner_id: 'user-1',
  title: 'Велосипед Stels',
  price: 1_250_000,
  status: 'published',
  created_at: '2026-08-01T10:00:00Z',
  ...overrides,
})

export const makePhoto = (overrides: Partial<ItemPhoto> = {}): ItemPhoto => ({
  id: 'photo-1',
  item_id: 'item-1',
  url: 'https://cdn.test/photo-1.jpg',
  position: 0,
  created_at: '2026-08-01T10:00:00Z',
  ...overrides,
})

export const makeFavorite = (
  overrides: Partial<FavoriteEntry> = {},
): FavoriteEntry => ({
  item_id: 'item-1',
  owner_id: 'user-2',
  title: 'Велосипед Stels',
  price: 1_250_000,
  status: 'published',
  added_at: '2026-08-02T10:00:00Z',
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
  signOut: () => {},
  setUser: () => {},
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
