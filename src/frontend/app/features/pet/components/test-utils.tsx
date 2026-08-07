import type { ReactElement } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { I18nextProvider } from 'react-i18next'
import { createMemoryRouter, RouterProvider } from 'react-router'
import { render, type RenderResult } from '@testing-library/react'
import { initI18n } from '#/i18n'
import type { Pet } from '#/features/pet/types'

export const makePet = (overrides: Partial<Pet> = {}): Pet => ({
  id: 'pet-1',
  user_id: 'user-1',
  name: 'Ноти',
  stage: 'baby',
  state: 'neutral',
  level: 3,
  xp: 15,
  next_level_xp: 22,
  satiety: 70,
  happiness: 80,
  energy: 90,
  streak_days: 4,
  freezes: 1,
  is_hatched: true,
  hatched_at: '2026-08-01T10:00:00Z',
  last_checkin_date: null,
  ...overrides,
})

export const renderWithProviders = (ui: ReactElement): RenderResult => {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  })
  const router = createMemoryRouter([{ path: '/', element: ui }], {
    initialEntries: ['/'],
  })

  return render(
    <I18nextProvider i18n={initI18n()}>
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    </I18nextProvider>,
  )
}
