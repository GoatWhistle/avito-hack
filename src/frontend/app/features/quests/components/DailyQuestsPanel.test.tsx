import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { renderWithProviders } from '#/features/items/components/test-utils'
import type { QuestItem } from '#/features/quests/types'
import { DailyQuestsPanel } from './DailyQuestsPanel'

const today = vi.fn()
const claim = vi.fn()

vi.mock('#/features/quests/repository', () => ({
  questsRepository: {
    today: (signal?: AbortSignal) => today(signal),
    claim: () => claim(),
  },
}))

const makeQuest = (overrides: Partial<QuestItem> = {}): QuestItem => ({
  id: 'favorite_three',
  action: 'favorite',
  target: 3,
  reward_xp: 5,
  progress_current: 0,
  completed: false,
  claimed: false,
  ...overrides,
})

beforeEach(() => {
  vi.clearAllMocks()
  today.mockResolvedValue([makeQuest()])
  claim.mockResolvedValue([])
})

describe('DailyQuestsPanel', () => {
  it('renders each quest with its progress', async () => {
    today.mockResolvedValue([
      makeQuest({ progress_current: 1 }),
      makeQuest({
        id: 'sell_one',
        action: 'item_sold',
        target: 1,
        reward_xp: 15,
      }),
    ])

    renderWithProviders(<DailyQuestsPanel />)

    const rows = await screen.findAllByTestId('quest-row')
    expect(rows).toHaveLength(2)
    expect(screen.getByText('1 / 3')).toBeInTheDocument()
  })

  it('shows the human-readable quest goal', async () => {
    renderWithProviders(<DailyQuestsPanel />)

    expect(
      await screen.findByText('Добавьте 3 объявления в избранное'),
    ).toBeInTheDocument()
  })

  it('marks a finished quest as completed', async () => {
    today.mockResolvedValue([
      makeQuest({ progress_current: 3, completed: true }),
    ])

    renderWithProviders(<DailyQuestsPanel />)

    await waitFor(() => {
      expect(screen.getByTestId('quest-row')).toHaveAttribute(
        'data-completed',
        'true',
      )
    })
  })

  it('offers claiming only when a reward is pending', async () => {
    today.mockResolvedValue([makeQuest({ progress_current: 1 })])

    renderWithProviders(<DailyQuestsPanel />)

    await screen.findByTestId('quest-row')
    expect(screen.queryByTestId('claim-quests')).not.toBeInTheDocument()
  })

  it('claims pending rewards on demand', async () => {
    today.mockResolvedValue([
      makeQuest({ progress_current: 3, completed: true }),
    ])
    claim.mockResolvedValue([
      makeQuest({ progress_current: 3, completed: true, claimed: true }),
    ])

    renderWithProviders(<DailyQuestsPanel />)

    await userEvent.click(await screen.findByTestId('claim-quests'))

    await waitFor(() => {
      expect(claim).toHaveBeenCalledTimes(1)
    })
    await waitFor(() => {
      expect(screen.queryByTestId('claim-quests')).not.toBeInTheDocument()
    })
  })

  it('reports a loading failure', async () => {
    today.mockRejectedValue(new Error('boom'))

    renderWithProviders(<DailyQuestsPanel />)

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Не удалось загрузить задания',
    )
  })
})
