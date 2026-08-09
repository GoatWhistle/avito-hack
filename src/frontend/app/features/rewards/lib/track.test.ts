import { describe, expect, it } from 'vitest'
import {
  makeCatalogItem,
  makeMyReward,
} from '#/features/rewards/components/test-utils'
import {
  buildRewardTrack,
  isActivatable,
  trackStateOf,
  trackSummary,
} from './track'

describe('trackStateOf', () => {
  it('marks a reward locked when the condition is not met', () => {
    const reward = makeCatalogItem({ unlocked: false, claimed: false })

    expect(trackStateOf(reward, null)).toBe('locked')
  })

  it('marks an unlocked but unclaimed reward as available', () => {
    const reward = makeCatalogItem({ unlocked: true, claimed: false })

    expect(trackStateOf(reward, null)).toBe('available')
  })

  it('prefers the granted record over catalog flags', () => {
    const reward = makeCatalogItem({ unlocked: true, claimed: false })

    expect(trackStateOf(reward, makeMyReward({ status: 'granted' }))).toBe(
      'granted',
    )
    expect(trackStateOf(reward, makeMyReward({ status: 'activated' }))).toBe(
      'activated',
    )
  })

  it('treats a past expiry date as expired', () => {
    const reward = makeCatalogItem()
    const granted = makeMyReward({
      status: 'granted',
      expires_at: '2020-01-01T00:00:00Z',
    })

    expect(trackStateOf(reward, granted)).toBe('expired')
  })

  it('falls back to activated for a claimed catalog reward', () => {
    expect(trackStateOf(makeCatalogItem({ claimed: true }), null)).toBe(
      'activated',
    )
  })
})

describe('isActivatable', () => {
  it('allows only available and granted rewards to be activated', () => {
    expect(isActivatable('available')).toBe(true)
    expect(isActivatable('granted')).toBe(true)
    expect(isActivatable('activated')).toBe(false)
    expect(isActivatable('locked')).toBe(false)
    expect(isActivatable('expired')).toBe(false)
  })
})

describe('buildRewardTrack', () => {
  it('orders entries by level and attaches granted codes', () => {
    const entries = buildRewardTrack(
      [
        makeCatalogItem({ id: 'b', condition_value: 5, progress_target: 5 }),
        makeCatalogItem({ id: 'a', condition_value: 1, progress_target: 1 }),
      ],
      [makeMyReward({ reward_id: 'a', code: 'CODE-A' })],
    )

    expect(entries.map((entry) => entry.level)).toEqual([1, 5])
    expect(entries[0].code).toBe('CODE-A')
    expect(entries[0].state).toBe('granted')
    expect(entries[1].code).toBe('')
  })

  it('works without any granted rewards', () => {
    const entries = buildRewardTrack([makeCatalogItem({ unlocked: false })])

    expect(entries).toHaveLength(1)
    expect(entries[0].state).toBe('locked')
    expect(entries[0].granted).toBeNull()
  })
})

describe('trackSummary', () => {
  it('counts claimed and available rewards', () => {
    const entries = buildRewardTrack(
      [
        makeCatalogItem({ id: 'a', condition_value: 1, unlocked: true }),
        makeCatalogItem({ id: 'b', condition_value: 2, unlocked: true }),
        makeCatalogItem({ id: 'c', condition_value: 3, unlocked: false }),
      ],
      [makeMyReward({ reward_id: 'a', status: 'activated' })],
    )

    expect(trackSummary(entries)).toEqual({
      claimed: 1,
      available: 1,
      total: 3,
    })
  })
})
