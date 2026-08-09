import { screen, waitFor } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it } from 'vitest'
import { changeLocale, initI18n } from '#/i18n'
import { renderWithProviders } from '#/features/items/components/test-utils'
import { toRewardProgress } from '#/features/rewards/lib'
import type { BadgeItem, RewardCatalogItem } from '#/features/rewards/types'
import { AchievementCard } from './AchievementCard'
import { RewardCard } from './RewardCard'

const badge: BadgeItem = {
  id: 'explorer',
  name: 'Исследователь',
  description: 'Добавить 5 объявлений в избранное',
  icon_url: '',
  earned_at: '',
  progress_current: 2,
  progress_target: 5,
}

const reward: RewardCatalogItem = {
  id: 'xl_listing_discount_50',
  title: 'Скидка 50% на XL-объявление',
  description: 'Скидка на услугу «XL-объявление»',
  kind: 'promo',
  condition_type: 'level',
  condition_value: 11,
  unlocked: false,
  claimed: false,
  status: '',
  progress_current: 4,
  progress_target: 11,
}

beforeEach(() => {
  initI18n()
})

afterEach(async () => {
  await changeLocale('ru')
})

describe('catalog text follows the interface language', () => {
  it('swaps badge name and description when the language changes', async () => {
    renderWithProviders(
      <ul>
        <AchievementCard badge={badge} />
      </ul>,
    )

    expect(await screen.findByText('Исследователь')).toBeInTheDocument()
    expect(
      screen.getByText('Добавить 5 объявлений в избранное'),
    ).toBeInTheDocument()

    await changeLocale('en')

    await waitFor(() => {
      expect(screen.getByText('Explorer')).toBeInTheDocument()
    })
    expect(screen.getByText('Add 5 listings to favorites')).toBeInTheDocument()
    expect(screen.queryByText('Исследователь')).not.toBeInTheDocument()
  })

  it('swaps reward title and description when the language changes', async () => {
    renderWithProviders(<RewardCard entry={toRewardProgress(reward)} />)

    expect(
      await screen.findByText('Скидка 50% на XL-объявление'),
    ).toBeInTheDocument()

    await changeLocale('en')

    await waitFor(() => {
      expect(screen.getByText('50% off an XL listing')).toBeInTheDocument()
    })
    expect(
      screen.getByText('Discount on the "XL listing" service'),
    ).toBeInTheDocument()
    expect(
      screen.queryByText('Скидка 50% на XL-объявление'),
    ).not.toBeInTheDocument()
  })

  it('keeps backend text for an id missing from the catalog', async () => {
    renderWithProviders(
      <ul>
        <AchievementCard badge={{ ...badge, id: 'brand_new_badge' }} />
      </ul>,
    )

    await changeLocale('en')

    await waitFor(() => {
      expect(screen.getByText('Исследователь')).toBeInTheDocument()
    })
    expect(
      screen.getByText('Добавить 5 объявлений в избранное'),
    ).toBeInTheDocument()
  })
})
