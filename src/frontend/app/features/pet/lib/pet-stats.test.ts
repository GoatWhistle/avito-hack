import { describe, expect, it } from 'vitest'
import {
  canCheckInToday,
  levelProgress,
  MAX_LEVEL,
  nextStep,
  statTone,
  statViews,
} from './pet-stats'
import type { Pet } from '#/features/pet/types'

const basePet: Pet = {
  id: 'pet-1',
  user_id: 'user-1',
  name: 'Ноти',
  stage: 'baby',
  state: 'neutral',
  level: 3,
  xp: 15,
  next_level_xp: 22,
  satiety: 70,
  happiness: 80,
  energy: 90,
  streak_days: 4,
  freezes: 1,
  is_hatched: true,
  hatched_at: '2026-08-01T10:00:00Z',
  last_checkin_date: null,
}

describe('statTone', () => {
  it.each([
    [10, 'critical'],
    [29, 'critical'],
    [30, 'warning'],
    [59, 'warning'],
    [60, 'good'],
    [100, 'good'],
  ])('maps %i to %s', (value, expected) => {
    expect(statTone(value)).toBe(expected)
  })
})

describe('statViews', () => {
  it('returns three clamped stats', () => {
    const views = statViews({ ...basePet, satiety: 120, energy: -5 })

    expect(views.map((view) => view.key)).toEqual([
      'satiety',
      'happiness',
      'energy',
    ])
    expect(views[0].value).toBe(100)
    expect(views[2].value).toBe(0)
    expect(views[2].tone).toBe('critical')
  })
})

describe('levelProgress', () => {
  it('computes the remaining xp and percent', () => {
    const progress = levelProgress(basePet)

    expect(progress.xpToNext).toBe(7)
    expect(progress.percent).toBe(68)
    expect(progress.isMaxLevel).toBe(false)
  })

  it('reports max level when the level cap is reached', () => {
    const progress = levelProgress({
      ...basePet,
      level: MAX_LEVEL,
      xp: 500,
      next_level_xp: 0,
    })

    expect(progress.isMaxLevel).toBe(true)
    expect(progress.percent).toBe(100)
    expect(progress.xpToNext).toBe(0)
  })
})

describe('canCheckInToday', () => {
  const now = new Date('2026-08-07T12:00:00Z')

  it('allows check-in when there is no previous one', () => {
    expect(canCheckInToday(basePet, now)).toBe(true)
  })

  it('blocks check-in when already done today', () => {
    const pet = { ...basePet, last_checkin_date: now.toISOString() }

    expect(canCheckInToday(pet, now)).toBe(false)
  })

  it('allows check-in on the following day', () => {
    const pet = { ...basePet, last_checkin_date: '2026-08-06T12:00:00Z' }

    expect(canCheckInToday(pet, now)).toBe(true)
  })
})

describe('nextStep', () => {
  it('prioritises hatching', () => {
    expect(
      nextStep({ pet: { ...basePet, is_hatched: false }, canCheckIn: false }),
    ).toBe('hatch')
  })

  it('prioritises check-in when available', () => {
    expect(nextStep({ pet: basePet, canCheckIn: true })).toBe('checkIn')
  })

  it('suggests feeding when hungry', () => {
    expect(
      nextStep({ pet: { ...basePet, satiety: 10 }, canCheckIn: false }),
    ).toBe('feedHungry')
  })

  it('suggests rest when energy is gone', () => {
    expect(
      nextStep({ pet: { ...basePet, energy: 5 }, canCheckIn: false }),
    ).toBe('rest')
  })

  it('suggests cheering up on low mood', () => {
    expect(
      nextStep({ pet: { ...basePet, happiness: 40 }, canCheckIn: false }),
    ).toBe('cheerUp')
  })

  it('reports max level for a legend', () => {
    expect(
      nextStep({ pet: { ...basePet, level: MAX_LEVEL }, canCheckIn: false }),
    ).toBe('maxLevel')
  })

  it('falls back to keep going', () => {
    expect(nextStep({ pet: basePet, canCheckIn: false })).toBe('keepGoing')
  })
})
