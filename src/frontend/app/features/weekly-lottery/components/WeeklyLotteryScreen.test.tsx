import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { WeeklyLotteryScreen } from './WeeklyLotteryScreen'
import {
  makePrizes,
  makeRun,
  makeState,
  renderLottery,
} from './test-utils'

const state = vi.fn()
const prizes = vi.fn()
const start = vi.fn()
const reveal = vi.fn()

vi.mock('#/features/weekly-lottery/repository', () => ({
  weeklyLotteryRepository: {
    state: (signal?: AbortSignal) => state(signal),
    prizes: (signal?: AbortSignal) => prizes(signal),
    start: () => start(),
    reveal: (runId: string, slot: number) => reveal(runId, slot),
  },
}))

beforeEach(() => {
  vi.clearAllMocks()
  state.mockResolvedValue(makeState())
  prizes.mockResolvedValue(makePrizes())
  start.mockResolvedValue(makeRun())
})

describe('WeeklyLotteryScreen', () => {
  it('starts a weekly board with nine closed slots', async () => {
    renderLottery(<WeeklyLotteryScreen />)

    await userEvent.click(
      await screen.findByRole('button', { name: 'Открыть поле' }),
    )

    expect(start).toHaveBeenCalledOnce()
    expect(
      await screen.findAllByRole('button', { name: /Открыть ячейку/ }),
    ).toHaveLength(9)
  })

  it('does not render symbols hidden by the backend', async () => {
    state.mockResolvedValue(makeState({ run: makeRun() }))

    renderLottery(<WeeklyLotteryScreen />)

    await screen.findByText('Открывайте ячейки')
    expect(screen.queryByText('Велосипед')).not.toBeInTheDocument()
    expect(screen.queryByText('Смартфон')).not.toBeInTheDocument()
  })

  it('reveals the selected slot and sends its index', async () => {
    const run = makeRun()
    state.mockResolvedValue(makeState({ run }))
    reveal.mockResolvedValue({
      index: 4,
      symbol: 'bicycle',
      run: makeRun({
        slots: run.slots.map((slot) =>
          slot.index === 4
            ? { index: 4, opened: true, symbol: 'bicycle' }
            : slot,
        ),
      }),
    })

    renderLottery(<WeeklyLotteryScreen />)
    await userEvent.click(
      await screen.findByRole('button', { name: 'Открыть ячейку 5' }),
    )

    await waitFor(() => {
      expect(reveal).toHaveBeenCalledWith('run123456789', 4)
    })
    expect(
      await screen.findByRole('button', { name: 'Ячейка 5: Велосипед' }),
    ).toBeDisabled()
  })

  it('shows the persisted promo code after a win', async () => {
    const slots = makeRun().slots.map((slot) =>
      slot.index < 3
        ? { index: slot.index, opened: true as const, symbol: 'bicycle' as const }
        : slot,
    )
    state.mockResolvedValue(
      makeState({
        available: false,
        run: makeRun({
          state: 'won',
          slots,
          prize: {
            id: 'weekly_bicycle_5',
            title: 'Скидка 5% на спорт и отдых',
            description: 'Скидка на товар из категории спорта и отдыха',
            benefit_type: 'percent_discount',
            benefit_value: 5,
            scope_type: 'category',
            scope_value: 'sport',
            code: 'PROMO-SIGNED',
            expires_at: '2026-09-11T12:00:00Z',
          },
        }),
      }),
    )

    renderLottery(<WeeklyLotteryScreen />)

    expect(
      await screen.findByText('Три совпадения — вы выиграли!'),
    ).toBeInTheDocument()
    expect(screen.getByText('PROMO-SIGNED')).toBeInTheDocument()
  })

  it('restores a completed losing board without enabling its slots', async () => {
    const symbols = [
      'bicycle',
      'bicycle',
      'smartphone',
      'smartphone',
      'sofa',
      'sofa',
      'delivery',
      'delivery',
      'promotion',
    ] as const
    state.mockResolvedValue(
      makeState({
        available: false,
        run: makeRun({
          state: 'lost',
          slots: symbols.map((symbol, index) => ({
            index,
            opened: true as const,
            symbol,
          })),
        }),
      }),
    )

    renderLottery(<WeeklyLotteryScreen />)

    expect(
      await screen.findByText('На этой неделе без совпадений'),
    ).toBeInTheDocument()
    expect(
      screen.getAllByRole('button', { name: /Ячейка \d:/ }),
    ).toHaveLength(9)
    for (const slot of screen.getAllByRole('button', { name: /Ячейка \d:/ })) {
      expect(slot).toBeDisabled()
    }
  })
})
