import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { AxiosError, AxiosHeaders } from 'axios'
import { renderWithProviders } from '#/features/items/components/test-utils'
import { RewardTrack } from './RewardTrack'
import { makeCatalogItem, makeMyReward } from './test-utils'

const catalog = vi.fn()
const mine = vi.fn()
const activate = vi.fn()

vi.mock('#/features/rewards/repository', () => ({
  rewardsRepository: {
    catalog: (signal?: AbortSignal) => catalog(signal),
    mine: (signal?: AbortSignal) => mine(signal),
    activate: (id: string) => activate(id),
    badges: vi.fn(),
  },
}))

const conflictError = () => {
  const headers = new AxiosHeaders()

  return new AxiosError('Request failed', 'ERR_BAD_REQUEST', undefined, null, {
    status: 409,
    statusText: 'Conflict',
    headers,
    config: { headers },
    data: {
      error: { code: 'conflict', message: 'reward has already been activated' },
    },
  })
}

beforeEach(() => {
  vi.clearAllMocks()
  catalog.mockResolvedValue([
    makeCatalogItem({
      id: 'level-1',
      title: 'Статус «Новичок»',
      condition_value: 1,
      progress_target: 1,
      unlocked: true,
    }),
    makeCatalogItem({
      id: 'level-9',
      title: 'Скидка 30%',
      condition_value: 9,
      progress_target: 9,
      unlocked: true,
      claimed: true,
    }),
    makeCatalogItem({
      id: 'level-15',
      title: 'Бесплатная доставка',
      condition_value: 15,
      progress_target: 15,
      unlocked: false,
      progress_current: 3,
    }),
  ])
  mine.mockResolvedValue([])
  activate.mockResolvedValue({
    reward_id: 'level-1',
    code: 'AVITO-NEW-2026',
    status: 'activated',
  })
})

describe('RewardTrack', () => {
  it('shows a skeleton while loading', () => {
    catalog.mockReturnValue(new Promise(() => {}))
    renderWithProviders(<RewardTrack />)

    expect(screen.getByRole('status', { busy: true })).toBeInTheDocument()
  })

  it('renders the level track sorted with a state for every reward', async () => {
    renderWithProviders(<RewardTrack />)

    const cards = await screen.findAllByTestId('reward-level-card')

    expect(cards).toHaveLength(3)
    expect(cards.map((card) => card.dataset.state)).toEqual([
      'available',
      'activated',
      'locked',
    ])
    expect(screen.getByText('Забрано 1 из 3')).toBeInTheDocument()
  })

  it('offers a claim button only for claimable rewards', async () => {
    renderWithProviders(<RewardTrack />)

    await screen.findAllByTestId('reward-level-card')

    expect(screen.getAllByRole('button', { name: 'Забрать' })).toHaveLength(1)
    expect(screen.getByText('Получено')).toBeInTheDocument()
    expect(screen.getByText('3/15 ур.')).toBeInTheDocument()
  })

  it('reveals the promo code returned by a successful activation', async () => {
    renderWithProviders(<RewardTrack />)

    await userEvent.click(
      await screen.findByRole('button', { name: 'Забрать' }),
    )

    await waitFor(() => expect(activate).toHaveBeenCalledWith('level-1'))
    expect(await screen.findByText('AVITO-NEW-2026')).toBeInTheDocument()
    expect(
      screen.getByText('Одноразовый код — работает один раз'),
    ).toBeInTheDocument()
  })

  it('copies the promo code to the clipboard', async () => {
    const writeText = vi.fn().mockResolvedValue(undefined)
    Object.assign(navigator, { clipboard: { writeText } })

    renderWithProviders(<RewardTrack />)
    await userEvent.click(
      await screen.findByRole('button', { name: 'Забрать' }),
    )
    await screen.findByText('AVITO-NEW-2026')

    await userEvent.click(
      screen.getByRole('button', { name: 'Скопировать код' }),
    )

    expect(writeText).toHaveBeenCalledWith('AVITO-NEW-2026')
    expect(await screen.findByText('Код скопирован')).toBeInTheDocument()
  })

  it('explains a 409 conflict on repeated activation', async () => {
    activate.mockRejectedValue(conflictError())
    renderWithProviders(<RewardTrack />)

    await userEvent.click(
      await screen.findByRole('button', { name: 'Забрать' }),
    )

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Этот промокод уже активирован',
    )
    expect(screen.queryByTestId('reward-code')).not.toBeInTheDocument()
  })

  it('shows an already granted code without asking to claim again', async () => {
    mine.mockResolvedValue([
      makeMyReward({
        reward_id: 'level-9',
        status: 'activated',
        code: 'OLD-CODE-9',
      }),
    ])
    renderWithProviders(<RewardTrack />)

    expect(await screen.findByText('OLD-CODE-9')).toBeInTheDocument()
  })

  it('shows an empty state when no rewards exist', async () => {
    catalog.mockResolvedValue([])
    renderWithProviders(<RewardTrack />)

    expect(
      await screen.findByText('Награды пока не добавлены'),
    ).toBeInTheDocument()
  })

  it('shows an error with a retry action', async () => {
    catalog.mockRejectedValue(new Error('offline'))
    renderWithProviders(<RewardTrack />)

    expect(await screen.findByRole('alert')).toBeInTheDocument()

    catalog.mockResolvedValue([makeCatalogItem({ id: 'level-1' })])
    await userEvent.click(screen.getByRole('button', { name: 'Повторить' }))

    expect(await screen.findByTestId('reward-level-card')).toBeInTheDocument()
  })
})
