import type {
  BadgeItem,
  MyRewardItem,
  RewardCatalogItem,
} from '#/features/rewards/types'

export const makeCatalogItem = (
  overrides: Partial<RewardCatalogItem> = {},
): RewardCatalogItem => ({
  id: 'reward-1',
  title: 'Скидка 10%',
  description: 'Промокод на продвижение',
  kind: 'promo',
  condition_type: 'level',
  condition_value: 3,
  unlocked: true,
  claimed: false,
  status: '',
  progress_current: 5,
  progress_target: 3,
  ...overrides,
})

export const makeMyReward = (
  overrides: Partial<MyRewardItem> = {},
): MyRewardItem => ({
  reward_id: 'reward-1',
  title: 'Скидка 10%',
  description: 'Промокод на продвижение',
  kind: 'promo',
  status: 'granted',
  code: 'PROMO-1234',
  granted_at: '2026-08-01T10:00:00Z',
  activated_at: null,
  expires_at: null,
  ...overrides,
})

export const makeBadge = (overrides: Partial<BadgeItem> = {}): BadgeItem => ({
  id: 'badge-1',
  name: 'Друг Енота',
  description: 'Кормить питомца 7 дней подряд',
  icon_url: '',
  earned_at: '2026-08-01T10:00:00Z',
  progress_current: 0,
  progress_target: 0,
  ...overrides,
})
