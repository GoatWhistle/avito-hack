import { screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { groupRewards } from '#/features/rewards'
import type { RewardCatalogItem } from '#/features/rewards'
import type * as RewardsModule from '#/features/rewards'
import { nextRewards, PetRewardsPanel } from './PetRewardsPanel'
import { renderWithProviders } from './test-utils'

const catalog = vi.fn()

vi.mock('#/features/rewards', async (importOriginal) => {
  const actual = await importOriginal<typeof RewardsModule>()

  return { ...actual, useRewardCatalog: () => catalog() }
})

const makeReward = (
  overrides: Partial<RewardCatalogItem> = {},
): RewardCatalogItem => ({
  id: 'reward-1',
  title: 'Промокод на 500 ₽',
  description: 'Скидка на продвижение',
  kind: 'promo',
  condition_type: 'level',
  condition_value: 10,
  unlocked: false,
  claimed: false,
  status: '',
  progress_current: 5,
  progress_target: 10,
  ...overrides,
})

beforeEach(() => {
  vi.clearAllMocks()
})

describe('nextRewards', () => {
  it('drops unlocked and claimed rewards', () => {
    const groups = groupRewards([
      makeReward({ id: 'a', unlocked: true }),
      makeReward({ id: 'b', claimed: true }),
      makeReward({ id: 'c' }),
    ])

    expect(nextRewards(groups).map((entry) => entry.reward.id)).toEqual(['c'])
  })

  it('sorts the closest rewards first and limits the list', () => {
    const groups = groupRewards([
      makeReward({ id: 'far', progress_current: 1, progress_target: 10 }),
      makeReward({ id: 'near', progress_current: 9, progress_target: 10 }),
      makeReward({ id: 'mid', progress_current: 5, progress_target: 10 }),
      makeReward({ id: 'tail', progress_current: 0, progress_target: 10 }),
    ])

    expect(nextRewards(groups).map((entry) => entry.reward.id)).toEqual([
      'near',
      'mid',
      'far',
    ])
  })
})

describe('PetRewardsPanel', () => {
  it('lists the closest rewards with their progress', () => {
    catalog.mockReturnValue({
      data: groupRewards([makeReward()]),
      isPending: false,
      isError: false,
    })

    renderWithProviders(<PetRewardsPanel />)

    expect(
      screen.getByRole('region', { name: 'Ближайшие награды' }),
    ).toBeInTheDocument()
    expect(screen.getByText('Промокод на 500 ₽')).toBeInTheDocument()
    expect(screen.getByText('5 / 10')).toBeInTheDocument()
    expect(screen.getByRole('meter')).toHaveAttribute('aria-valuenow', '50')
  })

  it('reports an empty state when everything is unlocked', () => {
    catalog.mockReturnValue({
      data: groupRewards([makeReward({ unlocked: true })]),
      isPending: false,
      isError: false,
    })

    renderWithProviders(<PetRewardsPanel />)

    expect(screen.getByText('Все награды уже открыты')).toBeInTheDocument()
  })

  it('reports a load failure without breaking the panel', () => {
    catalog.mockReturnValue({
      data: undefined,
      isPending: false,
      isError: true,
    })

    renderWithProviders(<PetRewardsPanel />)

    expect(screen.getByText('Не удалось загрузить награды')).toBeInTheDocument()
  })

  it('links to the full rewards page', () => {
    catalog.mockReturnValue({
      data: groupRewards([makeReward()]),
      isPending: false,
      isError: false,
    })

    renderWithProviders(<PetRewardsPanel />)

    expect(screen.getByRole('link', { name: 'Все' })).toHaveAttribute(
      'href',
      '/pet/rewards',
    )
  })
})
