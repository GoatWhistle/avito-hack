import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { PetScreen } from './PetScreen'
import { makePet, renderWithProviders } from './test-utils'
import type { DailySummary } from '#/features/pet/types'

const state = vi.fn()
const summaryToday = vi.fn()

vi.mock('#/features/pet/repository', () => ({
  petRepository: {
    state: () => state(),
    stroke: () => Promise.resolve(makePet()),
    checkIn: () => Promise.reject(new Error('unused')),
    summaryToday: () => summaryToday(),
  },
}))

const summary: DailySummary = {
  id: 'sum-1',
  date: '2026-08-07',
  message: 'Сегодня мы отлично поработали — три объявления стали лучше!',
  advice: {
    text: 'Добавьте фото в объявление о велосипеде',
    item_id: 'item-42',
    action: 'add_photo',
  },
  generated_by: 'llm',
  facts: {
    total_xp: 34,
    actions: [],
    level: 3,
    previous_level: 3,
    leveled_up: false,
    xp: 15,
    next_level_xp: 22,
    xp_to_next_level: 7,
    rewards: [],
    badges: [],
    stage: 'baby',
    state: 'neutral',
    satiety: 70,
    happiness: 80,
    energy: 90,
    streak_days: 4,
    streak_broken: false,
    issues_count: 1,
  },
}

const noSocket = { events: { enabled: false as const } }

beforeEach(() => {
  vi.clearAllMocks()
  state.mockResolvedValue(makePet())
})

describe('daily summary', () => {
  it('renders the message, facts and a link to the item', async () => {
    summaryToday.mockResolvedValue(summary)
    renderWithProviders(<PetScreen {...noSocket} />)

    const card = await screen.findByTestId('daily-summary')
    expect(card).toHaveTextContent('три объявления стали лучше')
    expect(card).toHaveTextContent('+34')
    expect(card).toHaveTextContent('Добавьте фото')

    expect(screen.getByRole('link', { name: /исправить/i })).toHaveAttribute(
      'href',
      '/items/item-42',
    )
  })

  it('omits the fix link when the advice has no item', async () => {
    summaryToday.mockResolvedValue({
      ...summary,
      advice: { text: 'Отметьтесь завтра', item_id: null, action: 'check_in' },
    })
    renderWithProviders(<PetScreen {...noSocket} />)

    await screen.findByTestId('daily-summary')
    expect(screen.queryByRole('link', { name: /исправить/i })).toBeNull()
  })

  it('can be dismissed', async () => {
    summaryToday.mockResolvedValue(summary)
    renderWithProviders(<PetScreen {...noSocket} />)

    const card = await screen.findByTestId('daily-summary')
    await userEvent.click(
      screen.getAllByRole('button', { name: /закрыть/i })[0],
    )

    expect(card).not.toBeInTheDocument()
  })

  it('renders nothing when there is no summary for today', async () => {
    summaryToday.mockResolvedValue(null)
    renderWithProviders(<PetScreen {...noSocket} />)

    await screen.findByRole('heading', { name: 'Ноти', level: 1 })
    expect(screen.queryByTestId('daily-summary')).toBeNull()
  })
})
