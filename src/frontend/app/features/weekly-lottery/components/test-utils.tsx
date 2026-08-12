import type { ReactElement } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render } from '@testing-library/react'
import { I18nextProvider } from 'react-i18next'
import { createMemoryRouter, RouterProvider } from 'react-router'
import { initI18n } from '#/i18n'
import type {
  LotteryPrizeCatalogItem,
  LotteryRun,
  LotterySlot,
  LotteryState,
} from '#/features/weekly-lottery/types'

if (typeof window !== 'undefined') {
  window.localStorage.setItem('avito-hack.locale', 'ru')
}

export const closedSlots = (): LotterySlot[] =>
  Array.from(
    { length: 9 },
    (_, index): LotterySlot => ({ index, opened: false }),
  )

export const makeRun = (overrides: Partial<LotteryRun> = {}): LotteryRun => ({
  id: 'run123456789',
  state: 'active',
  slots: closedSlots(),
  created_at: '2026-08-12T12:00:00Z',
  ...overrides,
})

export const makeState = (
  overrides: Partial<LotteryState> = {},
): LotteryState => ({
  available: true,
  week_start: '2026-08-10T00:00:00+03:00',
  next_available_at: '2026-08-17T00:00:00+03:00',
  ...overrides,
})

export const makePrizes = (): LotteryPrizeCatalogItem[] => [
  {
    id: 'weekly_bicycle_5',
    symbol: 'bicycle',
    title: 'Скидка 5% на спорт и отдых',
    description: 'Скидка на товар из категории спорта и отдыха',
    benefit_type: 'percent_discount',
    benefit_value: 5,
    scope_type: 'category',
    scope_value: 'sport',
  },
]

export const renderLottery = (ui: ReactElement) => {
  const i18n = initI18n()
  if (i18n.language !== 'ru') void i18n.changeLanguage('ru')

  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  })
  const router = createMemoryRouter(
    [
      { path: '/play/weekly-lottery', element: ui },
      { path: '*', element: <div /> },
    ],
    { initialEntries: ['/play/weekly-lottery'] },
  )

  return render(
    <I18nextProvider i18n={i18n}>
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    </I18nextProvider>,
  )
}
