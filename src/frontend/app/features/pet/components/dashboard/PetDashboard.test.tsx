import { screen, within } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import type * as RewardsModule from '#/features/rewards'
import {
  makeCheckedInPet,
  makePet,
  renderWithProviders,
} from '../test-utils'
import { PetDashboard } from './PetDashboard'

const state = vi.fn()
const summaryToday = vi.fn()

vi.mock('#/features/pet/repository', () => ({
  petRepository: {
    state: () => state(),
    stroke: () => Promise.resolve(makePet()),
    feed: () => Promise.resolve(makePet()),
    checkIn: () => Promise.reject(new Error('unused')),
    summaryToday: () => summaryToday(),
  },
}))

vi.mock('#/features/rewards', async (importOriginal) => {
  const actual = await importOriginal<typeof RewardsModule>()

  return { ...actual, useBadges: () => ({ data: [{ id: 'b1' }] }) }
})

const noSocket = { events: { enabled: false as const } }

const renderDashboard = (children = 'panel') =>
  renderWithProviders(<PetDashboard {...noSocket}>{children}</PetDashboard>)

beforeEach(() => {
  vi.clearAllMocks()
  state.mockResolvedValue(makePet())
  summaryToday.mockResolvedValue(null)
})

describe('PetDashboard', () => {
  it('renders the pet card, menu and stat bars in the left sidebar', async () => {
    renderDashboard()
    await screen.findByRole('heading', { name: 'Ноти' })

    expect(screen.getByText('Малыш · ур. 3')).toBeInTheDocument()

    const nav = screen.getAllByRole('navigation', {
      name: 'Разделы питомца',
    })[0]
    expect(within(nav).getByRole('link', { name: /Дом/ })).toBeInTheDocument()
    expect(
      within(nav).getByRole('link', { name: /Таблица лидеров/ }),
    ).toBeInTheDocument()

    const stats = screen.getByRole('list', { name: 'Параметры питомца' })
    for (const label of ['Энергия', 'Счастье', 'Сытость']) {
      expect(within(stats).getByText(label)).toBeInTheDocument()
    }
    expect(within(stats).queryByText('Стрик (дней)')).not.toBeInTheDocument()
  })

  it('keeps xp, badges and the streak inside the pet card', async () => {
    renderDashboard()
    await screen.findByRole('heading', { name: 'Ноти' })

    const card = screen.getByTestId('pet-card')
    expect(within(card).getByText('15 / 22 XP')).toBeInTheDocument()
    expect(within(card).getByText('1 бейджей')).toBeInTheDocument()
    expect(within(card).getByText('4 дн.')).toBeInTheDocument()
    expect(
      within(card).getByRole('progressbar', {
        name: 'Прогресс до уровня 4, осталось 7 XP',
      }),
    ).toBeInTheDocument()
  })

  it('drops the standalone top hud above the raccoon', async () => {
    renderDashboard()
    await screen.findByRole('heading', { name: 'Ноти' })

    expect(
      screen.queryByText('15 XP · до ур. 4: 7 XP'),
    ).not.toBeInTheDocument()
    expect(
      screen.queryByRole('heading', { name: 'Ноти · ур. 3' }),
    ).not.toBeInTheDocument()
  })

  it('greets with the pet name and stage', async () => {
    renderDashboard()

    expect(await screen.findByTestId('pet-speech')).toHaveTextContent(
      'Привет! Я Ноти. Стадия: Малыш',
    )
  })

  it('renders the routed panel content on the right', async () => {
    renderDashboard('panel content')
    await screen.findByRole('heading', { name: 'Ноти' })

    const panel = screen.getByRole('complementary', { name: 'Раздел питомца' })
    expect(within(panel).getByText('panel content')).toBeInTheDocument()
  })

  it('offers neither a check-in nor a separate stroke button', async () => {
    renderDashboard()
    await screen.findByRole('heading', { name: 'Ноти' })

    expect(
      screen.queryByRole('button', { name: 'Отметиться' }),
    ).not.toBeInTheDocument()
    expect(
      screen.queryByRole('button', { name: 'Погладить' }),
    ).not.toBeInTheDocument()
  })

  it('shows the streak count once, with the current streak week', async () => {
    state.mockResolvedValue(
      makePet({ last_checkin_date: new Date().toISOString(), streak_days: 4 }),
    )
    renderDashboard()
    await screen.findByRole('heading', { name: 'Ноти' })

    const tracker = screen.getByRole('region', { name: 'Стрик' })
    expect(within(tracker).getByText('4')).toBeInTheDocument()
    expect(tracker).toHaveTextContent('1-я неделя серии')
    expect(tracker).not.toHaveTextContent('следующая отметка завтра')
  })

  it('fills only the days of the current streak week', async () => {
    state.mockResolvedValue(
      makePet({ last_checkin_date: new Date().toISOString(), streak_days: 41 }),
    )
    renderDashboard()
    await screen.findByRole('heading', { name: 'Ноти' })

    const tracker = screen.getByRole('region', { name: 'Стрик' })
    expect(tracker).toHaveTextContent('6-я неделя серии')
    expect(
      within(tracker).getAllByLabelText(/отмечено сегодня$/),
    ).toHaveLength(6)
  })

  it('banners the streak when the backend applied a check-in', async () => {
    state.mockResolvedValue(makeCheckedInPet({ streak_days: 5 }))
    renderDashboard()

    expect(await screen.findByTestId('celebration-banner')).toHaveTextContent(
      'Серия 5 дней · +10 XP за сегодня',
    )
  })

  it('offers the feed button next to the satiety stat', async () => {
    renderDashboard()
    await screen.findByRole('heading', { name: 'Ноти' })

    const stats = screen.getByRole('list', { name: 'Параметры питомца' })
    expect(within(stats).getByTestId('feed-button')).toBeInTheDocument()
  })

  it('marks the streak tracker as done for today', async () => {
    state.mockResolvedValue(
      makePet({ last_checkin_date: new Date().toISOString(), streak_days: 4 }),
    )
    renderDashboard()
    await screen.findByRole('heading', { name: 'Ноти' })

    const tracker = screen.getByRole('region', { name: 'Стрик' })
    expect(
      within(tracker).getByLabelText('День 4, отмечено сегодня'),
    ).toBeInTheDocument()
    expect(
      within(tracker).queryByRole('listitem', { current: 'step' }),
    ).not.toBeInTheDocument()
  })

  it('lets the raccoon itself be stroked', async () => {
    renderDashboard()
    await screen.findByRole('heading', { name: 'Ноти' })

    expect(
      screen.getByRole('button', { name: /нажмите, чтобы погладить/i }),
    ).toBeInTheDocument()
  })

  it('renders one week of streak days and marks today as the target', async () => {
    renderDashboard()
    await screen.findByRole('heading', { name: 'Ноти' })

    const tracker = screen.getByRole('region', { name: 'Стрик' })
    expect(within(tracker).getAllByRole('listitem')).toHaveLength(7)
    expect(
      within(tracker).getByRole('listitem', { current: 'step' }),
    ).toHaveAccessibleName('День 5')
  })

  it('keeps the profile out of the dashboard navigation', async () => {
    renderDashboard()
    await screen.findByRole('heading', { name: 'Ноти' })

    for (const nav of screen.getAllByRole('navigation', {
      name: 'Разделы питомца',
    })) {
      expect(
        within(nav).queryByRole('link', { name: /Профиль/ }),
      ).not.toBeInTheDocument()
    }

    expect(
      document.querySelector('a[href="/pet/profile"]'),
    ).not.toBeInTheDocument()
  })

  it('links back to the avito listings', async () => {
    renderDashboard()

    expect(
      await screen.findByRole('link', { name: /К объявлениям/ }),
    ).toHaveAttribute('href', '/items')
  })

  it('shows the dashboard right away for a freshly created baby', async () => {
    state.mockResolvedValue(makePet({ stage: 'baby' }))
    renderDashboard()

    expect(await screen.findByTestId('pet-speech')).toHaveTextContent(
      'Стадия: Малыш',
    )
  })
})
